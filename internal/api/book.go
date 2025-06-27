package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/recommend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func BookHandler(r *gin.Engine) {
	r.GET("/book/id/:id", func(ctx *gin.Context) {
		bookid := ctx.Params.ByName("id")

		id, err := strconv.Atoi(bookid)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"body": nil, "error": err.Error()})
			return
		}

		collection := db.Client.Database("books").Collection("works")

		filter := bson.M{"work.bookid": id}

		opts := options.FindOne().SetProjection(bson.M{"_id": 0})

		book := models.Book{}

		err = collection.FindOne(ctx, filter, opts).Decode(&book)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"body": nil, "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": book.Work, "error": nil})

	})

	r.GET("/book/title", func(ctx *gin.Context) {
		title := ctx.Query("title")

		if len(title) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"body": nil, "error": "title string cannot be empty"})
			return
		}

		collection := db.Client.Database("books").Collection("works")

		book, err := recommend.SearchBook(collection, title)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"body": nil, "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": book.Work, "error": nil})

	})
}
