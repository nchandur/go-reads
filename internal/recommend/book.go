package recommend

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/nchandur/go-reads/internal/vectordb"
	"github.com/qdrant/go-client/qdrant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SearchBookByTitle(ctx context.Context, collection *mongo.Collection, title string) (models.Book, error) {
	exactRegex := fmt.Sprintf("^%s$", regexp.QuoteMeta(title))
	exactFilter := bson.M{
		"work.title": bson.M{
			"$regex":   exactRegex,
			"$options": "i",
		},
	}

	var exactResult models.Book
	err := collection.FindOne(ctx, exactFilter).Decode(&exactResult)
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

	if err = collection.FindOne(ctx, substringFilter, opts).Decode(&exactResult); err != nil {
		return exactResult, err
	}

	return exactResult, nil

}

func SearchBookByID(ctx context.Context, collection *mongo.Collection, id int) (models.Book, error) {
	filter := bson.M{
		"work.bookid": id,
	}

	opts := options.FindOne().SetProjection(bson.M{"_id": 0})

	var book models.Book
	err := collection.FindOne(ctx, filter, opts).Decode(&book)

	if err != nil {
		return models.Book{}, err
	}

	return book, nil

}

func RecommendByTitle(ctx context.Context, title string, n int) ([]models.RecommendedBook, error) {
	collection := db.Client.Database("books").Collection("works")

	book, err := SearchBookByTitle(ctx, collection, title)

	if err != nil {
		return nil, err
	}

	topK := uint64(n + 1)

	points, err := vectordb.Client.Query(ctx, &qdrant.QueryPoints{
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

func RecommendByContext(ctx context.Context, contextString string, n int) ([]models.RecommendedBook, error) {
	vec, err := ollama.Embed(contextString)

	if err != nil {
		return nil, err
	}

	topK := uint64(n)

	points, err := vectordb.Client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "books",
		Query:          qdrant.NewQuery(vec...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          &topK,
	})

	if err != nil {
		return nil, err
	}

	return getDocs(points)
}
