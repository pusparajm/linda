package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type BullyCommand struct {
	config CmdConfig
}

func NewBullyCommand(cfg CmdConfig) *BullyCommand {
	c := new(BullyCommand)
	c.config = cfg
	return c
}

func (c *BullyCommand) GetName() string {
	return c.config.Name
}

func (c *BullyCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	if strings.ContainsAny(msg.Text, c.config.Letters) {
		return true
	}

	_, ok := containsToken(msg.Text, c.config.Tokens)
	return ok
}

func (c *BullyCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	d.Talk(c.config.Response)
}
