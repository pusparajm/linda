package config

// Command configuration
type Command struct {
	// Command type
	Type string `json:"type,omitempty"`

	// Specific name of command type instance
	Name string `json:"name,omitempty"`

	// Command description
	Description string `json:"description,omitempty"`

	// Regular expression to react
	Expression string `json:"expression,omitempty"`

	// Hooks list applied to each param
	Hooks map[string][]int `json:"hooks,omitempty"`

	// Additional parameters
	Params map[string]string `json:"params,omitempty"`

	// Response template
	Response string `json:"response,omitempty"`

	// URL to request
	Url string `json:"url,omitempty"`
}
