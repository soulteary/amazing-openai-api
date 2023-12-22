package define

import "net/url"

type ModelConfig struct {
	Name     string `yaml:"name" json:"name"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Model    string `yaml:"model" json:"model"`
	Version  string `yaml:"version" json:"version"`
	Key      string `yaml:"key" json:"key"`
	URL      *url.URL
}

type ModelAlias [][]string
