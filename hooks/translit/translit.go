package translit

import (
	"github.com/fiam/gounidecode/unidecode"
)

// Translit hook
type Translit struct{}

// New Translit hook
func New() *Translit {
	return new(Translit)
}

// Transliterate input
func (hook *Translit) Fire(input string) string {
	return unidecode.Unidecode(input)
}
