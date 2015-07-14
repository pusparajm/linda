package commands

import (
	"github.com/kpashka/linda/commands/artist"
	"github.com/kpashka/linda/commands/bully"
	"github.com/kpashka/linda/commands/postman"
	"github.com/kpashka/linda/commands/proxy"
	"github.com/kpashka/linda/commands/snitch"
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

const (
	TypeArtist  = "Artist"
	TypeBully   = "Bully"
	TypePostman = "Postman"
	TypeProxy   = "Proxy"
	TypeSnitch  = "Snitch"
)

// Command interface
type Command interface {
	Run(user *commons.User, params []string) (string, error)
}

// Creates new Command instance
func New(cfg config.Command) Command {
	switch cfg.Type {
	case TypeArtist:
		return artist.New(cfg)
	case TypeBully:
		return bully.New(cfg)
	case TypePostman:
		return postman.New(cfg)
	case TypeProxy:
		return proxy.New(cfg)
	case TypeSnitch:
		return snitch.New(cfg)
	default:
		return nil
	}
}
