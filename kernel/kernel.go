package kernel

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/linda/adapters"
	"github.com/kpashka/linda/commands"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/hooks"
)

// Bot object
type Linda struct {
	adapter     adapters.Adapter
	commands    []commands.Command
	configs     []config.Command
	cfg         *config.Bot
	expressions []*regexp.Regexp
	hooks       map[string]hooks.Hook
}

// Create new Bot instance
func NewLinda(cfg *config.Bot) *Linda {
	linda := new(Linda)
	linda.cfg = cfg
	linda.commands = []commands.Command{}
	linda.hooks = map[string]hooks.Hook{}
	return linda
}

// Add command to bot
func (linda *Linda) AddCommand(cfg config.Command) *Linda {
	cmd := commands.New(cfg)
	if cmd != nil {
		r, err := regexp.Compile(cfg.Expression)
		if err != nil {
			log.Fatalf("Invalid expression in command %s", cfg.Name)
		}

		linda.expressions = append(linda.expressions, r)
		linda.commands = append(linda.commands, cmd)
		linda.configs = append(linda.configs, cfg)

		// Prepare hook instances
		for hook, _ := range cfg.Hooks {
			if _, ok := linda.hooks[hook]; !ok {
				log.Infof("Preloading hook %s", hook)
				linda.hooks[hook] = hooks.New(hook)
			}
		}
	}

	return linda
}

// Add multiple commands to bot
func (linda *Linda) AddCommands(configurations ...config.Command) *Linda {
	for _, cfg := range configurations {
		linda.AddCommand(cfg)
	}

	return linda
}

// Start listening to words and handling commands
func (linda *Linda) Start() {
	// Init linda
	linda.init()

	// Init backend
	err := linda.adapter.Init()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Error initializing backend")
	}

	// Handle process interruption
	linda.handleInterrupt()

	// Say hello if necessary
	err = linda.salute(linda.cfg.Params.Salutes.Greeting)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error while saying hello")
	}

	log.Infof("Listening to events...")
	events := make(chan *commons.Event)
	go linda.adapter.Listen(events)

	for {
		select {
		case e := <-events:
			linda.handleEvent(e)
		}
	}
}

func (linda *Linda) applyHooks(id int, params []string) []string {
	// Apply hooks to params if necessary
	for hook, hookParams := range linda.configs[id].Hooks {
		// Find hook
		if _, ok := linda.hooks[hook]; ok {
			log.Infof("Applying hook %s in command %s", hook, linda.configs[id].Name)

			// Range over params
			for _, num := range hookParams {
				if num < len(params) {
					params[num] = linda.hooks[hook].Fire(params[num])
				}
			}
		}
	}

	return params
}

func (linda *Linda) runCommand(id int, cmd commands.Command, e *commons.Event, params []string) {
	// Prepare
	params = linda.applyHooks(id, params)
	user := commons.NewUser(e.Username, linda.getNickname(e.Username))

	// Execute command
	response, err := cmd.Run(user, params)
	if err != nil {
		log.WithFields(log.Fields{
			"command": linda.configs[id].Name,
			"error":   err.Error(),
		}).Errorf("Error executing command")

		// Send error back to chat if tracing mode is active
		if linda.cfg.Params.Tracing {
			response = err.Error()
		} else {
			return
		}
	}

	// Send response
	err = linda.adapter.SendMessage(response, e)
	if err != nil {
		log.WithFields(log.Fields{
			"command": linda.configs[id].Name,
			"error":   err.Error(),
		}).Errorf("Error sending message")
		return
	}

	// Don't execute other commands if necessary
	if linda.cfg.Params.ExecutionMode == config.ExecutionModeFirst {
		return
	}
}

// Handle incoming events and pass them to commands
func (linda *Linda) handleEvent(e *commons.Event) {
	// Status change on top-priority
	if e.Type == commons.EventTypeStatusChange {
		linda.handleStatusChange(e)
		return
	}

	// Unless, execute commands
	for id, cmd := range linda.commands {
		// Define if should react to command by expression
		shouldReact, params := linda.shouldReact(id, cmd, e)
		if shouldReact {
			// Override for help command
			if linda.configs[id].Type == commands.TypeSnitch {
				linda.runCommand(id, cmd, e, linda.getDescriptions())
				continue
			}

			// Run other commands as usual
			linda.runCommand(id, cmd, e, params)
		}
	}
}

// Handle program interruption
func (linda *Linda) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		// Close signal channel
		<-c

		// Say goodbye if necessary
		err := linda.salute(linda.cfg.Params.Salutes.Greeting)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error while saying goodbye")
		}

		// Log and exit
		log.Infof("Interrupted...")
		os.Exit(1)
	}()
}

func (linda *Linda) getNickname(username string) string {
	if nickname, ok := linda.cfg.Params.Nicknames[username]; ok {
		return nickname
	}

	return username
}

// Temporary Slack-only builtin presence change handler
func (linda *Linda) handleStatusChange(e *commons.Event) {
	// Search for nicknames
	username := linda.getNickname(e.Username)

	var err error
	switch e.Status {
	case "active":
		err = linda.salute(fmt.Sprintf(linda.cfg.Params.Salutes.UserEntered, username))
	case "away":
		err = linda.salute(fmt.Sprintf(linda.cfg.Params.Salutes.UserLeft, username))
	}

	if err != nil {
		log.WithField("error", err.Error()).Errorf("Error while saluting user")
	}
}

func (linda *Linda) getDescriptions() []string {
	descriptions := []string{}

	for _, cfg := range linda.configs {
		description := fmt.Sprintf("*%s* - %s", cfg.Name, cfg.Description)
		descriptions = append(descriptions, description)
	}

	return descriptions
}

// Initialize bot adapter
func (linda *Linda) init() {
	// Logger setup
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(config.StringToLogLevel(linda.cfg.Params.LogLevel))
	log.Infof("Initializing...")

	// Set default execution mode
	linda.cfg.Params.ExecutionMode = config.GetExecutionMode(linda.cfg.Params.ExecutionMode)

	// Init adapter
	linda.adapter = adapters.New(linda.cfg.Adapter)
	if linda.adapter == nil {
		log.Fatal("Incorrect adapter configuration")
	}

	// Add commands
	linda.AddCommand(config.NewHelpCommand())
	linda.AddCommands(linda.cfg.Commands...)
}

func (linda *Linda) salute(message string) error {
	// Don't salute if shy
	if linda.cfg.Params.Shy {
		return nil
	}

	// Don't salute if empty message
	if len(message) == 0 {
		return nil
	}

	return linda.adapter.SendMessage(message, nil)
}

// Defines whether bot should react to event
func (linda *Linda) shouldReact(id int, cmd commands.Command, e *commons.Event) (bool, []string) {
	expression := linda.expressions[id]
	matches := expression.FindStringSubmatch(e.Text)

	if len(matches) > 0 {
		return true, matches
	}

	return false, matches
}
