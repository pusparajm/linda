package slack

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/dumbslut/config"
	"github.com/kpashka/dumbslut/event"
	wrapper "github.com/nlopes/slack"
)

const (
	DefaultChannel = "general"
)

// Slack backend object
type Slack struct {
	api        *wrapper.Slack
	cfg        config.Backend
	channels   []wrapper.Channel
	chSender   chan wrapper.OutgoingMessage
	chReceiver chan wrapper.SlackEvent
	msgParams  wrapper.PostMessageParameters
	userId     string
	users      []wrapper.User
}

// Create new Slack backend instance
func New(cfg config.Backend) *Slack {
	backend := new(Slack)
	backend.cfg = cfg
	return backend
}

// Initialize Slack backend
func (backend *Slack) Init() error {
	// Initialize Slack API
	backend.api = wrapper.New(backend.cfg.Token)

	// Set default channel
	if backend.cfg.Channel == "" {
		backend.cfg.Channel = DefaultChannel
	}

	// Set default message parameters
	backend.msgParams = wrapper.PostMessageParameters{
		AsUser:    true,
		LinkNames: 1,
		Markdown:  true,
	}

	// Handle realtime messages
	err := backend.handleRealtimeMessages()
	if err != nil {
		return err
	}

	// Set own user ID
	info := backend.api.GetInfo()
	backend.userId = info.User.Id

	// Sync channels
	err = backend.syncChannels()
	if err != nil {
		return err
	}

	// Sync users
	err = backend.syncUsers()
	if err != nil {
		return err
	}

	return nil
}

// Listen to incoming events
func (backend *Slack) Listen(events chan *event.Event) {
	for {
		select {
		case msg := <-backend.chReceiver:
			switch msg.Data.(type) {
			case wrapper.HelloEvent:
				slackEvent := msg.Data.(wrapper.HelloEvent)
				log.WithField("backend", backend.cfg.Type).Debugf("Hello: %v", slackEvent)

			case *wrapper.MessageEvent:
				slackEvent := msg.Data.(*wrapper.MessageEvent)
				channelName := backend.getChannelName(slackEvent.ChannelId)

				// Operate only on selected channel
				if channelName == backend.cfg.Channel {
					log.WithField("backend", backend.cfg.Type).Debugf("Message: %v", slackEvent)
					e := MessageToEvent(slackEvent)
					events <- e
				}

			case *wrapper.PresenceChangeEvent:
				slackEvent := msg.Data.(*wrapper.PresenceChangeEvent)
				log.WithField("backend", backend.cfg.Type).Debugf("Presence Change: %v", slackEvent)

				// Get username for presence change event
				if slackEvent.UserId != backend.userId {
					e := PresenceToEvent(slackEvent)
					backend.syncUsers()
					e.Username = backend.getUsername(slackEvent.UserId)
					events <- e
				}

			case wrapper.LatencyReport:
				// Skip latency report

			case *wrapper.SlackWSError:
				err := msg.Data.(*wrapper.SlackWSError)
				log.WithField("backend", backend.cfg.Type).Debugf("Error: %d - %s", err.Code, err.Msg)

			default:
				log.WithField("backend", backend.cfg.Type).Debugf("Unexpected: %v", msg.Data)
			}
		}
	}
}

// Send message
func (backend *Slack) SendMessage(msg string, e *event.Event) error {
	if e != nil && e.UserId == backend.userId {
		return nil
	}

	channel := fmt.Sprintf("#%s", backend.cfg.Channel)
	_, _, err := backend.api.PostMessage(channel, msg, backend.msgParams)
	return err
}

// Get channel name from synchronized data
func (backend *Slack) getChannelName(channelId string) string {
	for _, channel := range backend.channels {
		if channelId == channel.Id {
			return channel.Name
		}
	}

	return ""
}

// Get username from synchronized data
func (backend *Slack) getUsername(userId string) string {
	for _, u := range backend.users {
		if userId == u.Id {
			return u.Name
		}
	}

	return ""
}

// Handle realtime messages
func (backend *Slack) handleRealtimeMessages() error {
	backend.chSender = make(chan wrapper.OutgoingMessage)
	backend.chReceiver = make(chan wrapper.SlackEvent)

	// Protocol and origin are optional
	wsAPI, err := backend.api.StartRTM("", "http://example.com:8080")
	if err != nil {
		return err
	}

	go wsAPI.HandleIncomingEvents(backend.chReceiver)
	go wsAPI.Keepalive(20 * time.Second)

	go func(wsAPI *wrapper.SlackWS, chSender chan wrapper.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, backend.chSender)

	return nil
}

// Synchronize channels
func (backend *Slack) syncChannels() error {
	channels, err := backend.api.GetChannels(true)
	if err != nil {
		log.WithFields(log.Fields{
			"backend": backend.cfg.Type,
			"error":   err.Error(),
		}).Error("Error at channel sync")
		return err
	}

	backend.channels = channels
	return nil
}

// Synchronize users
func (backend *Slack) syncUsers() error {
	// React to active/away
	users, err := backend.api.GetUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"backend": backend.cfg.Type,
			"error":   err.Error(),
		}).Error("Error at user sync")
		return err
	}

	backend.users = users
	return nil
}
