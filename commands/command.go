package commands

import (
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"

	"github.com/kpashka/linda/commands/artist"
	"github.com/kpashka/linda/commands/bully"
	"github.com/kpashka/linda/commands/copycat"
	"github.com/kpashka/linda/commands/help"
	"github.com/kpashka/linda/commands/postman"
	"github.com/kpashka/linda/commands/proxy"
)

const (
	TypeArtist  = "artist"
	TypeBully   = "bully"
	TypeCopycat = "copycat"
	TypeHelp    = "help"
	TypePostman = "postman"
	TypeProxy   = "proxy"
)

// Command interface
type Command interface {
	Run(user *commons.User, params []string) (string, error)
}

// Creates new Command instance
func New(id string, cfg config.Command) Command {
	switch cfg.Type {
	case TypeArtist:
		return artist.New(id, cfg)
	case TypeBully:
		return bully.New(id, cfg)
	case TypeCopycat:
		return copycat.New(id, cfg)
	case TypeHelp:
		return help.New(id, cfg)
	case TypePostman:
		return postman.New(id, cfg)
	case TypeProxy:
		return proxy.New(id, cfg)
	default:
		return nil
	}
}
