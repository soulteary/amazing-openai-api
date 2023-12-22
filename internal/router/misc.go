package router

import "github.com/gin-gonic/gin"

func Hi(c *gin.Context) {
	c.Status(200)
}

func RegisterMiscRoute(r *gin.Engine) {
	r.GET("/", Hi)
	r.GET("/health", Hi)
	r.GET("/ping", Hi)
}
