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

func Deduplicate(ctx context.Context) error {

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

func Invalid(ctx context.Context) error {

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

func uniqueGenres(genres []string) []string {
	res := []string{}

	hash := make(map[string]bool)

	for _, genre := range genres {
		hash[genre] = true
	}

	for key := range hash {
		res = append(res, key)
	}

	return res
}

func MakeAuthorCollection(ctx context.Context) error {
	collection := db.Client.Database("books").Collection("works")

	pipeline := mongo.Pipeline{
		{
			{Key: "$sort", Value: bson.D{
				{Key: "work.stars", Value: 1},
				{Key: "work.ratings", Value: 1},
				{Key: "work.reviews", Value: 1},
			}},
		},
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$work.author"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "stars", Value: bson.D{{Key: "$avg", Value: "$work.stars"}}},
				{Key: "reviews", Value: bson.D{{Key: "$avg", Value: "$work.reviews"}}},
				{Key: "ratings", Value: bson.D{{Key: "$avg", Value: "$work.ratings"}}},
				{Key: "bookids", Value: bson.D{{Key: "$push", Value: "$work.bookid"}}},
				{Key: "temp_genres", Value: bson.D{{Key: "$push", Value: "$work.genres"}}},
			}},
		},
		{
			{Key: "$unwind", Value: "$temp_genres"},
		},
		{
			{Key: "$unwind", Value: "$temp_genres"},
		},
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "count", Value: bson.D{{Key: "$first", Value: "$count"}}},
				{Key: "stars", Value: bson.D{{Key: "$first", Value: "$stars"}}},
				{Key: "reviews", Value: bson.D{{Key: "$first", Value: "$reviews"}}},
				{Key: "ratings", Value: bson.D{{Key: "$first", Value: "$ratings"}}},
				{Key: "bookids", Value: bson.D{{Key: "$first", Value: "$bookids"}}},
				{Key: "genres", Value: bson.D{{Key: "$addToSet", Value: "$temp_genres"}}},
			}},
		},
		{
			{Key: "$project", Value: bson.D{
				{Key: "name", Value: "$_id"},
				{Key: "count", Value: 1},
				{Key: "stars", Value: bson.D{{Key: "$round", Value: bson.A{"$stars", 3}}}},
				{Key: "reviews", Value: bson.D{{Key: "$round", Value: bson.A{"$reviews", 0}}}},
				{Key: "ratings", Value: bson.D{{Key: "$round", Value: bson.A{"$ratings", 0}}}},
				{Key: "books", Value: bson.D{{Key: "$slice", Value: bson.A{"$bookids", 5}}}},
				{Key: "genres", Value: bson.D{{Key: "$slice", Value: bson.A{"$genres", 10}}}},
				{Key: "_id", Value: 0},
			}},
		},
		{
			{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var authors []models.Author

	count := 0
	for cursor.Next(ctx) {
		var author models.Author

		err := cursor.Decode(&author)

		if err != nil {
			return err
		}
		count++
		author.AuthorID = count
		author.Genres = uniqueGenres(author.Genres)

		authors = append(authors, author)

	}

	authorCollection := db.Client.Database("books").Collection("author")

	_, err = authorCollection.InsertMany(ctx, authors)

	if err != nil {
		return err
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := Invalid(ctx); err != nil {
		log.Fatal(err)
	}

	if err := Deduplicate(ctx); err != nil {
		log.Fatal(err)
	}

	if err := MakeAuthorCollection(ctx); err != nil {
		log.Fatal(err)
	}

}
