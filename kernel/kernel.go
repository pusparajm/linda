package kernel

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/kpashka/linda/adapters"
	"github.com/kpashka/linda/commands"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/filters"
)

const (
	// Default web port to listen
	DefaultWebPort = "8080"
)

// Bot object
type Linda struct {
	adapter     adapters.Adapter
	commands    map[string]commands.Command
	configs     map[string]config.Command
	cfg         *config.Bot
	expressions map[string]*regexp.Regexp
	filters     map[string]filters.Filter
}

// Create new Bot instance
func NewLinda(cfg *config.Bot) *Linda {
	linda := new(Linda)
	linda.cfg = cfg

	linda.commands = map[string]commands.Command{}
	linda.configs = map[string]config.Command{}
	linda.expressions = map[string]*regexp.Regexp{}
	linda.filters = map[string]filters.Filter{}
	return linda
}

// Add command to bot
func (linda *Linda) AddCommand(id string, cfg config.Command) *Linda {
	cmd := commands.New(id, cfg)

	if cmd != nil {
		r, err := regexp.Compile(cfg.Expression)
		if err != nil {
			log.Fatalf("Invalid regular expression in command [%s]", id)
		}

		// Save all instances
		linda.expressions[id] = r
		linda.commands[id] = cmd
		linda.configs[id] = cfg

		// Prepare filters if necessary
		for _, filtersList := range cfg.Filters {
			for _, name := range filtersList {
				if name != "" {
					log.Infof("Preloading hook [%s]", name)
					linda.filters[name] = filters.New(name)
				}
			}
		}

	} else {
		log.Errorf("Failed to load command [%s]", id)
	}

	return linda
}

// Add multiple commands to bot
func (linda *Linda) AddCommands(configurations map[string]config.Command) *Linda {
	for id, cfg := range configurations {
		log.Infof("Initializing command [%s]", id)
		linda.AddCommand(id, cfg)
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
		log.WithField("error", err.Error()).Fatal("Error initializing adapter")
	}

	// Init web-listener
	//linda.initWeb()

	// Handle process interruption
	linda.handleInterrupt()

	// Say hello if necessary
	err = linda.salute(linda.cfg.Params.Salutes.Greeting)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error when greeting")
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

// Init basic web interface
func (linda *Linda) initWeb() {
	port := strconv.Itoa(linda.cfg.Params.HttpPort)

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = DefaultWebPort
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Linda! I'm currently running on %s, channel: %s", linda.cfg.Adapter.Type, linda.cfg.Adapter.Channel)
	})

	log.Infof("Web interface available at http://localhost:%s", port)
	go http.ListenAndServe(":"+port, nil)
}

// Initialize bot adapter
func (linda *Linda) init() {
	// Logger setup
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(config.StringToLogLevel(linda.cfg.Params.LogLevel))
	log.Infof("Initializing Linda...")

	// Set default execution mode
	linda.cfg.Params.ExecutionMode = config.GetExecutionMode(linda.cfg.Params.ExecutionMode)

	// Init adapter
	linda.adapter = adapters.New(linda.cfg.Adapter)
	if linda.adapter == nil {
		log.Fatal("Incorrect adapter configuration")
	}

	// Add commands
	linda.AddCommand("help", config.NewHelpCommand())
	linda.AddCommands(linda.cfg.Commands)
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

	return linda.sendMessage(message, nil)
}

func (linda *Linda) sendMessage(message string, event *commons.Event) error {
	return linda.adapter.SendMessage(message, event)
}

// Handle incoming events and pass them to commands
func (linda *Linda) handleEvent(e *commons.Event) {
	// Ignore own events
	if e.UserId == linda.adapter.BotId() {
		return
	}

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
			log.Infof("Found pattern for command [%s]", id)

			// Override for help command
			if linda.configs[id].Type == commands.TypeHelp {
				linda.runCommand(id, cmd, e, linda.getDescriptions())
			} else {
				// Run other commands as usual
				linda.runCommand(id, cmd, e, params)
			}

			// Don't execute other commands if necessary
			if linda.cfg.Params.ExecutionMode == config.ExecutionModeFirst {
				return
			}
		}
	}
}

// Defines whether bot should react to event
func (linda *Linda) shouldReact(id string, cmd commands.Command, e *commons.Event) (bool, []string) {
	expression := linda.expressions[id]
	matches := expression.FindStringSubmatch(e.Text)

	if len(matches) > 0 {
		return true, matches
	}

	return false, matches
}

// Execute specified command
func (linda *Linda) runCommand(id string, cmd commands.Command, e *commons.Event, params []string) {
	// Prepare
	params = linda.applyFilters(id, params)
	user := commons.NewUser(e.Username, linda.getNickname(e.Username))

	// Execute command
	response, err := cmd.Run(user, params)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Errorf("Error executing command [%s]", id)

		// Send error back to chat if tracing mode is active
		if linda.cfg.Params.Tracing {
			response = err.Error()
		} else {
			return
		}
	}

	// Send response
	err = linda.sendMessage(response, e)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Errorf("Error sending message from command [%s]", id)
		return
	}
}

// Apply specified filter
func (linda *Linda) applyFilter(filter, param string) string {
	if _, ok := linda.filters[filter]; ok {
		param = linda.filters[filter](param)
	}

	return param
}

// Apply filters for command parameters
func (linda *Linda) applyFilters(id string, params []string) []string {
	length := len(params)
	if length == 0 {
		return params
	}

	for p, hooks := range linda.configs[id].Filters {
		for _, name := range hooks {
			if p < length && name != "" {
				log.Infof("Applying filter [%s] for param [%d] in command [%s]", name, p, id)
				params[p] = linda.applyFilter(name, params[p])
			}
		}
	}

	return params
}

// Get user nickname from config
func (linda *Linda) getNickname(username string) string {
	if nickname, ok := linda.cfg.Params.Nicknames[username]; ok {
		return nickname
	}

	return username
}

// Get each command description for `help` command
func (linda *Linda) getDescriptions() []string {
	descriptions := []string{}
	markdown := linda.adapter.Markdown()

	for id, _ := range linda.commands {
		description := linda.getDescription(id, markdown)
		if description != "" {
			descriptions = append(descriptions, description)
		}
	}

	return descriptions
}

// Get specified command description
func (linda *Linda) getDescription(id string, markdown bool) string {
	if len(linda.configs[id].Description) == 0 {
		return linda.configs[id].Description
	}

	template := "[%s] - %s"
	if markdown {
		template = "*%s* - %s"
	}

	return fmt.Sprintf(template, id, linda.configs[id].Description)
}

// Temporary Slack-only builtin presence change handler
func (linda *Linda) handleStatusChange(e *commons.Event) {
	// Search for nicknames
	username := linda.getNickname(e.Username)

	var err error
	switch e.Status {
	case "active":
		err = linda.salute(fmt.Sprintf(linda.cfg.Params.Salutes.UserActive, username))
	case "away":
		err = linda.salute(fmt.Sprintf(linda.cfg.Params.Salutes.UserAway, username))
	}

	if err != nil {
		log.WithField("error", err.Error()).Errorf("Error while saluting user")
	}
}

// Handle program interruption / termination
func (linda *Linda) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		// Close signal channel
		<-c

		// Say goodbye if necessary
		err := linda.salute(linda.cfg.Params.Salutes.Farewell)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error while saying goodbye")
		}

		// Log and exit
		log.Infof("Interrupted...")
		os.Exit(1)
	}()
}
