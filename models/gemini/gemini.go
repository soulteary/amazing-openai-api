package gemini

import (
	"fmt"
	"net/url"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/internal/fn"
)

// refs: https://ai.google.dev/models/gemini?hl=zh-cn
var (
	ModelConfig = map[string]define.ModelConfig{}
)

func Init() (err error) {
	var modelConfig define.ModelConfig

	// gemini openai api endpoint
	endpoint := fn.GetStringOrDefaultFromEnv(ENV_GEMINI_ENDPOINT, DEFAULT_REST_API_ENTRYPOINT)
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse endpoint error: %w", err)
	}
	modelConfig.URL = u
	modelConfig.Endpoint = endpoint

	// gemini openai api version
	apiVersion := fn.GetStringOrDefaultFromEnv(ENV_GEMINI_API_VER, DEFAULT_GEMINI_API_VER)
	// google api versions supported
	// https://ai.google.dev/docs/api_versions?hl=zh-cn
	if apiVersion != "v1" && apiVersion != "v1beta" {
		apiVersion = DEFAULT_GEMINI_API_VER
	} else {
		apiVersion = "/" + apiVersion
	}
	modelConfig.Version = apiVersion

	// gemini openai api key, allow override by request header
	apikey := fn.GetStringOrDefaultFromEnv(ENV_GEMINI_API_KEY, "")
	modelConfig.Key = apikey

	// gemini openai api model
	model := fn.GetStringOrDefaultFromEnv(ENV_GEMINI_MODEL, DEFAULT_GEMINI_MODEL)
	if model == "" {
		model = DEFAULT_GEMINI_MODEL
	}
	modelConfig.Model = model

	ModelConfig[model] = modelConfig

	// gemini openai api model alias
	alias := fn.ExtractModelAlias(fn.GetStringOrDefaultFromEnv(ENV_GEMINI_MODEL_ALIAS, ""))
	for _, pair := range alias {
		if model == pair[0] {
			modelConfig.Model = pair[1]
		}
		ModelConfig[pair[0]] = modelConfig
	}
	return nil
}
