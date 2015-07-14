package snitch

import (
	"fmt"
	"strings"

	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
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
func (c *Snitch) Run(user *commons.User, params []string) (string, error) {
	response := fmt.Sprintf("\nAvailable commands:\n%s", strings.Join(params, "\n"))
	return response, nil
}
