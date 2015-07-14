package cyber

import (
	"strings"
)

// Cyber hook
type Cyber struct{}

// New Cyber hook
func New() *Cyber {
	return new(Cyber)
}

// Cyber input
func (hook *Cyber) Fire(input string) string {
	letters := []string{}
	for _, i := range input {
		letters = append(letters, string(i))
	}

	raw := strings.Join(letters, " ")
	return strings.ToUpper(raw)
}
