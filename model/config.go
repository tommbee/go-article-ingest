package model

// Config object
type Config struct {
	BaseURL      string `json:"baseUrl"`
	BaseElement  string `json:"baseElement"`
	TitleElement string `json:"titleElement"`
	DateElement  string `json:"dateElement"`
	LinkElement  string `json:"linkElement"`
	DateFormat   string `json:"dateFormat"`
}

// Configs - a list of config objects
type Configs struct {
	Configs []Config `json:"sources"`
}
