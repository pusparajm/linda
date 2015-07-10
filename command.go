package main

import (
	"strings"

	"github.com/nlopes/slack"
)

const (
	CommandTypeTommi   = "Tommi"
	CommandTypeArthur  = "Arthur"
	CommandTypeWeather = "Weather"
)

type Command interface {
	Trigger(d *DumbSlut, msg *slack.MessageEvent) bool
	Respond(d *DumbSlut, msg *slack.MessageEvent)
}

func NewCommand(cfg CmdConfig) Command {
	switch cfg.Type {
	case CommandTypeTommi:
		return NewTommiCommand(cfg)
	case CommandTypeArthur:
		return NewArthurCommand(cfg)
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
		if strings.Contains(toLower, token) {
			return token, true
		}
	}

	return "", false
}
