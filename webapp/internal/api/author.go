package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/recommend"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AuthorHandler(r *gin.Engine) {
	r.GET("/authors/:id", func(ctx *gin.Context) {
		authorid := ctx.Params.ByName("id")

		id, err := strconv.Atoi(authorid)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"body": nil, "error": err.Error()})
			return
		}

		collection := db.Client.Database("books").Collection("author")

		filter := bson.M{"authorid": id}

		opts := options.FindOne().SetProjection(bson.M{"_id": 0})

		author := models.Author{}

		err = collection.FindOne(ctx, filter, opts).Decode(&author)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"body": nil, "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": author, "error": nil})

	})

	r.GET("/authors", func(ctx *gin.Context) {
		name := ctx.Query("name")

		if len(name) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
			return
		}

		collection := db.Client.Database("books").Collection("author")

		author, err := recommend.SearchAuthor(ctx, collection, name)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": author, "error": nil})

	})

	r.GET("/authors/recommendations", func(ctx *gin.Context) {
		name := ctx.Query("name")

		if len(name) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
			return
		}

		n := ctx.Query("n")

		topK, err := strconv.Atoi(n)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": nil})
			return
		}

		collection := db.Client.Database("books").Collection("author")

		author, err := recommend.SearchAuthor(ctx, collection, name)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		authors, err := recommend.RecommendAuthor(ctx, collection, author.Name, topK)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "body": gin.H{"matched": author, "recommended": authors}})

	})

}
