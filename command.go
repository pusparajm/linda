package main

import (
	"strings"

	"github.com/nlopes/slack"
)

const (
	CommandTypeArthur  = "Arthur"
	CommandTypeFeed    = "Feed"
	CommandTypeTommi   = "Tommi"
	CommandTypeWeather = "Weather"
)

type Command interface {
	Trigger(d *DumbSlut, msg *slack.MessageEvent) bool
	Respond(d *DumbSlut, msg *slack.MessageEvent)
}

func NewCommand(cfg CmdConfig) Command {
	switch cfg.Type {
	case CommandTypeArthur:
		return NewArthurCommand(cfg)
	case CommandTypeFeed:
		return NewFeedCommand(cfg)
	case CommandTypeTommi:
		return NewTommiCommand(cfg)
	case CommandTypeWeather:
		return NewWeatherCommand(cfg)
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
