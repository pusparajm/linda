package telegram

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/dumbslut/config"
	"github.com/kpashka/dumbslut/event"
	"github.com/tucnak/telebot"
)

// Telegram backend object
type Telegram struct {
	bot    *telebot.Bot
	cfg    config.Backend
	userId string
}

// Create new Telegram backend instance
func New(cfg config.Backend) *Telegram {
	backend := new(Telegram)
	backend.cfg = cfg
	return backend
}

// Initialize Telegram backend
func (backend *Telegram) Init() error {
	bot, err := telebot.NewBot(backend.cfg.Token)
	if err != nil {
		return err
	}

	backend.bot = bot
	return nil
}

// Listen to incoming events
func (backend *Telegram) Listen(events chan *event.Event) {
	messages := make(chan telebot.Message)

	for {
		backend.bot.Listen(messages, time.Second)

		for message := range messages {
			log.WithField("backend", backend.cfg.Type).Debugf("Message: %v", message)
			events <- event.FromTelegramMessage(message)
		}

		time.Sleep(time.Second)
	}
}

// Send message
func (backend *Telegram) SendMessage(msg string, e *event.Event) error {
	if e != nil && e.TgMsg.Sender.ID == backend.bot.Identity.ID {
		return nil
	}

	options := telebot.SendOptions{}
	if e != nil {
		options.ReplyTo = e.TgMsg
	}

	return backend.bot.SendMessage(e.TgMsg.Chat, msg, &options)
}
