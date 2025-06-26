package scrape

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

// fetches the genres the book belongs to
func extractGenres(page *rod.Page) ([]string, error) {
	var genres []string
	moreButton, err := page.Timeout(5*time.Second).ElementR(".Button.Button--tag.Button--medium", `\.\.\.more`)
	if err == nil && moreButton != nil {
		err = moreButton.Click("left", 1)
		if err != nil {
			return nil, fmt.Errorf("failed to click '...more' button: %v", err)
		}
		page.MustWaitIdle()
	} else {
		return nil, fmt.Errorf("no '...more' button found or it took too long")
	}

	genreButtons, err := page.Timeout(25 * time.Second).Elements(".BookPageMetadataSection__genreButton")
	if err != nil || len(genreButtons) == 0 {
		return genres, fmt.Errorf("no genre buttons found: %v", err)
	}

	for _, el := range genreButtons {
		text, _ := el.Text()
		genres = append(genres, text)
	}

	return genres, nil
}
