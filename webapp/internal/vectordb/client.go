package vectordb

import (
	"context"

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

	err := Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     768,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return err
	}

	return nil
}
