package vectordb

import (
	"context"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

var Client *qdrant.Client

func Connect() error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "qdrant",
		Port: 6334,
	})

	if err != nil {
		return err
	}

	Client = client
	return nil
}

func Disconnect() error {
	err := Client.Close()

	if err != nil {
		return err
	}

	return nil
}

func CreateCollection(collection string) error {
	dimStr := os.Getenv("OLLAMA_EMBEDDING_DIM")

	dimInt, err := strconv.Atoi(dimStr)

	if err != nil {
		return err
	}

	dim := uint64(dimInt)

	err = Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     dim,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return err
	}

	return nil
}
