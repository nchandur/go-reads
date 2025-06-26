package scrape

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/nchandur/go-reads/internal/models"
)

// util function to safely retrieve the text in the given HTML tag
func getTextFromSelector(page *rod.Page, selector string) (string, error) {
	el, err := page.Timeout(25 * time.Second).Element(selector)
	if err != nil || el == nil {
		return "", fmt.Errorf("element not found")
	}
	text, _ := el.Text()
	return text, nil

}

// fetch data for the book in the given URL
func FetchBookData(url string, errLog *log.Logger) models.Work {
	page := rod.New().MustConnect().MustPage()
	err := rod.Try(func() {
		page.Timeout(25 * time.Second).MustNavigate(url)
	})
	if errors.Is(err, context.DeadlineExceeded) {
		errLog.Printf("page timeout: %v", err)
		return models.Work{}
	}
	defer page.Close()

	page.MustWaitLoad()

	var book models.Work

	book.Url = url

	book.BookID, err = getBookID(url)

	if err != nil {
		errLog.Printf("id not extracted: %v", err)
	}

	book.Author, err = getTextFromSelector(page, `span[class="ContributorLink__name"]`)

	if err != nil {
		errLog.Printf("author not extracted: %v", err)
	}

	book.Title, err = getTextFromSelector(page, `h1[data-testid="bookTitle"]`)

	if err != nil {
		book.Title = ""
		errLog.Printf("title not extracted: %v", err)
	}

	book.Summary, err = extractSummary(page)

	if err != nil {
		errLog.Printf("summary not extracted: %v", err)
	}

	book.Genres, err = extractGenres(page)

	if err != nil {
		errLog.Printf("genres not extracted: %v", err)
	}

	stars, err := getTextFromSelector(page, `div.RatingStatistics__rating`)

	if err != nil {
		errLog.Printf("stars not extracted: %v", err)
	} else {
		book.Stars, err = strconv.ParseFloat(stars, 64)
	}

	if err != nil {
		errLog.Printf("stars not extracted: %v", err)
	}

	book.Ratings, err = fetchRatings(page)

	if err != nil {
		errLog.Printf("ratings not extracted: %v", err)
	}

	book.Reviews, err = fetchReviews(page)

	if err != nil {
		errLog.Printf("reviews not extracted: %v", err)
	}

	book.Format, err = fetchFormat(page)

	if err != nil {
		errLog.Printf("format not extracted: %v", err)
	}

	book.Published, err = fetchDate(page)

	if err != nil {
		errLog.Printf("publication not extracted: %v", err)
	}

	return book

}
