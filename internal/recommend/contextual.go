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
		return nil, err
	}

	return getDocs(points)
}
