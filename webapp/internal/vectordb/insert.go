package vectordb

import (
	"context"
	"strings"

	"github.com/nchandur/go-reads/internal/models"
	"github.com/qdrant/go-client/qdrant"
)

func InsertDoc(collection string, id uint64, book models.Book) error {

	genreStr := strings.Join(book.Work.Genres, ", ")

	_, err := Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collection,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(id),
				Vectors: qdrant.NewVectors(book.Embedding...),
				Payload: qdrant.NewValueMap(map[string]any{
					"bookid":  book.Work.BookID,
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

	if err != nil {
		return err
	}

	return nil

}
