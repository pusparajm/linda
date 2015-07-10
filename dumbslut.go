package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

const (
	DefaultChannel = "#general"
)

type DumbSlut struct {
	api        *slack.Slack
	chSender   chan slack.OutgoingMessage
	chReceiver chan slack.SlackEvent
	commands   []Command
	config     *Config
	msgParams  slack.PostMessageParameters
	userId     string
}

func NewDumbSlut(config *Config) *DumbSlut {
	d := new(DumbSlut)
	d.config = config
	d.commands = []Command{}
	return d
}

func (d *DumbSlut) Start() {
	log.Infof("Started...")
	d.init()

	// Greet everyone
	d.salute(d.config.Salutes.Greeting)

	// Set own user ID
	info := d.api.GetInfo()
	d.userId = info.User.Id

	// Main handler
	d.listenAndServe()
}

func (d *DumbSlut) AddCommand(cmd Command) {
	if cmd != nil {
		d.commands = append(d.commands, cmd)
	}
}

func (d *DumbSlut) AddCommands(commands ...Command) {
	for _, cmd := range commands {
		d.AddCommand(cmd)
	}
}

// Start talking or i'll shoot
func (d *DumbSlut) Talk(message string) {
	if message == "" {
		return
	}

	channelId, timestamp, err := d.api.PostMessage(d.config.Channel, message, d.msgParams)
	if err != nil {
		log.Errorf("Error sending message: %s", err.Error())
	}

	log.Infof("Message successfully sent to %s at %s", channelId, timestamp)
}

func (d *DumbSlut) init() {
	// Default values
	log.SetLevel(stringToLogLevel(d.config.LogLevel))
	d.msgParams = slack.PostMessageParameters{AsUser: true}

	if d.config.Channel == "" {
		d.config.Channel = DefaultChannel
	}

	// Init API
	d.api = slack.New(d.config.ApiToken)
	d.api.SetDebug(d.config.Debug)

	// Add commands
	for _, cfg := range d.config.Commands {
		d.AddCommand(NewCommand(cfg))
	}

	// Handle stuff
	d.handleRealtimeMessages()
	d.handleInterrupt()
}

func (d *DumbSlut) salute(message string) {
	if !d.config.Shy {
		d.Talk(message)
	}
}

func (d *DumbSlut) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		d.salute(d.config.Salutes.Farewell)
		log.Infof("Finished...")
		os.Exit(1)
	}()
}

func (d *DumbSlut) handleRealtimeMessages() {
	d.chSender = make(chan slack.OutgoingMessage)
	d.chReceiver = make(chan slack.SlackEvent)

	wsAPI, err := d.api.StartRTM("", "http://example.com")
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	go wsAPI.HandleIncomingEvents(d.chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, d.chSender)
}

func (d *DumbSlut) handleCommands(msg *slack.MessageEvent) {
	if msg.UserId == d.userId {
		return
	}

	// Check trigger and respond
	for _, command := range d.commands {
		if command.Trigger(d, msg) {
			log.Infof("Triggered by %s command", command.GetName())
			command.Execute(d, msg)
			log.Infof("Executed %s command", command.GetName())
		}
	}
}

func (d *DumbSlut) handlePresenceChange(event *slack.PresenceChangeEvent) {
	// React to active/away
	users, err := d.api.GetUsers()
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("Error at statusReaction()")
		return
	}

	// Search for users
	var user slack.User
	for _, u := range users {
		if event.UserId == u.Id && u.Id != d.userId {
			user = u
			break
		}
	}

	if user.Id != "" {
		username := user.RealName
		if nickname, ok := d.config.Nicknames[user.Name]; ok {
			username = nickname
		}

		switch event.Presence {
		case "active":
			d.salute(fmt.Sprintf(d.config.Salutes.UserActive, username))
		case "away":
			d.salute(fmt.Sprintf(d.config.Salutes.UserAway, username))
		}
	}
}

func (d *DumbSlut) listenAndServe() {
	for {
		select {
		case msg := <-d.chReceiver:
			switch msg.Data.(type) {
			case slack.HelloEvent:
				// Ignore hello
			case *slack.MessageEvent:
				event := msg.Data.(*slack.MessageEvent)
				log.Debugf("Message: %v\n", event)
				d.handleCommands(event)
			case *slack.PresenceChangeEvent:
				event := msg.Data.(*slack.PresenceChangeEvent)
				log.Debugf("Presence Change: %v\n", event)
				d.handlePresenceChange(event)
			case slack.LatencyReport:
				//event := msg.Data.(slack.LatencyReport)
				//log.Debugf("Current latency: %v\n", event.Value)
			case *slack.SlackWSError:
				//error := msg.Data.(*slack.SlackWSError)
				//log.Debugf("Error: %d - %s\n", error.Code, error.Msg)
			default:
				log.Debugf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
