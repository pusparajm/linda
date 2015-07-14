package slack

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/event"
	wrapper "github.com/nlopes/slack"
)

const (
	DefaultChannel = "general"
)

// Slack adapter object
type Slack struct {
	api        *wrapper.Slack
	cfg        config.Adapter
	channels   []wrapper.Channel
	chSender   chan wrapper.OutgoingMessage
	chReceiver chan wrapper.SlackEvent
	msgParams  wrapper.PostMessageParameters
	userId     string
	users      []wrapper.User
}

// Create new Slack adapter instance
func New(cfg config.Adapter) *Slack {
	adapter := new(Slack)
	adapter.cfg = cfg
	return adapter
}

// Initialize Slack adapter
func (adapter *Slack) Init() error {
	// Initialize Slack API
	adapter.api = wrapper.New(adapter.cfg.Token)

	// Set default channel
	if adapter.cfg.Channel == "" {
		adapter.cfg.Channel = DefaultChannel
	}

	// Set default message parameters
	adapter.msgParams = wrapper.PostMessageParameters{
		AsUser:    true,
		LinkNames: 1,
		Markdown:  true,
	}

	// Handle realtime messages
	err := adapter.handleRealtimeMessages()
	if err != nil {
		return err
	}

	// Set own user ID
	info := adapter.api.GetInfo()
	adapter.userId = info.User.Id

	// Sync channels
	err = adapter.syncChannels()
	if err != nil {
		return err
	}

	// Sync users
	err = adapter.syncUsers()
	if err != nil {
		return err
	}

	return nil
}

// Listen to incoming events
func (adapter *Slack) Listen(events chan *event.Event) {
	for {
		select {
		case msg := <-adapter.chReceiver:
			switch msg.Data.(type) {
			case wrapper.HelloEvent:
				slackEvent := msg.Data.(wrapper.HelloEvent)
				log.WithField("adapter", adapter.cfg.Type).Debugf("Hello: %v", slackEvent)

			case *wrapper.MessageEvent:
				slackEvent := msg.Data.(*wrapper.MessageEvent)
				channelName := adapter.getChannelName(slackEvent.ChannelId)

				// Operate only on selected channel
				if channelName == adapter.cfg.Channel {
					log.WithField("adapter", adapter.cfg.Type).Debugf("Message: %v", slackEvent)
					e := event.FromSlackMessage(slackEvent)
					events <- e
				}

			case *wrapper.PresenceChangeEvent:
				slackEvent := msg.Data.(*wrapper.PresenceChangeEvent)
				log.WithField("adapter", adapter.cfg.Type).Debugf("Presence Change: %v", slackEvent)

				// Get username for presence change event
				if slackEvent.UserId != adapter.userId {
					e := event.FromSlackPresenceChange(slackEvent)
					adapter.syncUsers()
					e.Username = adapter.getUsername(slackEvent.UserId)
					events <- e
				}

			case wrapper.LatencyReport:
				// Skip latency report

			case *wrapper.SlackWSError:
				err := msg.Data.(*wrapper.SlackWSError)
				log.WithField("adapter", adapter.cfg.Type).Debugf("Error: %d - %s", err.Code, err.Msg)

			default:
				log.WithField("adapter", adapter.cfg.Type).Debugf("Unexpected: %v", msg.Data)
			}
		}
	}
}

// Send message
func (adapter *Slack) SendMessage(msg string, e *event.Event) error {
	if e != nil && e.Type == event.TypeMessage && e.SlackMsg.UserId == adapter.userId {
		return nil
	}

	if e != nil && e.Type == event.TypeStatusChange && e.SlackPce.UserId == adapter.userId {
		return nil
	}

	channel := fmt.Sprintf("#%s", adapter.cfg.Channel)
	_, _, err := adapter.api.PostMessage(channel, msg, adapter.msgParams)
	return err
}

// Get channel name from synchronized data
func (adapter *Slack) getChannelName(channelId string) string {
	for _, channel := range adapter.channels {
		if channelId == channel.Id {
			return channel.Name
		}
	}

	return ""
}

// Get username from synchronized data
func (adapter *Slack) getUsername(userId string) string {
	for _, u := range adapter.users {
		if userId == u.Id {
			return u.Name
		}
	}

	return ""
}

// Handle realtime messages
func (adapter *Slack) handleRealtimeMessages() error {
	adapter.chSender = make(chan wrapper.OutgoingMessage)
	adapter.chReceiver = make(chan wrapper.SlackEvent)

	// Protocol and origin are optional
	wsAPI, err := adapter.api.StartRTM("", "http://example.com:8080")
	if err != nil {
		return err
	}

	go wsAPI.HandleIncomingEvents(adapter.chReceiver)
	go wsAPI.Keepalive(20 * time.Second)

	go func(wsAPI *wrapper.SlackWS, chSender chan wrapper.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, adapter.chSender)

	return nil
}

// Synchronize channels
func (adapter *Slack) syncChannels() error {
	channels, err := adapter.api.GetChannels(true)
	if err != nil {
		log.WithFields(log.Fields{
			"adapter": adapter.cfg.Type,
			"error":   err.Error(),
		}).Error("Error at channel sync")
		return err
	}

	adapter.channels = channels
	return nil
}

// Synchronize users
func (adapter *Slack) syncUsers() error {
	// React to active/away
	users, err := adapter.api.GetUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"adapter": adapter.cfg.Type,
			"error":   err.Error(),
		}).Error("Error at user sync")
		return err
	}

	adapter.users = users
	return nil
}
