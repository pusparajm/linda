package command

import (
	"github.com/kpashka/dumbslut/command/artist"
	"github.com/kpashka/dumbslut/command/bully"
	"github.com/kpashka/dumbslut/command/postman"
	"github.com/kpashka/dumbslut/command/proxy"
	"github.com/kpashka/dumbslut/config"
)

const (
	CommandTypeArtist  = "Artist"
	CommandTypeBully   = "Bully"
	CommandTypePostman = "Postman"
	CommandTypeProxy   = "Proxy"
)

// Command interface
type Command interface {
	Run(params []string) (string, error)
}

// Creates new Command instance
func New(cfg config.Command) Command {
	switch cfg.Type {
	case CommandTypeArtist:
		return artist.New(cfg)
	case CommandTypeBully:
		return bully.New(cfg)
	case CommandTypePostman:
		return postman.New(cfg)
	case CommandTypeProxy:
		return proxy.New(cfg)
	default:
		return nil
	}
}
