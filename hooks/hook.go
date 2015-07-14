package hooks

import (
	"github.com/kpashka/linda/hooks/translit"
)

const (
	TypeTranslit = "translit"
)

// Hook interface
type Hook interface {
	Fire(input string) string
}

// Creates new Hook instance
func New(hook string) Hook {
	switch hook {
	case TypeTranslit:
		return translit.New()
	default:
		return nil
	}
}
