package scrape

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// fetch the URLS for books in a list in Goodreads
func FetchBookLinks(url string) ([]string, error) {
	page := rod.New().MustConnect().MustPage()

	err := rod.Try(func() {
		page.Timeout(25 * time.Second).MustNavigate(url)
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf("page timeout: %v", err)
	}

	defer page.Close()

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	var links []string

	res, err := page.Timeout(25 * time.Second).Elements("a.bookTitle")

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		href, err := r.Attribute("href")

		if err == nil && href != nil {
			link := fmt.Sprintf("https://www.goodreads.com%s", *href)
			links = append(links, strings.TrimSpace(link))
		}
	}

	return links, nil
}
