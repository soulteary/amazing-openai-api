package yi

import (
	"net/http"

	"github.com/soulteary/amazing-openai-api/internal/define"
)

type RequestConverter interface {
	Name() string
	Convert(req *http.Request, config *define.ModelConfig) (*http.Request, error)
}

type StripPrefixConverter struct {
	Prefix string
}

func (c *StripPrefixConverter) Name() string {
	return "StripPrefix"
}

func (c *StripPrefixConverter) Convert(req *http.Request, config *define.ModelConfig) (*http.Request, error) {
	req.Host = config.URL.Host
	req.URL.Scheme = config.URL.Scheme
	req.URL.Host = config.URL.Host
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add(HeaderAPIVer, config.Version)
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func NewStripPrefixConverter(prefix string) *StripPrefixConverter {
	return &StripPrefixConverter{
		Prefix: prefix,
	}
}
