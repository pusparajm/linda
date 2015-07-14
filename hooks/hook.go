package hooks

import (
	"github.com/kpashka/linda/hooks/cyber"
	"github.com/kpashka/linda/hooks/translit"
)

const (
	TypeCyber    = "cyber"
	TypeTranslit = "translit"
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
	default:
		return nil
	}
}
