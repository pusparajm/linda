package main

import (
	"github.com/nlopes/slack"
)

type TommiCommand struct {
	config CmdConfig
}

func NewTommiCommand(cfg CmdConfig) *TommiCommand {
	c := new(TommiCommand)
	c.config = cfg
	return c
}

func (c *TommiCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	_, ok := containsToken(msg.Text, c.config.Tokens)
	return ok
}

func (c *TommiCommand) Respond(d *DumbSlut, msg *slack.MessageEvent) {
	d.Talk(c.config.Response)
}
