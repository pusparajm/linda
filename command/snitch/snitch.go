package snitch

import (
	"fmt"
	"strings"

	"github.com/kpashka/dumbslut/config"
)

// Help command, nothing valuable here
type Snitch struct {
	cfg config.Command
}

// Create new Snitch instance
func New(cfg config.Command) *Snitch {
	c := new(Snitch)
	c.cfg = cfg
	return c
}

// Return response
func (c *Snitch) Run(params []string) (string, error) {
	response := fmt.Sprintf("Available commands:\n%s", strings.Join(params, "\n"))
	return response, nil
}
