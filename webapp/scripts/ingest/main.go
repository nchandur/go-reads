package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/vectordb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	err := vectordb.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer vectordb.Disconnect()

	err = vectordb.CreateCollection("books")

	if err != nil {
		fmt.Println(err)
	}

	err = db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	if err != nil {
		log.Fatal(err)
	}

	collection := db.Client.Database("books").Collection("works")

	cur, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	count := 0
	failed := 0
	start := time.Now()

	for cur.Next(context.Background()) {
		count++
		var book models.Book

		if err = cur.Decode(&book); err != nil {
			continue
		}

		err = vectordb.InsertDoc("books", uint64(count), book)

		if err != nil {
			failed++
		}

		fmt.Printf("\r%d books processed.", count)

	}

	fmt.Printf("\n%d books failed during ingestion\nTime Taken for Embedding Books: %v\n", failed, time.Since(start))

}
