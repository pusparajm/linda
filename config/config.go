package config

// Combined config for better readability
type Config struct {
	// Backend configuration
	Backend struct {
		// Backend type - Slack, Telegram etc.
		Type string `json:"type,omitempty"`

		// Bot API token
		Token string `json:"token,omitempty"`

		// Channel to work
		Channel string `json:"channel,omitempty"`
	} `json:"backend,omitempty"`

	// Commands configuration
	Commands []struct {
		// Command type
		Type string `json:"type,omitempty"`

		// Specific name of command type instance
		Name string `json:"name,omitempty"`

		// Command description
		Description string `json:"description,omitempty"`

		// Regular expression to react
		Expression string `json:"expression,omitempty"`

		// Additional parameters
		Params map[string]string `json:"params,omitempty"`

		// Response template
		Response string `json:"response,omitempty"`

		// URL to request
		Url string `json:"url,omitempty"`
	} `json:"commands,omitempty"`

	// Additional bot parameters
	Params struct {
		// Command execution mode - first, all
		ExecutionMode string `json:"execution_mode,omitempty"`

		// Log level for Logrus library
		LogLevel string `json:"log_level,omitempty"`

		// Nicknames map for users
		Nicknames map[string]string `json:"nicknames,omitempty"`

		// Reactions to presence events
		Salutes struct {
			// After bot is connected to channel
			Greeting string `json:"greeting,omitempty"`

			// Before bot is disconnected to channel
			Farewell string `json:"farewell,omitempty"`

			// When some user entered channel
			UserEntered string `json:"user_entered,omitempty"`

			// When some user left channel
			UserLeft string `json:"user_left,omitempty"`
		} `json:"salutes,omitempty"`

		// Salutes override switch
		Shy bool `json:"shy,omitempty"`
	} `json:"params,omitempty"`
}
