package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/linda/backend"
	"github.com/kpashka/linda/command"
	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/event"
)

// Bot object
type Linda struct {
	backend     backend.Backend
	commands    []command.Command
	configs     []config.Command
	cfg         *config.Bot
	expressions []*regexp.Regexp
}

// Create new Bot instance
func NewLinda(cfg *config.Bot) *Linda {
	linda := new(Linda)
	linda.cfg = cfg
	linda.commands = []command.Command{}
	return linda
}

// Add command to bot
func (linda *Linda) AddCommand(cfg config.Command) *Linda {
	cmd := command.New(cfg)
	if cmd != nil {
		r, err := regexp.Compile(cfg.Expression)
		if err != nil {
			log.Fatalf("Invalid expression in command %s", cfg.Name)
		}

		linda.expressions = append(linda.expressions, r)
		linda.commands = append(linda.commands, cmd)
		linda.configs = append(linda.configs, cfg)
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
	err := linda.backend.Init()
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
	events := make(chan *event.Event)
	go linda.backend.Listen(events)

	for {
		select {
		case e := <-events:
			linda.handleEvent(e)
		}
	}
}

func (linda *Linda) runCommand(id int, cmd command.Command, e *event.Event, params []string) {
	// Execute command
	response, err := cmd.Run(params)
	if err != nil {
		log.WithFields(log.Fields{
			"command": linda.configs[id].Name,
			"error":   err.Error(),
		}).Errorf("Error executing command")
		return
	}

	// Send response
	err = linda.backend.SendMessage(response, e)
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
func (linda *Linda) handleEvent(e *event.Event) {
	// Status change on top-priority
	if e.Type == event.TypeStatusChange {
		linda.handleStatusChange(e)
		return
	}

	// Unless, execute commands
	for id, cmd := range linda.commands {
		// Define if should react to command by expression
		shouldReact, params := linda.shouldReact(id, cmd, e)
		if shouldReact {
			// Override for help command
			if linda.configs[id].Type == command.TypeSnitch {
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

// Temporary Slack-only builtin presence change handler
func (linda *Linda) handleStatusChange(e *event.Event) {
	// Search for nicknames
	username := e.Username
	if nickname, ok := linda.cfg.Params.Nicknames[username]; ok {
		username = nickname
	}

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

// Initialize bot backend
func (linda *Linda) init() {
	// Logger setup
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(config.StringToLogLevel(linda.cfg.Params.LogLevel))
	log.Infof("Initializing...")

	// Set default execution mode
	linda.cfg.Params.ExecutionMode = config.GetExecutionMode(linda.cfg.Params.ExecutionMode)

	// Init backend
	linda.backend = backend.New(linda.cfg.Backend)
	if linda.backend == nil {
		log.Fatal("Incorrect backend configuration")
	}

	// Add commands
	linda.AddCommand(config.NewHelpCommand())
	linda.AddCommands(linda.cfg.Commands...)
}

func (linda *Linda) salute(message string) error {
	// Don't salute if shy
	if linda.cfg.Params.Shy || len(message) == 0 {
		return nil
	}

	return linda.backend.SendMessage(message, nil)
}

// Defines whether bot should react to event
func (linda *Linda) shouldReact(id int, cmd command.Command, e *event.Event) (bool, []string) {
	expression := linda.expressions[id]
	matches := expression.FindStringSubmatch(e.Text)

	if len(matches) > 0 {
		return true, matches
	}

	return false, matches
}
