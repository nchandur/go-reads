package recommend

import (
	"github.com/nchandur/go-reads/internal/models"
	"github.com/qdrant/go-client/qdrant"
)

func getDocs(points []*qdrant.ScoredPoint) ([]models.RecommendedBook, error) {

	books := []models.RecommendedBook{}

	for _, point := range points {

		var book models.RecommendedBook

		book.BookID = point.GetPayload()["bookid"].GetIntegerValue()
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

	return books[1:], nil
}
