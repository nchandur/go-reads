package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/recommend"
)

func RecommendHandler(r *gin.Engine) {
	r.GET("/recommend/title", func(ctx *gin.Context) {
		title := ctx.Query("title")
		n := ctx.Query("n")

		if len(title) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "title string cannot be empty", "body": nil})
			return
		}

		topK, err := strconv.Atoi(n)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": nil})
			return
		}

		if (topK <= 0) || (topK > 25) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "n should lie between 0 and 25", "body": nil})
			return
		}

		books, err := recommend.RecommendByTitle(title, topK)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": books, "error": nil})

	})

}
