package slack

import (
	"github.com/kpashka/dumbslut/event"
	"github.com/nlopes/slack"
)

// Convert slack.MessageEvent to event.Event
func MessageToEvent(msg *slack.MessageEvent) *event.Event {
	e := event.New(event.TypeMessage)

	e.Channel = msg.ChannelId
	e.Text = msg.Text
	e.Timestamp = msg.Timestamp
	e.UserId = msg.UserId
	e.Username = msg.Username
	return e
}

// Convert slack.PresenceChangeEvent to event.Event
func PresenceToEvent(msg *slack.PresenceChangeEvent) *event.Event {
	e := event.New(event.TypeStatusChange)
	e.Status = msg.Presence
	e.UserId = msg.UserId
	return e
}
