package help

import (
	"fmt"
	"strings"

	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

// Help command, nothing valuable here
type Help struct {
	id  string
	cfg config.Command
}

// Create new Help instance
func New(id string, cfg config.Command) *Help {
	c := new(Help)
	c.id = id
	c.cfg = cfg
	return c
}

// Return response
func (c *Help) Run(user *commons.User, params []string) (string, error) {
	response := fmt.Sprintf("\nAvailable commands:\n%s\n", strings.Join(params, "\n"))
	return response, nil
}
