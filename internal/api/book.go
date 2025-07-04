package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/recommend"
)

func BookHandler(r *gin.Engine) {
	r.GET("/books/:id", func(ctx *gin.Context) {
		bookid := ctx.Params.ByName("id")

		id, err := strconv.Atoi(bookid)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"body": nil, "error": err.Error()})
			return
		}

		collection := db.Client.Database("books").Collection("works")

		book, err := recommend.SearchBookByID(ctx, collection, id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"body": nil, "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": book.Work, "error": nil})

	})

	r.GET("/books", func(ctx *gin.Context) {
		title := ctx.Query("title")

		if len(title) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"body": nil, "error": "title string cannot be empty"})
			return
		}

		collection := db.Client.Database("books").Collection("works")

		book, err := recommend.SearchBookByTitle(ctx, collection, title)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"body": nil, "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": book.Work, "error": nil})

	})

	r.GET("/books/recommendations", func(ctx *gin.Context) {
		title := ctx.Query("title")
		topK := ctx.Query("n")

		n, err := strconv.Atoi(topK)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": nil})
			return
		}

		collection := db.Client.Database("books").Collection("works")
		book, err := recommend.SearchBookByTitle(ctx, collection, title)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		books, err := recommend.RecommendByTitle(ctx, book.Work.Title, n)

		ctx.JSON(http.StatusOK, gin.H{"body": gin.H{"matched": book.Work, "recommended": books}, "error": nil})

	})
}
