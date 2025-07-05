package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/nchandur/go-reads/internal/vectordb"
)

func HealthChecker(r *gin.Engine) {

	r.GET("/health/mongodb", func(ctx *gin.Context) {
		if err := db.Client.Ping(ctx, nil); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "body": "MongoDB service is up!"})

	})

	r.GET("/health/ollama", func(ctx *gin.Context) {
		if _, err := ollama.Embed("hello world"); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "body": "Ollama service is up!"})
	})

	r.GET("/health/vectordb", func(ctx *gin.Context) {
		if _, err := vectordb.Client.HealthCheck(ctx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "body": "VectorDB service is up!"})
	})

}
