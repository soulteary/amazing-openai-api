package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/internal/fn"
	"github.com/soulteary/amazing-openai-api/internal/network"
)

const (
	HeaderAuthKey = "api-key"
	HeaderAPIVer  = "api-version"
)

func ProxyWithConverter(requestConverter RequestConverter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, x-requested-with")
			c.Status(200)
			return
		}
		Proxy(c, requestConverter)
	}
}

var maskURL = regexp.MustCompile(`key=.+`)

// Proxy Gemini
func Proxy(c *gin.Context, requestConverter RequestConverter) {
	var body []byte
	director := func(req *http.Request) {
		if req.Body == nil {
			network.SendError(c, errors.New("request body is empty"))
			return
		}
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var openaiPayload define.OpenAI_Payload
		err := json.Unmarshal(body, &openaiPayload)
		if err != nil {
			network.SendError(c, errors.Wrap(err, "parse payload error"))
			return
		}

		model := strings.TrimSpace(openaiPayload.Model)
		if model == "" {
			model = DEFAULT_GEMINI_MODEL
		}

		config, ok := ModelConfig[model]
		if ok {
			fmt.Println("rewrite model ", model, "to", config.Model)
			openaiPayload.Model = config.Model
		}

		var payload GoogleGeminiPayload
		for _, data := range openaiPayload.Messages {
			var message GeminiPayloadContents
			if strings.ToLower(data.Role) == "user" {
				message.Role = "user"
			} else {
				message.Role = "model"
			}
			message.Parts = append(message.Parts, GeminiPayloadParts{
				Text: strings.TrimSpace(data.Content),
			})
			payload.Contents = append(payload.Contents, message)
		}

		// set default safety settings
		var safetySettings []GeminiSafetySettings
		safetySettings = append(safetySettings, GeminiSafetySettings{
			Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
			Threshold: "BLOCK_ONLY_HIGH",
		})
		payload.SafetySettings = safetySettings

		// set default generation config
		payload.GenerationConfig.StopSequences = []string{"Title"}
		payload.GenerationConfig.Temperature = openaiPayload.Temperature
		payload.GenerationConfig.MaxOutputTokens = openaiPayload.MaxTokens
		payload.GenerationConfig.TopP = openaiPayload.TopP
		// payload.GenerationConfig.TopK = openaiPayload.TopK

		// get deployment from request
		deployment, err := GetDeploymentByModel(model)
		if err != nil {
			network.SendError(c, err)
			return
		}
		// get auth token from header or deployemnt config
		token := deployment.Key
		if token == "" {
			rawToken := req.Header.Get("Authorization")
			token = strings.TrimPrefix(rawToken, "Bearer ")
		}
		if token == "" {
			network.SendError(c, errors.New("token is empty"))
			return
		}
		req.Header.Set("Authorization", token)

		repack, err := json.Marshal(payload)
		if err != nil {
			network.SendError(c, errors.Wrap(err, "repack payload error"))
			return
		}

		originURL := req.URL.String()
		req, err = requestConverter.Convert(req, deployment, repack)
		if err != nil {
			network.SendError(c, errors.Wrap(err, "convert request error"))
			return
		}

		log.Printf("proxying request [%s] %s -> %s", model, originURL, maskURL.ReplaceAllString(req.URL.String(), "key=******"))
	}

	proxy := &httputil.ReverseProxy{Director: director}
	transport, err := network.NewProxyFromEnv(
		fn.GetStringOrDefaultFromEnv("ENV_GEMINI_SOCKS_PROXY", ""),
		fn.GetStringOrDefaultFromEnv("ENV_GEMINI_HTTP_PROXY", ""),
	)
	if err != nil {
		network.SendError(c, errors.Wrap(err, "get proxy error"))
		return
	}
	if transport != nil {
		proxy.Transport = transport
	}

	proxy.ServeHTTP(c.Writer, c.Request)

	// issue: https://github.com/Chanzhaoyu/chatgpt-web/issues/831
	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte{'\n'}); err != nil {
			log.Printf("rewrite response error: %v", err)
		}
	}

	if c.Writer.Status() != 200 {
		log.Printf("encountering error with body: %s", string(body))
	}
}

func GetDeploymentByModel(model string) (*define.ModelConfig, error) {
	deploymentConfig, exist := ModelConfig[model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("deployment config for %s not found", model))
	}
	return &deploymentConfig, nil
}
