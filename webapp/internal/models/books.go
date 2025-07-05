package models

import (
	"fmt"
	"time"
)

type Book struct {
	Source    string    `json:"source"`
	Work      Work      `json:"work"`
	Embedding []float32 `json:"embedding"`
}

type Work struct {
	BookID    int       `json:"bookid" bson:"bookid"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Summary   string    `json:"summary"`
	Genres    []string  `json:"genres"`
	Stars     float64   `json:"stars"`
	Ratings   int       `json:"ratings"`
	Reviews   int       `json:"reviews"`
	Format    Format    `json:"format"`
	Published time.Time `json:"published"`
	Url       string    `json:"url"`
}

type Format struct {
	PageNo int    `json:"page_no"`
	Type   string `json:"type"`
}

func (b *Book) Display() {
	fmt.Printf("BookID: %d\nTitle: %s\nAuthor: %s\nSummary: %s\nGenres: %v\nStars: %f\nRatings: %d, Reviews: %d, Published: %v\nPageNo: %d, Type: %s\n", b.Work.BookID, b.Work.Title, b.Work.Author, b.Work.Summary, b.Work.Genres, b.Work.Stars, b.Work.Ratings, b.Work.Reviews, b.Work.Published, b.Work.Format.PageNo, b.Work.Format.Type)
}
