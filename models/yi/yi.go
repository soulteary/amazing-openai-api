package yi

import (
	"fmt"
	"net/url"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/internal/fn"
)

var (
	ModelConfig = map[string]define.ModelConfig{}
)

func Init() (err error) {
	var modelConfig define.ModelConfig

	// yi api endpoint
	endpoint := fn.GetStringOrDefaultFromEnv(ENV_YI_ENDPOINT, "")
	if endpoint == "" {
		return fmt.Errorf("missing environment variable %s", ENV_YI_ENDPOINT)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse endpoint error: %w", err)
	}
	modelConfig.URL = u
	modelConfig.Endpoint = endpoint

	// yi api version
	apiVersion := fn.GetStringOrDefaultFromEnv(ENV_YI_API_VER, DEFAULT_YI_API_VER)
	if apiVersion == "" {
		apiVersion = DEFAULT_YI_API_VER
	}
	modelConfig.Version = apiVersion

	// yi api key, allow override by request header
	apikey := fn.GetStringOrDefaultFromEnv(ENV_YI_API_KEY, "")
	modelConfig.Key = apikey

	// yi api model
	model := fn.GetStringOrDefaultFromEnv(ENV_YI_MODEL, DEFAULT_YI_MODEL)
	if model == "" {
		model = DEFAULT_YI_MODEL
	}
	modelConfig.Model = model

	ModelConfig[model] = modelConfig

	// yi api model alias
	alias := fn.ExtractModelAlias(fn.GetStringOrDefaultFromEnv(ENV_YI_MODEL_ALIAS, ""))
	for _, pair := range alias {
		if model == pair[0] {
			modelConfig.Model = pair[1]
		}
		ModelConfig[pair[0]] = modelConfig
	}
	return nil
}
