package main

import (
	"fmt"
	"log"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/recommend"
	"github.com/nchandur/go-reads/internal/vectordb"
)

func main() {
	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = db.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	err = vectordb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = vectordb.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	// books, err := recommend.RecommendByTitle("The Fellowship of The Ring", 10)

	books, err := recommend.RecommendByContext("frodo baggins destroys the ring", 5)

	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		fmt.Println(book.Title, book.Score)
	}

}
