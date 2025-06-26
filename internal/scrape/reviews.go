package scrape

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

// fetch the number of reviews left for the book
func fetchReviews(page *rod.Page) (int, error) {
	reviews, err := getTextFromSelector(page, `span[data-testid="reviewsCount"]`)

	if err != nil {
		return -1, err
	}

	re := regexp.MustCompile(`[^0-9]+`)
	reviews = re.ReplaceAllString(reviews, "")

	intReviews, err := strconv.Atoi(strings.TrimSpace(reviews))

	if err != nil {
		return -1, err
	}

	return intReviews, nil

}
