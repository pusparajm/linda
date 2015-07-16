package copycat

import (
	"fmt"

	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

// Help command, nothing valuable here
type Copycat struct {
	id  string
	cfg config.Command
}

// Create new Copycat instance
func New(id string, cfg config.Command) *Copycat {
	c := new(Copycat)
	c.id = id
	c.cfg = cfg
	return c
}

// Return response
func (c *Copycat) Run(user *commons.User, params []string) (string, error) {
	response := fmt.Sprintf("@%s: %s", user.Username, params[1])
	return response, nil
}
