package telegram

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
	"github.com/tucnak/telebot"
)

// Telegram adapter object
type Telegram struct {
	bot    *telebot.Bot
	cfg    config.Adapter
	userId string
}

// Create new Telegram adapter instance
func New(cfg config.Adapter) *Telegram {
	adapter := new(Telegram)
	adapter.cfg = cfg
	return adapter
}

// Initialize Telegram adapter
func (adapter *Telegram) Init() error {
	bot, err := telebot.NewBot(adapter.cfg.Token)
	if err != nil {
		return err
	}

	adapter.bot = bot
	return nil
}

// Listen to incoming events
func (adapter *Telegram) Listen(events chan *commons.Event) {
	messages := make(chan telebot.Message)

	for {
		adapter.bot.Listen(messages, time.Second)

		for message := range messages {
			log.WithField("adapter", adapter.cfg.Type).Debugf("Message: %v", message)
			events <- commons.NewEvent().FromTelegramMessage(message)
		}

		time.Sleep(time.Second)
	}
}

// Send message
func (adapter *Telegram) SendMessage(msg string, e *commons.Event) error {
	if e != nil && e.TgMsg.Sender.ID == adapter.bot.Identity.ID {
		return nil
	}

	options := telebot.SendOptions{}
	if e != nil {
		options.ReplyTo = e.TgMsg
	}

	return adapter.bot.SendMessage(e.TgMsg.Chat, msg, &options)
}
