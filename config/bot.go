package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Bot configuration
type Bot struct {
	// Backend configuration
	Backend Backend `json:"backend,omitempty"`

	// Commands configuration
	Commands []Command `json:"commands,omitempty"`

	// Additional bot parameters
	Params Params `json:"params,omitempty"`
}

// Create new bot configuration instance
func New() *Bot {
	bot := new(Bot)
	return bot
}

// Load configuration from specified location
func (bot *Bot) Load(location string) error {
	var err error
	if strings.Contains(location, "http://") || strings.Contains(location, "https://") {
		err = bot.loadFromUrl(location)
	} else {
		err = bot.loadFromFile(location)
	}

	return err
}

// Load configuration from file
func (bot *Bot) loadFromFile(location string) error {
	// Read config file
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return err
	}

	// Create config instance
	err = json.Unmarshal(bytes, bot)
	if err != nil {
		return err
	}

	return nil
}

// Load configuration from URL
func (bot *Bot) loadFromUrl(location string) error {
	response, err := http.Get(location)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, bot)
	if err != nil {
		return err
	}

	return nil
}
