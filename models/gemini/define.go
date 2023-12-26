package gemini

const (
	ENV_GEMINI_ENDPOINT    = "GEMINI_ENDPOINT"
	ENV_GEMINI_API_VER     = "GEMINI_API_VER"
	ENV_GEMINI_MODEL_ALIAS = "GEMINI_MODEL_ALIAS"
	ENV_GEMINI_API_KEY     = "GEMINI_API_KEY"
	ENV_GEMINI_MODEL       = "GEMINI_MODEL"

	ENV_GEMINI_HTTP_PROXY  = "GEMINI_HTTP_PROXY"
	ENV_GEMINI_SOCKS_PROXY = "GEMINI_SOCKS_PROXY"
)

const (
	DEFAULT_REST_API_VERSION_SHIM = "/v1"
	DEFAULT_REST_API_VERSION      = "/v1beta"
	DEFAULT_REST_API_ENTRYPOINT   = "https://generativelanguage.googleapis.com"
)

const (
	DEFAULT_GEMINI_API_VER = DEFAULT_REST_API_VERSION
	DEFAULT_GEMINI_MODEL   = "gemini-pro"
)

type OpenAIPayloadMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIPayload struct {
	MaxTokens       int                     `json:"max_tokens"`
	Model           string                  `json:"model"`
	Temperature     float64                 `json:"temperature"`
	TopP            float64                 `json:"top_p"`
	PresencePenalty float64                 `json:"presence_penalty"`
	Messages        []OpenAIPayloadMessages `json:"messages"`
	Stream          bool                    `json:"stream"`
}

type GoogleGeminiPayload struct {
	Contents         []GeminiPayloadContents `json:"contents"`
	SafetySettings   []GeminiSafetySettings  `json:"safetySettings"`
	GenerationConfig GeminiGenerationConfig  `json:"generationConfig"`
}

type GeminiSafetySettings struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiGenerationConfig struct {
	StopSequences   []string `json:"stopSequences"`
	Temperature     float64  `json:"temperature,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	TopP            float64  `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
}

// gemini response
type GeminiSafetyRatings struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type PromptFeedback struct {
	SafetyRatings []GeminiSafetyRatings `json:"safetyRatings"`
}

type GeminiPayloadParts struct {
	Text string `json:"text"`
}

type GeminiPayloadContents struct {
	Parts []GeminiPayloadParts `json:"parts"`
	Role  string               `json:"role"`
}

type GeminiCandidates struct {
	Content       GeminiPayloadContents `json:"content"`
	FinishReason  string                `json:"finishReason"`
	Index         int                   `json:"index"`
	SafetyRatings []GeminiSafetyRatings `json:"safetyRatings"`
}

type GeminiResponse struct {
	Candidates     []GeminiCandidates `json:"candidates"`
	PromptFeedback PromptFeedback     `json:"promptFeedback"`
}
