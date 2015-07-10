package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

const (
	DefaultChannel     = "#general"
	DefaultMsgGreeting = "Hello, boys!"
	DefaultMsgFarewell = "Bye, boys!"
)

type DumbSlut struct {
	api        *slack.Slack
	chSender   chan slack.OutgoingMessage
	chReceiver chan slack.SlackEvent
	commands   []Command
	config     *Config
	msgParams  slack.PostMessageParameters
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
	d.salute(d.config.MsgGreeting)
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

func (d *DumbSlut) initCommands() {
	for _, cfg := range d.config.Commands {
		d.AddCommand(NewCommand(cfg))
	}
}

// Start talking or i'll shoot
func (d *DumbSlut) Talk(message string) {
	channelId, timestamp, err := d.api.PostMessage(d.config.Channel, message, d.msgParams)
	if err != nil {
		log.Errorf("Error sending message: %s", err.Error())
	}

	log.Infof("Message successfully sent to %s at %s", channelId, timestamp)
}

func (d *DumbSlut) init() {
	d.api = slack.New(d.config.ApiToken)
	d.api.SetDebug(d.config.Debug)

	d.setDefaults()
	d.initCommands()

	if d.config.LogLevel != "" {
		log.SetLevel(stringToLogLevel(d.config.LogLevel))
	}

	d.handleRealtimeMessages()
	d.handleInterrupt()
}

func (d *DumbSlut) setDefaults() {
	// Default message parameters
	d.msgParams = slack.PostMessageParameters{
		AsUser: true,
	}

	if d.config.Channel == "" {
		d.config.Channel = DefaultChannel
	}

	if d.config.MsgGreeting == "" {
		d.config.MsgGreeting = DefaultMsgGreeting
	}

	if d.config.MsgFarewell == "" {
		d.config.MsgFarewell = DefaultMsgFarewell
	}
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

func (d *DumbSlut) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		d.salute(d.config.MsgFarewell)
		log.Infof("Finished...")
		os.Exit(1)
	}()
}

func (d *DumbSlut) salute(message string) {
	if !d.config.Shy {
		d.salute(message)
	}
}

func (d *DumbSlut) handleCommands(msg *slack.MessageEvent) {
	// Don't analyze own words
	info := d.api.GetInfo()
	if msg.UserId == info.User.Id {
		return
	}

	// Check trigger and respond
	for _, command := range d.commands {
		if command.Trigger(d, msg) {
			command.Respond(d, msg)
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
				//event := msg.Data.(*slack.PresenceChangeEvent)
				//log.Debugf("Presence Change: %v\n", event)
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
