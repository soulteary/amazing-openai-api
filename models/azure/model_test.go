package azure_test

import (
	"net/http"
	"net/url"

	"testing"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/models/azure"
)

func TestStripPrefixConverter_Convert(t *testing.T) {
	prefix := "/api/v1"
	converter := azure.NewStripPrefixConverter(prefix)

	u, _ := url.Parse("https://example.com")

	modelConfig := &define.ModelConfig{
		Model:   "test-model",
		Version: "2023-04-01",
		URL:     u,
	}

	reqURL, _ := url.Parse("http://localhost:8080/api/v1/model/predict")
	req := &http.Request{
		URL:    reqURL,
		Header: http.Header{},
	}

	convertedReq, err := converter.Convert(req, modelConfig)
	if err != nil {
		t.Fatalf("Convert failed with error: %v", err)
	}

	expectedPath := "/openai/deployments/test-model/model/predict"
	if convertedReq.URL.Path != expectedPath {
		t.Errorf("Expected path '%s', but got '%s'", expectedPath, convertedReq.URL.Path)
	}

	if convertedReq.URL.Host != modelConfig.URL.Host {
		t.Errorf("Expected host '%s', but got '%s'", modelConfig.URL.Host, convertedReq.URL.Host)
	}

	if convertedReq.URL.Scheme != modelConfig.URL.Scheme {
		t.Errorf("Expected scheme '%s', but got '%s'", modelConfig.URL.Scheme, convertedReq.URL.Scheme)
	}

	expectedVersion := modelConfig.Version
	queryValues := convertedReq.URL.Query()
	if queryValues.Get(azure.HeaderAPIVer) != expectedVersion {
		t.Errorf("Expected API version query parameter '%s', but got '%s'", expectedVersion, queryValues.Get(azure.HeaderAPIVer))
	}
}
