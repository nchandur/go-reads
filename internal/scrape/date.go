package scrape

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// fetch the first published date (or expected date of publication) for the book

func fetchDate(page *rod.Page) (time.Time, error) {

	published, err := getTextFromSelector(page, `p[data-testid="publicationInfo"]`)

	if err != nil {
		return time.Time{}, err
	}

	re := regexp.MustCompile(`([A-Za-z]+\s+\d{1,2}\,\s+\d{1,4})`)
	match := re.FindString(published)

	if match == "" {
		return time.Time{}, fmt.Errorf("failed to parse date")
	}

	parts := strings.Split(match, ", ")
	year, err := strconv.Atoi(parts[1])

	if err != nil {
		log.Fatal(err)
	}

	paddedDate := parts[0] + ", " + fmt.Sprintf("%04d", year)

	layout := "January 2, 2006"

	date, err := time.Parse(layout, paddedDate)

	if err != nil {
		log.Fatal(err)
	}

	return date, nil
}
