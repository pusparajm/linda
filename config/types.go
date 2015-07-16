package config

const (
	// Execute only first matched command
	ExecutionModeFirst = "first"

	// Execute all matched commands
	ExecutionModeAll = "all"
)

// Bot configuration
type Bot struct {
	Adapter  Adapter
	Commands map[string]Command
	Params   Params
}

// Adapter configuration
type Adapter struct {
	Type    string
	Token   string
	Channel string
}

// Command configuration
type Command struct {
	Type        string
	Description string
	Expression  string
	Filters     [][]string
	Params      map[string]string
	Response    string
	Url         string
}

// Bot parameters
type Params struct {
	ExecutionMode string `toml:"execution_mode"`
	LogLevel      string `toml:"log_level"`
	Shy           bool   `toml:"shy_mode"`
	Tracing       bool   `toml:"trace_errors"`

	Nicknames map[string]string
	Salutes   struct {
		Greeting   string
		Farewell   string
		UserActive string
		UserAway   string
	}
}
