package backend

import (
	"github.com/kpashka/dumbslut/backend/slack"
	"github.com/kpashka/dumbslut/backend/telegram"
	"github.com/kpashka/dumbslut/config"
	"github.com/kpashka/dumbslut/event"
)

const (
	BackendTypeSlack    = "Slack"
	BackendTypeTelegram = "Telegram"
)

// Backend interface
type Backend interface {
	Init() error
	Listen(events chan *event.Event)
	SendMessage(msg string, e *event.Event) error
}

// Creates new Backend instance
func New(cfg config.Backend) Backend {
	switch cfg.Type {
	case BackendTypeSlack:
		return slack.New(cfg)
	case BackendTypeTelegram:
		return telegram.New(cfg)
	default:
		return nil
	}
}
