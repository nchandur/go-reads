package scrape

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

// fetch the number of ratings given to a book
func fetchRatings(page *rod.Page) (int, error) {
	ratings, err := getTextFromSelector(page, `span[data-testid="ratingsCount"]`)

	if err != nil {
		return -1, err
	}

	re := regexp.MustCompile(`[^0-9]+`)
	ratings = re.ReplaceAllString(ratings, "")

	intRatings, err := strconv.Atoi(strings.TrimSpace(ratings))

	if err != nil {
		return -1, err
	}

	return intRatings, nil

}
