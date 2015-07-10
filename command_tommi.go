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

func (c *TommiCommand) GetType() string {
	return CommandTypeTommi
}

func (c *TommiCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	_, ok := containsToken(msg.Text, c.config.Tokens)
	return ok
}

func (c *TommiCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	d.Talk(c.config.Response)
}
