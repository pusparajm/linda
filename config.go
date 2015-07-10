package main

type Config struct {
	ApiToken    string      `json:"api_token,omitempty"`
	Channel     string      `json:"channel,omitempty"`
	Debug       bool        `json:"debug,omitempty"`
	LogLevel    string      `json:"level,omitempty"`
	MsgGreeting string      `json:"greeting,omitempty"`
	MsgFarewell string      `json:"farewell,omitempty"`
	Shy         bool        `json:"shy,omitempty"`
	Commands    []CmdConfig `json:"commands,omitempty"`
}

type CmdConfig struct {
	Type     string   `json:"type,omitempty"`
	Letters  string   `json:"letters,omitempty"`
	Tokens   []string `json:"tokens,omitempty"`
	Response string   `json:"response,omitempty"`
	ApiUrl   string   `json:"api_url,omitempty"`
}
