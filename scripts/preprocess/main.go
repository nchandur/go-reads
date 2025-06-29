package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Deduplicate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	collection := db.Client.Database("books").Collection("works")

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "title", Value: "$work.title"},
				{Key: "author", Value: "$work.author"},
			}},
			{Key: "ids", Value: bson.D{bson.E{Key: "$addToSet", Value: "$_id"}}},
			{Key: "count", Value: bson.D{bson.E{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{bson.E{Key: "$gt", Value: 1}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("aggregation failed: %v", err)
	}
	defer cursor.Close(ctx)

	var idsToDelete []any

	for cursor.Next(ctx) {
		var doc struct {
			IDs   []any `bson:"ids"`
			Count int   `bson:"count"`
		}
		if err = cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode aggregation result: %v", err)
		}

		if doc.Count > 1 {
			idsToDelete = append(idsToDelete, doc.IDs[1:]...)
		}
	}
	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor iteration error: %v", err)
	}

	if len(idsToDelete) > 0 {
		filter := bson.M{"_id": bson.M{"$in": idsToDelete}}

		result, err := collection.DeleteMany(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to delete duplicate documents: %v", err)
		}
		fmt.Printf("%d duplicate documents dropped\n", result.DeletedCount)
	} else {
		fmt.Println("No duplicate documents found to drop.")
	}

	return nil
}

func Invalid() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)

	defer cancel()

	collection := db.Client.Database("books").Collection("works")

	librarianNoteRegex := regexp.MustCompile(`(?i)^Librarian's note`)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to retrieve documents: %v", err)
	}

	defer func() error {
		if err := cursor.Close(ctx); err != nil {
			return fmt.Errorf("error closing cursor: %v", err)
		}
		return nil
	}()

	count := 0
	docCount := 0

	for cursor.Next(ctx) {
		docCount++
		var doc models.Book
		err := cursor.Decode(&doc)
		if err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		title := doc.Work.Title
		summary := doc.Work.Summary
		ratings := doc.Work.Ratings

		ratingsFlag := ratings == -1
		boxedFlag := strings.Contains(strings.ToLower(title), "box set") ||
			strings.Contains(strings.ToLower(title), "boxed set") ||
			strings.Contains(strings.ToLower(title), "boxset")

		libraryFlag := false
		if len(summary) > 0 {
			libraryFlag = librarianNoteRegex.MatchString(summary)
		}

		if ratingsFlag || boxedFlag || libraryFlag {
			_, err := collection.DeleteOne(ctx, bson.M{"work.bookid": doc.Work.BookID})
			if err != nil {
				log.Printf("Error deleting document with ID %v: %v", doc.Work.BookID, err)
			} else {
				count++
			}
		}

	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor iteration error: %v", err)
	}

	fmt.Println()
	fmt.Printf("%d invalid documents dropped\n", count)
	return nil
}

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

	if err := Invalid(); err != nil {
		log.Fatal(err)
	}

	if err := Deduplicate(); err != nil {
		log.Fatal(err)
	}

}
