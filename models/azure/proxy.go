package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/internal/fn"
	"github.com/soulteary/amazing-openai-api/internal/network"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

type DeploymentInfo struct {
	Data   []map[string]interface{} `json:"data"`
	Object string                   `json:"object"`
}

func ModelProxy(c *gin.Context) {
	// Create a channel to receive the results of each request
	results := make(chan []map[string]interface{}, len(ModelConfig))

	// Send a request for each deployment in the map
	for _, deployment := range ModelConfig {
		go func(deployment define.ModelConfig) {
			// Create the request
			req, err := http.NewRequest(http.MethodGet, deployment.Endpoint+"/openai/deployments?api-version=2022-12-01", nil)
			if err != nil {
				log.Printf("error parsing response body for deployment %s: %v", deployment.Name, err)
				results <- nil
				return
			}

			// Set the auth header
			req.Header.Set(HeaderAuthKey, deployment.Key)

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("error sending request for deployment %s: %v", deployment.Name, err)
				results <- nil
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Printf("unexpected status code %d for deployment %s", resp.StatusCode, deployment.Name)
				results <- nil
				return
			}

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error reading response body for deployment %s: %v", deployment.Name, err)
				results <- nil
				return
			}

			// Parse the response body as JSON
			var deplotmentInfo DeploymentInfo
			err = json.Unmarshal(body, &deplotmentInfo)
			if err != nil {
				log.Printf("error parsing response body for deployment %s: %v", deployment.Name, err)
				results <- nil
				return
			}
			results <- deplotmentInfo.Data
		}(deployment)
	}

	// Wait for all requests to finish and collect the results
	var allResults []map[string]interface{}
	for i := 0; i < len(ModelConfig); i++ {
		result := <-results
		if result != nil {
			allResults = append(allResults, result...)
		}
	}
	var info = DeploymentInfo{Data: allResults, Object: "list"}
	combinedResults, err := json.Marshal(info)
	if err != nil {
		log.Printf("error marshalling results: %v", err)
		network.SendError(c, err)
		return
	}

	// Set the response headers and body
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(combinedResults))
}

// Proxy Azure OpenAI
func Proxy(c *gin.Context, requestConverter RequestConverter) {
	// preserve request body for error logging
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	bodyBytes, err := io.ReadAll(tee)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return
	}
	c.Request.Body = io.NopCloser(&buf)

	director := func(req *http.Request) {
		if req.Body == nil {
			network.SendError(c, errors.New("request body is empty"))
			return
		}

		// extract model from request url
		model := c.Param("model")
		if model == "" {
			// extract model from request body
			body, err := io.ReadAll(req.Body)
			defer req.Body.Close()
			if err != nil {
				network.SendError(c, errors.Wrap(err, "read request body error"))
				return
			}

			var requestData define.RequestData
			err = json.Unmarshal(body, &requestData)
			if err != nil {
				network.SendError(c, errors.Wrap(err, "parse payload error"))
				return
			}

			model = requestData.Model
			// TODO change alias to model
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

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
		req.Header.Set(HeaderAuthKey, token)
		req.Header.Del("Authorization")

		originURL := req.URL.String()
		req, err = requestConverter.Convert(req, deployment)
		if err != nil {
			network.SendError(c, errors.Wrap(err, "convert request error"))
			return
		}
		log.Printf("proxying request [%s] %s -> %s", model, originURL, req.URL.String())
	}

	proxy := &httputil.ReverseProxy{Director: director}
	transport, err := network.NewProxyFromEnv(
		fn.GetStringOrDefaultFromEnv("ENV_AZURE_SOCKS_PROXY", ""),
		fn.GetStringOrDefaultFromEnv("ENV_AZURE_HTTP_PROXY", ""),
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
		log.Printf("encountering error with body: %s", string(bodyBytes))
	}
}

func GetDeploymentByModel(model string) (*define.ModelConfig, error) {
	deploymentConfig, exist := ModelConfig[model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("deployment config for %s not found", model))
	}
	return &deploymentConfig, nil
}
