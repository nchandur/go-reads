package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"github.com/nchandur/go-reads/internal/ollama"
	"github.com/nchandur/go-reads/internal/scrape"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func createEmbedString(book models.Book) string {
	text := fmt.Sprintf("Title: %s\nAuthor: %s\nGenres: %v\nSummary: %s", book.Work.Title, book.Work.Author, book.Work.Genres, book.Work.Summary)

	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")

	text = reg.ReplaceAllString(text, "")

	return text
}

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

	if err != nil {
		errLog.Println(err)
		return
	}

	links, err := scrape.FetchBookLinks(pageURL)

	if err != nil {
		errLog.Printf("error fetching links from %s", pageURL)
		return
	} else {
		infoLog.Printf("fetched links from %s", pageURL)
	}

	collection := db.Client.Database("books").Collection("works")

	for idx, link := range links {

		if idx == 5 {
			break
		}

		book := scrape.FetchBookData(link, errLog)
		book.Url = link

		var doc models.Book

		doc.Source = pageURL
		doc.Work = book
		embedStr := createEmbedString(doc)

		doc.Embedding, err = ollama.Embed(embedStr)

		if err != nil {
			errLog.Printf("error embedding book: %v", err.Error())
		}

		var exists models.Book

		filter := bson.M{"bookid": doc.Work.BookID}

		err = collection.FindOne(context.Background(), filter).Decode(&exists)

		if err == mongo.ErrNoDocuments {
			_, insertErr := collection.InsertOne(context.Background(), doc)

			if insertErr != nil {
				errLog.Printf("error inserting book %s: %v\n", book.Title, err)
			} else {
				infoLog.Printf("%s pushed to DB\n", book.Title)
			}

		} else if err != nil {
			errLog.Printf("error checking for existing book %s: %v\n", book.Title, err)
		} else {
			infoLog.Printf("book %s already exists in DB. skipping insertion.\n", book.Title)
		}

		time.Sleep(5 * time.Second)
		fmt.Printf("\r%d books scraped", idx+1)

	}

	fmt.Println()
	infoLog.Printf("extraction complete")

}
