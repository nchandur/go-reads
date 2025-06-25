package recommend

import (
	"context"

	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/nchandur/go-reads/internal/vectordb"
	"github.com/qdrant/go-client/qdrant"
)

func RecommendByContext(contextString string, n int) ([]models.RecommendedBook, error) {
	vec, err := ollama.Embed(contextString)

	if err != nil {
		return nil, err
	}

	topK := uint64(n)

	points, err := vectordb.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "books",
		Query:          qdrant.NewQuery(vec...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          &topK,
	})

	if err != nil {
		return nil, nil
	}

	books := []models.RecommendedBook{}

	for _, point := range points {

		var book models.RecommendedBook

		book.BookID = point.GetPayload()["book_id"].GetIntegerValue()
		book.Title = point.GetPayload()["title"].GetStringValue()
		book.Author = point.GetPayload()["author"].GetStringValue()
		book.Summary = point.GetPayload()["summary"].GetStringValue()
		book.Genres = point.GetPayload()["genres"].GetStringValue()
		book.Stars = point.GetPayload()["stars"].GetDoubleValue()
		book.Ratings = point.GetPayload()["ratings"].GetIntegerValue()
		book.Reviews = point.GetPayload()["reviews"].GetIntegerValue()

		book.Score = point.GetScore()

		books = append(books, book)

	}

	return books, nil

}
