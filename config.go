package main

type Config struct {
	ApiToken string `json:"api_token,omitempty"`
	Channel  string `json:"channel,omitempty"`

	Debug    bool   `json:"debug,omitempty"`
	LogLevel string `json:"log_level,omitempty"`

	Shy       bool              `json:"shy,omitempty"`
	Nicknames map[string]string `json:"nicknames,omitempty"`
	Salutes   struct {
		Greeting   string `json:"greeting,omitempty"`
		Farewell   string `json:"farewell,omitempty"`
		UserActive string `json:"user_active,omitempty"`
		UserAway   string `json:"user_away,omitempty"`
	} `json:"salutes,omitempty"`

	Commands []CmdConfig `json:"commands,omitempty"`
}

type CmdConfig struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`

	Letters string   `json:"letters,omitempty"`
	Tokens  []string `json:"tokens,omitempty"`

	Response string `json:"response,omitempty"`
	Url      string `json:"url,omitempty"`
}
