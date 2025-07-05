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
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func createEmbedString(book models.Book) string {
	text := fmt.Sprintf("Title: %s\nAuthor: %s\nGenres: %v\nSummary: %s", book.Work.Title, book.Work.Author, book.Work.Genres, book.Work.Summary)

	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")

	text = reg.ReplaceAllString(text, "")

	return text
}

func main() {
	timeoutDuration := 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

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
		errLog.Println("Database connection error (if any occurred within db.Connect)")

	}

	links, err := scrape.FetchBookLinks(pageURL)

	if err != nil {
		errLog.Printf("error fetching links from %s: %v", pageURL, err)
		return
	} else {
		infoLog.Printf("fetched %d links from %s", len(links), pageURL)
	}

	collection := db.Client.Database("books").Collection("works")

	start := time.Now()

	for idx, link := range links {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				infoLog.Printf("Timeout reached after %v. Exiting loop prematurely.", timeoutDuration)
				fmt.Println("\nTimeout reached. Script exiting.")
			} else {
				infoLog.Printf("Context canceled for another reason. Exiting loop.")
				fmt.Println("\nScript canceled. Exiting.")
			}
			goto endLoop
		default:
		}

		book := scrape.FetchBookData(link, errLog)
		book.Url = link

		var doc models.Book

		doc.Source = pageURL
		doc.Work = book
		embedStr := createEmbedString(doc)

		doc.Embedding, err = ollama.Embed(embedStr)

		if err != nil {
			errLog.Printf("error embedding book %s: %v", book.Title, err.Error())
		}

		var exists models.Book

		filter := bson.M{"bookid": doc.Work.BookID}

		err = collection.FindOne(ctx, filter).Decode(&exists)

		if err == mongo.ErrNoDocuments {
			_, insertErr := collection.InsertOne(ctx, doc)

			if insertErr != nil {
				errLog.Printf("error inserting book %s: %v\n", book.Title, insertErr)
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

endLoop:
	fmt.Println()
	infoLog.Printf("extraction complete or terminated after %v", time.Since(start))
	fmt.Println("Total Time Elapsed: ", time.Since(start))
}
