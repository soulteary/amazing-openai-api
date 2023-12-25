package define

import "net/url"

type ModelConfig struct {
	Name     string `yaml:"name" json:"name"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Model    string `yaml:"model" json:"model"`
	Version  string `yaml:"version" json:"version"`
	Key      string `yaml:"key" json:"key"`
	URL      *url.URL
	Alias    string
}

type ModelAlias [][]string

// openai api payload
type RequestData struct {
	MaxTokens       int       `json:"max_tokens"`
	Model           string    `json:"model"`
	Temperature     float64   `json:"temperature"`
	TopP            float64   `json:"top_p"`
	PresencePenalty float64   `json:"presence_penalty"`
	Messages        []Message `json:"messages"`
	Stream          bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
