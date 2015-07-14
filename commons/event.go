package commons

import (
	"github.com/nlopes/slack"
	"github.com/tucnak/telebot"
)

const (
	EventTypeMessage = iota
	EventTypeStatusChange
)

// Abstract event (message or presence change) containing all necessary info for reply
type Event struct {
	Type int

	SlackMsg *slack.MessageEvent
	SlackPce *slack.PresenceChangeEvent
	TgMsg    telebot.Message

	Status   string
	Username string
	Text     string
}

func NewEvent() *Event {
	e := new(Event)
	return e
}

// Convert slack.MessageEvent to event.Event
func (e *Event) FromSlackMessage(msg *slack.MessageEvent) *Event {
	e.Type = EventTypeMessage
	e.SlackMsg = msg

	e.Text = msg.Text
	return e
}

// Convert slack.PresenceChangeEvent to event.Event
func (e *Event) FromSlackStatus(msg *slack.PresenceChangeEvent) *Event {
	e.Type = EventTypeStatusChange
	e.SlackPce = msg

	e.Status = msg.Presence
	return e
}

// Convert telebot.Message to event.Event
func (e *Event) FromTelegramMessage(msg telebot.Message) *Event {
	e.Type = EventTypeMessage

	e.TgMsg = msg
	e.Text = msg.Text
	e.Username = msg.Sender.Username
	return e
}
