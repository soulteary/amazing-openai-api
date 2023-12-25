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
type OpenAI_Payload struct {
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

type OpenAI_Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAI_Choices struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type OpeAI_Response struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int              `json:"created"`
	Model   string           `json:"model"`
	Usage   OpenAI_Usage     `json:"usage"`
	Choices []OpenAI_Choices `json:"choices"`
	// openai extra fields
	SystemFingerprint string `json:"system_fingerprint"`
}
