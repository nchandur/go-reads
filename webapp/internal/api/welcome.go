package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"body": "Welcome to GoReads!", "error": nil})
	})
}
