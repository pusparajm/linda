package adapters

import (
	"github.com/kpashka/linda/adapters/slack"
	"github.com/kpashka/linda/adapters/telegram"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

const (
	TypeSlack    = "slack"
	TypeTelegram = "telegram"
)

// Adapter interface
type Adapter interface {
	BotId() string
	Init() error
	Listen(events chan *commons.Event)
	Markdown() bool
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
