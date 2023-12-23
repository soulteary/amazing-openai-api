package gemini

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/soulteary/amazing-openai-api/internal/define"
)

type RequestConverter interface {
	Name() string
	Convert(req *http.Request, config *define.ModelConfig, payload []byte) (*http.Request, error)
}

type StripPrefixConverter struct {
	Prefix string
}

func (c *StripPrefixConverter) Name() string {
	return "StripPrefix"
}

func (c *StripPrefixConverter) Convert(req *http.Request, config *define.ModelConfig, payload []byte) (*http.Request, error) {
	req.Host = config.URL.Host
	req.URL.Scheme = config.URL.Scheme
	req.URL.Host = config.URL.Host
	req.URL.Path = fmt.Sprintf("%s/models/%s:generateContent", config.Version, config.Model)
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add("key", config.Key)
	req.URL.RawQuery = query.Encode()
	req.Body = io.NopCloser(bytes.NewBuffer(payload))
	req.ContentLength = int64(len(payload))
	return req, nil
}

func NewStripPrefixConverter(prefix string) *StripPrefixConverter {
	return &StripPrefixConverter{
		Prefix: prefix,
	}
}
