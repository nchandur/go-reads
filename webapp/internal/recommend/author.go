package recommend

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nchandur/go-reads/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func SearchAuthor(ctx context.Context, collection *mongo.Collection, author string) (models.Author, error) {
	exactRegex := fmt.Sprintf("^%s$", regexp.QuoteMeta(author))
	exactFilter := bson.M{
		"name": bson.M{
			"$regex":   exactRegex,
			"$options": "i",
		},
	}

	var exactResult models.Author
	err := collection.FindOne(ctx, exactFilter).Decode(&exactResult)
	if err == nil {
		return exactResult, nil
	}
	if err != mongo.ErrNoDocuments {
		return models.Author{}, err
	}

	substringFilter := bson.M{
		"name": bson.M{
			"$regex":   author,
			"$options": "i",
		},
	}

	opts := options.FindOne().
		SetSort(bson.D{{Key: "name", Value: 1}})

	if err = collection.FindOne(ctx, substringFilter, opts).Decode(&exactResult); err != nil {
		return exactResult, err
	}

	return exactResult, nil

}

func RecommendAuthor(ctx context.Context, collection *mongo.Collection, name string, topK int) ([]models.RecommendedAuthor, error) {

	author, err := SearchAuthor(ctx, collection, name)

	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{

		bson.D{{Key: "$set", Value: bson.D{
			{Key: "targetGenres", Value: author.Genres},
		}}},

		bson.D{{Key: "$match", Value: bson.D{
			{Key: "authorid", Value: bson.D{{Key: "$ne", Value: author.AuthorID}}},
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "genreMatchCount", Value: bson.D{
				{Key: "$size", Value: bson.D{
					{Key: "$setIntersection", Value: bson.A{"$genres", "$targetGenres"}},
				}},
			}},
		}}},

		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "overallSimilarityScore", Value: bson.D{
				{Key: "$sum", Value: bson.A{
					bson.D{{Key: "$multiply", Value: bson.A{"$genreMatchCount", 1.0}}},
				}},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "overallSimilarityScore", Value: -1},
		}}},
		bson.D{{Key: "$limit", Value: topK}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "authorid", Value: 1},
			{Key: "overallSimilarityScore", Value: 1},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var similarAuthors []models.RecommendedAuthor

	if err = cursor.All(ctx, &similarAuthors); err != nil {
		return nil, err
	}

	return similarAuthors, nil

}
