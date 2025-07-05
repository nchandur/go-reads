package models

type Author struct {
	AuthorID   int      `json:"authorid" bson:"authorid"`
	Name       string   `json:"name" bson:"name"`
	BookCount  int      `json:"count" bson:"count"`
	AvgStars   float64  `json:"stars" bson:"stars"`
	AvgRatings int      `json:"ratings" bson:"ratings"`
	AvgReviews int      `json:"reviews" bson:"reviews"`
	TopBooks   []int    `json:"books" bson:"books"`
	Genres     []string `json:"genres" bson:"genres"`
}
