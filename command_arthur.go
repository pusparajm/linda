package main

import (
	"strings"

	"github.com/nlopes/slack"
)

type ArthurCommand struct {
	config CmdConfig
}

func NewArthurCommand(cfg CmdConfig) *ArthurCommand {
	c := new(ArthurCommand)
	c.config = cfg
	return c
}

func (c *ArthurCommand) GetType() string {
	return CommandTypeArthur
}

func (c *ArthurCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	if strings.ContainsAny(msg.Text, c.config.Letters) {
		return true
	}

	return false
}

func (c *ArthurCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	d.Talk(c.config.Response)
}
