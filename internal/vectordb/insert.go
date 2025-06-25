package vectordb

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/qdrant/go-client/qdrant"
)

func createEmbedString(book models.Book) string {
	text := fmt.Sprintf("Title: %s\nAuthor: %s\nGenres: %v\nSummary: %s", book.Work.Title, book.Work.Author, book.Work.Genres, book.Work.Summary)

	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")

	text = reg.ReplaceAllString(text, "")

	return text
}

func InsertDoc(collection string, id uint64, book models.Book) error {

	vec, err := ollama.Embed(createEmbedString(book))

	if err != nil {
		return err
	}

	genreStr := strings.Join(book.Work.Genres, ", ")

	_, err = Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collection,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(id),
				Vectors: qdrant.NewVectors(vec...),
				Payload: qdrant.NewValueMap(map[string]any{
					"book_id": book.Work.BookID,
					"title":   book.Work.Title,
					"author":  book.Work.Author,
					"summary": book.Work.Summary,
					"genres":  genreStr,
					"stars":   book.Work.Stars,
					"ratings": book.Work.Ratings,
					"reviews": book.Work.Reviews,
				}),
			},
		},
	})

	return nil

}
