package config

const (
	ExecutionModeFirst = "first"
	ExecutionModeAll   = "all"
)

// Bot parameters
type Params struct {
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

	// Tracing - send errors back to chat
	Tracing bool `json:"tracing,omitempty"`
}
