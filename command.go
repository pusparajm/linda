package main

import (
	"strings"

	"github.com/nlopes/slack"
)

const (
	CommandTypeBully     = "Bully"
	CommandTypeFeed      = "Feed"
	CommandTypeSwastika  = "Swastika"
	CommandTypeTranslate = "Translate"
	CommandTypeWeather   = "Weather"
	CommandTypeWiki      = "Wiki"
)

type Command interface {
	GetName() string
	Trigger(d *DumbSlut, msg *slack.MessageEvent) bool
	Execute(d *DumbSlut, msg *slack.MessageEvent)
}

func NewCommand(cfg CmdConfig) Command {
	switch cfg.Type {
	case CommandTypeBully:
		return NewBullyCommand(cfg)
	case CommandTypeFeed:
		return NewFeedCommand(cfg)
	case CommandTypeSwastika:
		return NewSwastikaCommand(cfg)
	case CommandTypeTranslate:
		return NewTranslateCommand(cfg)
	case CommandTypeWeather:
		return NewWeatherCommand(cfg)
	case CommandTypeWiki:
		return NewWikiCommand(cfg)
	default:
		return nil
	}
}

func containsToken(text string, tokens []string) (string, bool) {
	if len(tokens) == 0 {
		return "", false
	}

	toLower := strings.ToLower(text)
	for _, token := range tokens {
		// Temporary - till I'll learn Armenian
		if strings.Contains(text, token) {
			return token, true
		}

		if strings.Contains(toLower, token) {
			return token, true
		}
	}

	return "", false
}
