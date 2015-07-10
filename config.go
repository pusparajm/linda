package main

type Config struct {
	ApiToken string `json:"api_token,omitempty"`
	Channel  string `json:"channel,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
	LogLevel string `json:"level,omitempty"`
	Shy      bool   `json:"shy,omitempty"`
	Salutes  struct {
		Greeting   string `json:"greeting,omitempty"`
		Farewell   string `json:"farewell,omitempty"`
		UserActive string `json:"user_active,omitempty"`
		UserAway   string `json:"user_away,omitempty"`
	} `json:"salutes,omitempty"`
	Nicknames map[string]string `json:"nicknames,omitempty"`
	Commands  []CmdConfig       `json:"commands,omitempty"`
}

type CmdConfig struct {
	Type     string   `json:"type,omitempty"`
	Letters  string   `json:"letters,omitempty"`
	Tokens   []string `json:"tokens,omitempty"`
	Response string   `json:"response,omitempty"`
	ApiUrl   string   `json:"api_url,omitempty"`
}
