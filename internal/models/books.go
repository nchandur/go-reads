package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	DocID  primitive.ObjectID `json:"_id"`
	Source string             `json:"source"`
	Work   struct {
		BookID  int      `json:"book_id"`
		Title   string   `json:"title"`
		Author  string   `json:"author"`
		Summary string   `json:"summary"`
		Genres  []string `json:"genres"`
		Stars   float64  `json:"stars"`
		Ratings int      `json:"ratings"`
		Reviews int      `json:"reviews"`
		Format  struct {
			PageNo int    `json:"page_no"`
			Type   string `json:"type"`
		} `json:"format"`
		Published time.Time `json:"published"`
		Url       string    `json:"url"`
	} `json:"work"`
}

func (b *Book) Display() {
	fmt.Printf("Title: %s\nAuthor: %s\nSummary: %s\nGenres: %v\nStars: %f\nRatings: %d, Reviews: %d\nPublished: %v\n", b.Work.Title, b.Work.Author, b.Work.Summary, b.Work.Genres, b.Work.Stars, b.Work.Ratings, b.Work.Reviews, b.Work.Published)
}
