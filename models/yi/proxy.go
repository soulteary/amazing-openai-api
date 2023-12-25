package yi

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
	HeaderAPIVer = "api-version"
)

var maskURL = regexp.MustCompile(`https?:\/\/.+\/v1\/`)

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

// Proxy YI
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
		req.Header.Set("Authorization", token)

		originURL := req.URL.String()
		req, err = requestConverter.Convert(req, deployment)
		if err != nil {
			network.SendError(c, errors.Wrap(err, "convert request error"))
			return
		}

		log.Printf("proxying request [%s] %s -> %s", model, originURL, maskURL.ReplaceAllString(req.URL.String(), "${YI-API-SERVER}/v1/"))
	}

	proxy := &httputil.ReverseProxy{Director: director}
	transport, err := network.NewProxyFromEnv(
		fn.GetStringOrDefaultFromEnv("ENV_YI_SOCKS_PROXY", ""),
		fn.GetStringOrDefaultFromEnv("ENV_YI_HTTP_PROXY", ""),
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
