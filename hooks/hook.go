package hooks

import (
	"github.com/kpashka/linda/hooks/cyber"
	"github.com/kpashka/linda/hooks/translit"
	"github.com/kpashka/linda/hooks/zalgo"
)

const (
	TypeCyber    = "cyber"
	TypeTranslit = "translit"
	TypeZalgo    = "zalgo"
)

// Hook interface
type Hook interface {
	Fire(input string) string
}

// Creates new Hook instance
func New(hook string) Hook {
	switch hook {
	case TypeCyber:
		return cyber.New()
	case TypeTranslit:
		return translit.New()
	case TypeZalgo:
		return zalgo.New()
	default:
		return nil
	}
}
