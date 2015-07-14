package adapters

import (
	"github.com/kpashka/linda/adapters/slack"
	"github.com/kpashka/linda/adapters/telegram"
	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/event"
)

const (
	TypeSlack    = "Slack"
	TypeTelegram = "Telegram"
)

// Adapter interface
type Adapter interface {
	Init() error
	Listen(events chan *event.Event)
	SendMessage(msg string, e *event.Event) error
}

// Creates new Adapter instance
func New(cfg config.Adapter) Adapter {
	switch cfg.Type {
	case TypeSlack:
		return slack.New(cfg)
	case TypeTelegram:
		return telegram.New(cfg)
	default:
		return nil
	}
}
