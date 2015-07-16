package bully

import (
	"github.com/kpashka/linda/commons"
	"github.com/kpashka/linda/config"
)

// Bully - simple command that react with static phrase to expression
type Bully struct {
	id  string
	cfg config.Command
}

// Create new Bully instance
func New(id string, cfg config.Command) *Bully {
	c := new(Bully)
	c.id = id
	c.cfg = cfg
	return c
}

// Return response
func (c *Bully) Run(user *commons.User, params []string) (string, error) {
	return c.cfg.Response, nil
}
