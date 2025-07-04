package models

import "fmt"

type RecommendedBook struct {
	BookID  int64   `json:"bookid"`
	Title   string  `json:"title"`
	Author  string  `json:"author"`
	Summary string  `json:"summary"`
	Genres  string  `json:"genres"`
	Stars   float64 `json:"stars"`
	Ratings int64   `json:"ratings"`
	Reviews int64   `json:"reviews"`
	Score   float32 `json:"score"`
}

func (r *RecommendedBook) Display() {
	fmt.Printf("Title: %s\nAuthor: %s\nSummary: %s\nGenres: %s\nStars: %f\nRatings: %d\nReviews: %d\nScore: %f\n", r.Title, r.Author, r.Summary, r.Genres, r.Stars, r.Ratings, r.Reviews, r.Score)
}

type RecommendedAuthor struct {
	AuthorID               int     `json:"authorid" bson:"authorid"`
	OverallSimilarityScore float64 `json:"score" bson:"overallSimilarityScore"`
}
