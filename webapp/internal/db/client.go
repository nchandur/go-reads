package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func Connect() error {
	uri := os.Getenv("MONGO_DB_URI")
	cli, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	Client = cli
	return nil

}

func Disconnect() error {
	err := Client.Disconnect(context.Background())

	if err != nil {
		return err
	}

	return nil
}
