package db

import (
	"context"

	"github.com/nchandur/go-reads/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InsertBooks(collection *mongo.Collection, document models.Book) error {

	_, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		return err
	}

	return nil
}
