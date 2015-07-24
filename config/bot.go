package config

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	// Environment variable for platform API token
	ApiTokenEnvVar = "LINDA_API_TOKEN"
)

// Create new bot configuration instance
func New() *Bot {
	bot := new(Bot)
	return bot
}

// Load configuration from specified location
func (bot *Bot) Load(location string) error {
	var bytes []byte
	var err error

	// Load file / url
	if strings.Contains(location, "http://") || strings.Contains(location, "https://") {
		bytes, err = bot.loadFromUrl(location)
	} else {
		bytes, err = ioutil.ReadFile(location)
	}

	if err != nil {
		return err
	}

	// Load config from bytes
	if _, err := toml.Decode(string(bytes), bot); err != nil {
		return err
	}

	// Try to grab API token from env variable if empty
	if bot.Adapter.Token == "" {
		bot.Adapter.Token = os.Getenv(ApiTokenEnvVar)
	}

	return err
}

// Load configuration from URL
func (bot *Bot) loadFromUrl(location string) ([]byte, error) {
	response, err := http.Get(location)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
