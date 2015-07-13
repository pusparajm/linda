package event

import (
	"github.com/nlopes/slack"
	"github.com/tucnak/telebot"
)

const (
	TypeMessage = iota
	TypeStatusChange
)

// Abstract event (message or presence change) containing all necessary info for reply
type Event struct {
	Type     int
	SlackMsg *slack.MessageEvent
	SlackPce *slack.PresenceChangeEvent
	TgMsg    telebot.Message

	Status   string
	Username string
	Text     string
}

// Convert slack.MessageEvent to event.Event
func FromSlackMessage(msg *slack.MessageEvent) *Event {
	e := new(Event)
	e.Type = TypeMessage
	e.SlackMsg = msg
	e.Text = msg.Text
	return e
}

// Convert slack.PresenceChangeEvent to event.Event
func FromSlackPresenceChange(msg *slack.PresenceChangeEvent) *Event {
	e := new(Event)
	e.Type = TypeStatusChange
	e.SlackPce = msg
	e.Status = msg.Presence
	return e
}

// Convert telebot.Message to event.Event
func FromTelegramMessage(msg telebot.Message) *Event {
	e := new(Event)
	e.Type = TypeMessage
	e.TgMsg = msg
	e.Text = msg.Text
	return e
}
