package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GenreHandler(r *gin.Engine) {
	r.GET("/genres/authors", func(ctx *gin.Context) {
		genre := ctx.Query("genre")
		topK := ctx.Query("n")

		n, err := strconv.Atoi(topK)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": nil})
			return
		}

		collection := db.Client.Database("books").Collection("author")

		pipeline := mongo.Pipeline{
			{
				{Key: "$match", Value: bson.D{{Key: "genres", Value: genre}}},
			},
			{
				{Key: "$sort", Value: bson.D{
					{Key: "count", Value: -1},
					{Key: "stars", Value: -1},
				}},
			},
			{
				{Key: "$limit", Value: n},
			},
		}

		cursor, err := collection.Aggregate(ctx, pipeline)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		defer cursor.Close(ctx)

		var authors []models.Author

		if err := cursor.All(ctx, &authors); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": authors, "error": nil})

	})

	r.GET("/genres/books", func(ctx *gin.Context) {
		genre := ctx.Query("genre")
		topK := ctx.Query("n")

		n, err := strconv.Atoi(topK)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "body": nil})
			return
		}

		collection := db.Client.Database("books").Collection("works")

		pipeline := mongo.Pipeline{
			{
				{Key: "$match", Value: bson.D{{Key: "work.genres", Value: genre}}},
			},
			{
				{Key: "$sort", Value: bson.D{{Key: "work.stars", Value: -1}}},
			},
			{
				{Key: "$limit", Value: n},
			},
		}

		cursor, err := collection.Aggregate(ctx, pipeline)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		defer cursor.Close(ctx)

		var books []models.Book

		if err := cursor.All(ctx, &books); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "body": nil})
			return
		}

		works := []models.Work{}

		for _, book := range books {
			works = append(works, book.Work)
		}

		ctx.JSON(http.StatusOK, gin.H{"body": works, "error": nil})

	})

}
