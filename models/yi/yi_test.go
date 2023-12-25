package yi_test

import (
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/soulteary/amazing-openai-api/models/yi"
)

func TestInit(t *testing.T) {
	// Load environment variables from a .env file for testing if needed
	// godotenv.Load("../path/to/your/.env.testing")

	t.Run("it should handle missing endpoint error", func(t *testing.T) {
		// Clear environment variable for endpoint
		os.Unsetenv(yi.ENV_YI_ENDPOINT)

		err := yi.Init()
		if err == nil || err.Error() != "missing environment variable "+yi.ENV_YI_ENDPOINT {
			t.Errorf("Expected missing endpoint environment variable error, got %v", err)
		}
	})

	t.Run("it should parse and assign endpoint successfully", func(t *testing.T) {
		expectedURL := "https://example.com/api"
		os.Setenv(yi.ENV_YI_ENDPOINT, expectedURL)

		err := yi.Init()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		modelCfg, exists := yi.ModelConfig[yi.DEFAULT_YI_MODEL]
		if !exists {
			t.Fatal("Model config does not exist after Init")
		}

		parsedURL, _ := url.Parse(expectedURL)
		if !reflect.DeepEqual(modelCfg.URL, parsedURL) {
			t.Errorf("Expected URL to be parsed correctly, got %+v", modelCfg.URL)
		}
	})

	// ... Additional tests for version, api key, model, aliasing, etc.

	// Reset environment variables after testing
	t.Cleanup(func() {
		os.Unsetenv(yi.ENV_YI_ENDPOINT)
		// ... Unset other environment variables used during the tests
	})
}
