package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/dumbslut/backend"
	"github.com/kpashka/dumbslut/command"
	"github.com/kpashka/dumbslut/config"
	"github.com/kpashka/dumbslut/event"
)

// Bot object
type Slut struct {
	backend     backend.Backend
	commands    []command.Command
	configs     []config.Command
	cfg         *config.Bot
	expressions []*regexp.Regexp
}

// Create new Bot instance
func NewSlut(cfg *config.Bot) *Slut {
	slut := new(Slut)
	slut.cfg = cfg
	slut.commands = []command.Command{}
	return slut
}

// Add command to bot
func (slut *Slut) AddCommand(cfg config.Command) *Slut {
	cmd := command.New(cfg)
	if cmd != nil {
		r, err := regexp.Compile(cfg.Expression)
		if err != nil {
			log.Fatalf("Invalid expression in command %s", cfg.Name)
		}

		slut.expressions = append(slut.expressions, r)
		slut.commands = append(slut.commands, cmd)
		slut.configs = append(slut.configs, cfg)
	}

	return slut
}

// Add multiple commands to bot
func (slut *Slut) AddCommands(configurations ...config.Command) *Slut {
	for _, cfg := range configurations {
		slut.AddCommand(cfg)
	}

	return slut
}

// Start listening to words and handling commands
func (slut *Slut) Start() {
	// Init slut
	slut.init()

	// Init backend
	err := slut.backend.Init()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Error initializing backend")
	}

	// Handle process interruption
	slut.handleInterrupt()

	// Say hello if necessary
	err = slut.salute(slut.cfg.Params.Salutes.Greeting)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error while saying hello")
	}

	log.Infof("Listening to events...")
	events := make(chan *event.Event)
	go slut.backend.Listen(events)

	for {
		select {
		case e := <-events:
			slut.handleEvent(e)
		}
	}
}

func (slut *Slut) runCommand(id int, cmd command.Command, e *event.Event, params []string) {
	// Execute command
	response, err := cmd.Run(params)
	if err != nil {
		log.WithFields(log.Fields{
			"command": slut.configs[id].Name,
			"error":   err.Error(),
		}).Errorf("Error executing command")
		return
	}

	// Send response
	err = slut.backend.SendMessage(response, e)
	if err != nil {
		log.WithFields(log.Fields{
			"command": slut.configs[id].Name,
			"error":   err.Error(),
		}).Errorf("Error sending message")
		return
	}

	// Don't execute other commands if necessary
	if slut.cfg.Params.ExecutionMode == config.ExecutionModeFirst {
		return
	}
}

// Handle incoming events and pass them to commands
func (slut *Slut) handleEvent(e *event.Event) {
	// Status change on top-priority
	if e.Type == event.TypeStatusChange {
		slut.handleStatusChange(e)
		return
	}

	// Unless, execute commands
	for id, cmd := range slut.commands {
		// Define if should react to command by expression
		shouldReact, params := slut.shouldReact(id, cmd, e)
		if shouldReact {
			// Override for help command
			if slut.configs[id].Type == command.TypeSnitch {
				slut.runCommand(id, cmd, e, slut.getDescriptions())
				continue
			}

			// Run other commands as usual
			slut.runCommand(id, cmd, e, params)
		}
	}
}

// Handle program interruption
func (slut *Slut) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		// Close signal channel
		<-c

		// Say goodbye if necessary
		err := slut.salute(slut.cfg.Params.Salutes.Greeting)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error while saying goodbye")
		}

		// Log and exit
		log.Infof("Interrupted...")
		os.Exit(1)
	}()
}

// Temporary Slack-only builtin presence change handler
func (slut *Slut) handleStatusChange(e *event.Event) {
	// Search for nicknames
	username := e.Username
	if nickname, ok := slut.cfg.Params.Nicknames[username]; ok {
		username = nickname
	}

	var err error
	switch e.Status {
	case "active":
		err = slut.salute(fmt.Sprintf(slut.cfg.Params.Salutes.UserEntered, username))
	case "away":
		err = slut.salute(fmt.Sprintf(slut.cfg.Params.Salutes.UserLeft, username))
	}

	if err != nil {
		log.WithField("error", err.Error()).Errorf("Error while saluting user")
	}
}

func (slut *Slut) getDescriptions() []string {
	descriptions := []string{}

	for _, cfg := range slut.configs {
		description := fmt.Sprintf("*%s* - %s", cfg.Name, cfg.Description)
		descriptions = append(descriptions, description)
	}

	return descriptions
}

// Initialize bot backend
func (slut *Slut) init() {
	// Logger setup
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(config.StringToLogLevel(slut.cfg.Params.LogLevel))
	log.Infof("Initializing...")

	// Set default execution mode
	slut.cfg.Params.ExecutionMode = config.GetExecutionMode(slut.cfg.Params.ExecutionMode)

	// Init backend
	slut.backend = backend.New(slut.cfg.Backend)
	if slut.backend == nil {
		log.Fatal("Incorrect backend configuration")
	}

	// Add commands
	slut.AddCommand(config.NewHelpCommand())
	slut.AddCommands(slut.cfg.Commands...)
}

func (slut *Slut) salute(message string) error {
	// Don't salute if shy
	if slut.cfg.Params.Shy {
		return nil
	}

	return slut.backend.SendMessage(message, nil)
}

// Defines whether bot should react to event
func (slut *Slut) shouldReact(id int, cmd command.Command, e *event.Event) (bool, []string) {
	expression := slut.expressions[id]
	matches := expression.FindStringSubmatch(e.Text)

	if len(matches) > 0 {
		return true, matches
	}

	return false, matches
}
