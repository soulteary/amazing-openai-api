package network

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Error ErrorDescription `json:"error"`
}

type ErrorDescription struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func SendError(c *gin.Context, err error) {
	c.JSON(500, ApiResponse{
		Error: ErrorDescription{
			Code:    "500",
			Message: err.Error(),
		},
	})
}
