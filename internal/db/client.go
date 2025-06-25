package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func Connect() error {
	cli, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:9001"))

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
