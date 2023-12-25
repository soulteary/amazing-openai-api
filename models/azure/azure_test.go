package azure_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/soulteary/amazing-openai-api/models/azure"
)

// Helper function to set environment variables for testing.
func setEnv(envMap map[string]string) error {
	for key, value := range envMap {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper function to unset environment variables after testing.
func unsetEnv(keys []string) {
	for _, key := range keys {
		os.Unsetenv(key)
	}
}

// TestInitMissingEndpoint tests if Init returns an error when ENV_AZURE_ENDPOINT is missing.
func TestInitMissingEndpoint(t *testing.T) {
	unsetEnv([]string{azure.ENV_AZURE_ENDPOINT}) // Ensure the environment variable is not set.

	err := azure.Init()
	if err == nil || err.Error() != fmt.Sprintf("missing environment variable %s", azure.ENV_AZURE_ENDPOINT) {
		t.Errorf("Expected missing endpoint error, got %v", err)
	}
}

// TestInitInvalidEndpoint tests if Init returns an error when ENV_AZURE_ENDPOINT is invalid.
func TestInitInvalidEndpoint(t *testing.T) {
	envMap := map[string]string{
		azure.ENV_AZURE_ENDPOINT: "http://invalid-endpoint", // Invalid schema or format
	}
	setEnv(envMap)
	defer unsetEnv([]string{azure.ENV_AZURE_ENDPOINT})

	err := azure.Init()
	if err == nil || err.Error() != fmt.Sprintf("invalid environment variable %s", azure.ENV_AZURE_ENDPOINT) {
		t.Errorf("Expected invalid endpoint error, got %v", err)
	}
}

// TestInitUnsupportedVersion tests if Init sets the default version when an unsupported version is passed.
func TestInitUnsupportedVersion(t *testing.T) {
	envMap := map[string]string{
		azure.ENV_AZURE_ENDPOINT: "https://valid-endpoint.openai.azure.com/",
		azure.ENV_AZURE_API_VER:  "unsupported-version",
	}
	setEnv(envMap)
	defer unsetEnv([]string{azure.ENV_AZURE_ENDPOINT, azure.ENV_AZURE_API_VER})

	err := azure.Init()
	if err != nil {
		t.Fatalf("Init failed with error: %v", err)
	}
	if azure.ModelConfig[azure.DEFAULT_AZURE_MODEL].Version != azure.DEFAULT_AZURE_API_VER {
		t.Errorf("Expected version to be set to default, got %v", azure.ModelConfig[azure.DEFAULT_AZURE_MODEL].Version)
	}
}

// TestInitSuccess tests if Init successfully initializes ModelConfig with the right values.
func TestInitSuccess(t *testing.T) {
	envMap := map[string]string{
		azure.ENV_AZURE_ENDPOINT:    "https://valid-endpoint.openai.azure.com/",
		azure.ENV_AZURE_API_VER:     "2023-03-15-preview",
		azure.ENV_AZURE_API_KEY:     "test-api-key",
		azure.ENV_AZURE_MODEL:       "test-model",
		azure.ENV_AZURE_MODEL_ALIAS: "alias1:test-model-alias",
	}
	setEnv(envMap)
	defer unsetEnv([]string{azure.ENV_AZURE_ENDPOINT, azure.ENV_AZURE_API_VER, azure.ENV_AZURE_API_KEY, azure.ENV_AZURE_MODEL, azure.ENV_AZURE_MODEL_ALIAS})

	err := azure.Init()
	if err != nil {
		t.Fatalf("Init failed with error: %v", err)
	}

	modelConfig, ok := azure.ModelConfig["test-model"]
	if !ok {
		t.Fatalf("Model 'test-model' not found in ModelConfig")
	}

	if modelConfig.Endpoint != "https://valid-endpoint.openai.azure.com/" {
		t.Errorf("Expected endpoint to match, got %v", modelConfig.Endpoint)
	}
	if modelConfig.Version != "2023-03-15-preview" {
		t.Errorf("Expected API version to match, got %v", modelConfig.Version)
	}

	if modelConfig.Model != "test-model" { // The alias should override the original model name.
		t.Errorf("Expected model to use alias, got %v", modelConfig.Model)
	}

	config, ok := azure.ModelConfig["alias1"]
	if !ok {
		t.Fatalf("Model 'alias1' not found in ModelConfig")
	}
	if config.Alias != "test-model-alias" {
		t.Errorf("Expected model to match, got %v", config.Alias)
	}
}
