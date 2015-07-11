package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
	"github.com/kpashka/dumbslut/swastika"
	"github.com/nlopes/slack"
)

const (
	SwastikaMaxLen = 10
)

type SwastikaCommand struct {
	config       CmdConfig
	currentToken string
}

func NewSwastikaCommand(cfg CmdConfig) *SwastikaCommand {
	c := new(SwastikaCommand)
	c.config = cfg
	return c
}

func (c *SwastikaCommand) GetName() string {
	return c.config.Name
}

func (c *SwastikaCommand) Trigger(d *DumbSlut, msg *slack.MessageEvent) bool {
	currentToken, ok := containsToken(msg.Text, c.config.Tokens)
	if !ok {
		return false
	}

	c.currentToken = currentToken
	return true
}

func (c *SwastikaCommand) Execute(d *DumbSlut, msg *slack.MessageEvent) {
	query := strings.TrimSpace(strings.Replace(msg.Text, c.currentToken, "", 1))
	length := utf8.RuneCountInString(query)

	if length < swastika.MinLen || length > SwastikaMaxLen {
		message := fmt.Sprintf("Input string size should be in range [%d, %d]", swastika.MinLen, SwastikaMaxLen)
		log.Error(message)
		d.TalkTo(message, msg.UserId)
		return
	}

	response, err := swastika.New(strings.ToUpper(query))
	if err != nil {
		log.Error(err.Error())
		d.TalkTo(err.Error(), msg.UserId)
	}

	message := fmt.Sprintf("\n```%s```", response)
	d.TalkTo(message, msg.UserId)
}
