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
	"strconv"
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

func parseRequestBody(reqBody io.ReadCloser) (openaiPayload define.OpenAI_Payload, err error) {
	if reqBody == nil {
		err = errors.New("request body is empty")
		return openaiPayload, err
	}
	body, _ := io.ReadAll(reqBody)

	// req.Body = io.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, &openaiPayload)
	return openaiPayload, err
}

func GetModelNameAndConfig(openaiPayload define.OpenAI_Payload) (string, define.ModelConfig, bool) {
	model := strings.TrimSpace(openaiPayload.Model)
	if model == "" {
		model = DEFAULT_GEMINI_MODEL
	}

	config, ok := ModelConfig[model]
	if ok {
		fmt.Println("rewrite model ", model, "to", config.Model)
		openaiPayload.Model = config.Model
	}
	return model, config, ok
}

func getDirector(req *http.Request, body []byte, c *gin.Context, requestConverter RequestConverter, openaiPayload define.OpenAI_Payload, model string) func(req *http.Request) {
	return func(req *http.Request) {
		// req.Body = io.NopCloser(bytes.NewBuffer(body))

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
}

// Proxy Gemini
func Proxy(c *gin.Context, requestConverter RequestConverter) {
	var body []byte

	openaiPayload, err := parseRequestBody(c.Request.Body)
	if err != nil {
		network.SendError(c, err)
		return
	}

	model, config, ok := GetModelNameAndConfig(openaiPayload)
	if ok {
		fmt.Println("rewrite model ", model, "to", config.Model)
		openaiPayload.Model = config.Model
	}

	fmt.Println(openaiPayload)

	proxy := &httputil.ReverseProxy{Director: getDirector(c.Request, body, c, requestConverter, openaiPayload, model)}
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

	proxy.ModifyResponse = func(response *http.Response) error {
		if response.StatusCode == http.StatusOK {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			var responsePayload GoogleGeminiPayload
			err = json.Unmarshal(body, &responsePayload)
			if err != nil {
				return err
			}

			var openaiResponse define.OpeAI_Response
			openaiResponse.ID = "gemini"
			openaiResponse.Object = "chat.completion"
			// openaiResponse.Object = "chat.completion.chunk"
			openaiResponse.Created = 2467348
			openaiResponse.Model = model

			var openaiMessage define.Message
			var openaiChoice define.OpenAI_Choices
			for _, data := range responsePayload.Contents {
				for _, part := range data.Parts {
					openaiMessage.Role = data.Role
					openaiMessage.Content = part.Text
				}
			}

			openaiResponse.Usage.CompletionTokens = 11
			openaiResponse.Usage.PromptTokens = 12
			openaiResponse.Usage.TotalTokens = 23

			openaiChoice.FinishReason = "stop"
			openaiChoice.Message = openaiMessage
			openaiChoice.Index = 0
			openaiResponse.Choices = append(openaiResponse.Choices, openaiChoice)

			repack, err := json.Marshal(openaiResponse)
			if err != nil {
				return err
			}

			response.Body = io.NopCloser(bytes.NewBuffer(repack))
			response.ContentLength = int64(len(repack))
			response.Header.Set("Content-Length", strconv.Itoa(len(repack)))
		}
		return nil
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
