package command

import (
	"github.com/kpashka/linda/command/artist"
	"github.com/kpashka/linda/command/bully"
	"github.com/kpashka/linda/command/postman"
	"github.com/kpashka/linda/command/proxy"
	"github.com/kpashka/linda/command/snitch"
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
	Run(params []string) (string, error)
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
