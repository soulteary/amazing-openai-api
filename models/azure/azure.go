package azure

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/internal/fn"
)

var (
	ModelConfig = map[string]define.ModelConfig{}
)

func Init() (err error) {
	var modelConfig define.ModelConfig

	// azure openai api endpoint
	endpoint := fn.GetStringOrDefaultFromEnv(ENV_AZURE_ENDPOINT, "")
	if endpoint == "" {
		return fmt.Errorf("missing environment variable %s", ENV_AZURE_ENDPOINT)
	}
	// Use a URL starting with `https://` and ending with `.openai.azure.com/`
	if !(strings.HasPrefix(endpoint, "https://") && strings.HasSuffix(endpoint, ".openai.azure.com/")) {
		return fmt.Errorf("invalid environment variable %s", ENV_AZURE_ENDPOINT)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse endpoint error: %w", err)
	}
	modelConfig.URL = u
	modelConfig.Endpoint = endpoint

	// azure openai api version
	apiVersion := fn.GetStringOrDefaultFromEnv(ENV_AZURE_API_VER, DEFAULT_AZURE_API_VER)
	// azure openai api versions supported
	// https://learn.microsoft.com/en-us/azure/ai-services/openai/reference
	if apiVersion != "2022-12-01" &&
		apiVersion != "2023-03-15-preview" &&
		apiVersion != "2023-05-15" &&
		apiVersion != "2023-06-01-preview" &&
		apiVersion != "2023-07-01-preview" &&
		apiVersion != "2023-08-01-preview" &&
		apiVersion != "2023-09-01-preview" {
		apiVersion = DEFAULT_AZURE_API_VER
	}
	modelConfig.Version = apiVersion

	// azure openai api key, allow override by request header
	apikey := fn.GetStringOrDefaultFromEnv(ENV_AZURE_API_KEY, "")
	modelConfig.Key = apikey

	// azure openai api model
	model := fn.GetStringOrDefaultFromEnv(ENV_AZURE_MODEL, DEFAULT_AZURE_MODEL)
	if model == "" {
		model = DEFAULT_AZURE_MODEL
	}
	modelConfig.Model = model

	modelConfig.Vision = fn.GetBoolOrDefaultFromEnv(ENV_AZURE_VISION, false)

	ModelConfig[model] = modelConfig

	// azure openai api model alias
	alias := fn.ExtractModelAlias(fn.GetStringOrDefaultFromEnv(ENV_AZURE_MODEL_ALIAS, ""))
	for _, pair := range alias {
		modelConfig.Alias = pair[1]
		ModelConfig[pair[0]] = modelConfig
	}
	return nil
}
