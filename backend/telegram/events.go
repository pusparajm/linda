package telegram

import (
	"strconv"
	"time"

	"github.com/kpashka/dumbslut/event"
	"github.com/tucnak/telebot"
)

// Convert telebot.Message to event.Event
func MessageToEvent(msg telebot.Message) *event.Event {
	e := event.New(event.TypeMessage)

	e.Channel = msg.Chat.Username
	e.Text = msg.Text
	e.Timestamp = time.Unix(int64(msg.Unixtime), 0).String()
	e.UserId = strconv.Itoa(msg.Sender.ID)
	e.Username = msg.Sender.Username
	return e
}
