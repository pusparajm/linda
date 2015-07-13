package bully

import (
	"github.com/kpashka/linda/config"
)

// Bully - simple command that react with static phrase to expression
type Bully struct {
	cfg config.Command
}

// Create new Bully instance
func New(cfg config.Command) *Bully {
	c := new(Bully)
	c.cfg = cfg
	return c
}

// Return response
func (c *Bully) Run(params []string) (string, error) {
	return c.cfg.Response, nil
}
