package adapters

import (
	"github.com/kpashka/linda/adapters/slack"
	"github.com/kpashka/linda/adapters/telegram"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

const (
	TypeSlack    = "Slack"
	TypeTelegram = "Telegram"
)

// Adapter interface
type Adapter interface {
	Init() error
	Listen(events chan *commons.Event)
	SendMessage(msg string, e *commons.Event) error
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
