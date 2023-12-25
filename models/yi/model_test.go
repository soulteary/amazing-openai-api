package yi_test

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/models/yi"
	"github.com/stretchr/testify/assert"
)

const HeaderAPIVer string = "X-API-Version"

func TestStripPrefixConverter_Name(t *testing.T) {
	converter := yi.NewStripPrefixConverter("/api")
	assert.Equal(t, "StripPrefix", converter.Name())
}

func TestStripPrefixConverter_Convert(t *testing.T) {
	prefix := "/api"
	converter := yi.NewStripPrefixConverter(prefix)
	modelConfig := &define.ModelConfig{
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.com",
		},
		Version: "v1",
	}

	req := httptest.NewRequest("GET", "http://localhost"+prefix+"/endpoint?param=value", nil)
	convertedReq, err := converter.Convert(req, modelConfig)

	assert.NoError(t, err)
	assert.NotNil(t, convertedReq)
	assert.Equal(t, "example.com", convertedReq.Host)
	assert.Equal(t, "https", convertedReq.URL.Scheme)
	assert.Equal(t, "example.com", convertedReq.URL.Host)

	// Ensure original path is maintained without the prefix
	assert.Contains(t, convertedReq.URL.Path, prefix)
}
