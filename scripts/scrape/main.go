package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/scrape"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	pageURL := os.Args[1]

	logFile, err := os.OpenFile("data/extraction.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	infoLog := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errLog := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	db.Connect()

	defer db.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)

	defer cancel()

	links, err := scrape.FetchBookLinks(pageURL)

	if err != nil {
		errLog.Println(err)
	} else {
		infoLog.Println("extracted links from list: ", pageURL)
	}

	collection := db.Client.Database("books").Collection("works")

	for idx, link := range links {

		if idx == 5 {
			break
		}

		book, err := scrape.Fetch(link)

		if err != nil {
			errLog.Println(err)
		}

		book.Source = pageURL

		var exists models.Book

		filter := bson.M{"work.bookid": book.Work.BookID}

		err = collection.FindOne(ctx, filter).Decode(&exists)

		if err == mongo.ErrNoDocuments {
			_, insertErr := collection.InsertOne(ctx, book)
			if insertErr != nil {
				errLog.Printf("error inserting book %s: %v\n", book.Work.Title, insertErr)
			} else {
				infoLog.Printf("%s pushed to DB\n", book.Work.Title)
			}
		} else if err != nil {
			errLog.Printf("error checking for existing book %s: %v\n", book.Work.Title, err)
		} else {
			infoLog.Printf("book %s already exists in DB. skipping insertion.\n", book.Work.Title)
		}
		time.Sleep(time.Duration(rand.IntN(10) + 5))

		fmt.Printf("\r%d books processed.", idx+1)

	}

}
