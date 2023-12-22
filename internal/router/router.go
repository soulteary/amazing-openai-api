package router

import (
	"github.com/gin-gonic/gin"
	"github.com/soulteary/amazing-openai-api/models/azure"
)

func RegisterModelRoute(r *gin.Engine, serviceType string) {
	// https://platform.openai.com/docs/api-reference
	apiBase := "/v1"

	if serviceType == "azure" {
		stripPrefixConverter := azure.NewStripPrefixConverter(apiBase)
		r.GET(stripPrefixConverter.Prefix+"/models", azure.ModelProxy)
		apiBasedRouter := r.Group(apiBase)
		{
			apiBasedRouter.Any("/completions", azure.ProxyWithConverter(stripPrefixConverter))
			apiBasedRouter.Any("/chat/completions", azure.ProxyWithConverter(stripPrefixConverter))
		}
	}
}
