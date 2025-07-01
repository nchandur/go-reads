package recommend

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/vectordb"
	"github.com/qdrant/go-client/qdrant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func createEmbedString(book models.Book) string {
	text := fmt.Sprintf("Title: %s\nAuthor: %s\nGenres: %v\nSummary: %s", book.Work.Title, book.Work.Author, book.Work.Genres, book.Work.Summary)

	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")

	text = reg.ReplaceAllString(text, "")

	return text
}

func SearchBook(collection *mongo.Collection, title string) (models.Book, error) {
	exactRegex := fmt.Sprintf("^%s$", regexp.QuoteMeta(title))
	exactFilter := bson.M{
		"work.title": bson.M{
			"$regex":   exactRegex,
			"$options": "i",
		},
	}

	var exactResult models.Book
	err := collection.FindOne(context.Background(), exactFilter).Decode(&exactResult)
	if err == nil {
		return exactResult, nil
	}
	if err != mongo.ErrNoDocuments {
		return models.Book{}, err
	}

	substringFilter := bson.M{
		"work.title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}

	opts := options.FindOne().
		SetSort(bson.D{{Key: "work.title", Value: 1}})

	if err = collection.FindOne(context.Background(), substringFilter, opts).Decode(&exactResult); err != nil {
		return exactResult, err
	}

	return exactResult, nil

}

func RecommendByTitle(title string, n int) ([]models.RecommendedBook, error) {
	collection := db.Client.Database("books").Collection("works")

	book, err := SearchBook(collection, title)

	if err != nil {
		return nil, err
	}

	topK := uint64(n + 1)

	points, err := vectordb.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "books",
		Query:          qdrant.NewQuery(book.Embedding...),
		Limit:          &topK,
		WithPayload:    qdrant.NewWithPayload(true),
	})

	if err != nil {
		log.Fatalf("failed to search points: %v", err)
	}
	return getDocs(points)

}
