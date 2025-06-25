package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/nchandur/go-reads/internal/vectordb"
	"github.com/qdrant/go-client/qdrant"
)

func main() {
	err := vectordb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to VectorDB: %v", err)
	}
	defer func() {
		if err := vectordb.Disconnect(); err != nil {
			log.Printf("Error disconnecting from VectorDB: %v", err)
		}
	}()

	text := "divergent"

	vec, err := ollama.Embed(text)

	if err != nil {
		log.Fatal(err)
	}

	topK := uint64(10)

	points, err := vectordb.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "books",
		Query:          qdrant.NewQuery(vec...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          &topK,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, point := range points {
		title := point.GetPayload()["title"]
		author := point.GetPayload()["author"]
		score := point.GetScore()

		fmt.Printf("Title: %v, Author: %v, Score: %v\n", title, author, score)

	}

}
