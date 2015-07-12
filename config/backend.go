package config

// Backend configuration
type Backend struct {
	// Backend type - Slack, Telegram etc.
	Type string `json:"type,omitempty"`

	// Bot API token
	Token string `json:"token,omitempty"`

	// Channel to work
	Channel string `json:"channel,omitempty"`
}
