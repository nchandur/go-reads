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
		if err := db.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	err = vectordb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := vectordb.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	collection := db.Client.Database("books").Collection("works")

	title := "harry potter"
	book, err := recommend.SearchBook(collection, title)

	if err != nil {
		log.Fatal(err)
	}

	books, err := recommend.RecommendByTitle(book.Work.Title, 5)

	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		book.Display()

		fmt.Println("---------------------------------------------------------------------------------------")
	}

}
