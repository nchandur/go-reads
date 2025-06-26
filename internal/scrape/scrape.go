package scrape

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/nchandur/go-reads/internal/models"
)

type Scraper struct {
	Page *rod.Page
}

func NewScraper(url string) (*Scraper, error) {
	page := rod.New().MustConnect().MustPage(url)
	err := rod.Try(func() {
		page.Timeout(25 * time.Second).MustNavigate(url)
	})
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}

	page.MustWaitLoad()
	return &Scraper{Page: page}, nil
}

func (s *Scraper) getTextFromSelector(selector string) (string, error) {
	el, err := s.Page.Timeout(25 * time.Second).Element(selector)
	if err != nil || el == nil {
		return "", fmt.Errorf("element not found")
	}
	text, _ := el.Text()
	return text, nil

}

func (s *Scraper) getBookID(url string) (int, error) {
	re := regexp.MustCompile(`[0-9]+`)

	match := re.FindStringSubmatch(url)

	if len(match) == 0 {
		return -1, fmt.Errorf("book id not found")
	}

	id, err := strconv.Atoi(match[0])

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Scraper) getAuthor() (string, error) {
	author, err := s.getTextFromSelector(`span[class="ContributorLink__name"]`)
	return author, err
}

func (s *Scraper) getTitle() (string, error) {
	title, err := s.getTextFromSelector(`h1[data-testid="bookTitle"]`)
	return title, err
}

func (s *Scraper) getSummary() (string, error) {
	firstButton, err := s.Page.Timeout(25 * time.Second).Element(".Button.Button--tertiary.Button--medium")
	if err != nil || firstButton == nil {
		return "", fmt.Errorf("failed find the button: %v", err)
	}

	err = firstButton.Click("left", 1)
	if err != nil {
		return "", fmt.Errorf("failed to click the button: %v", err)
	}

	s.Page.MustWaitIdle()

	truncEl, err := s.Page.Element(".TruncatedContent")
	if err != nil || truncEl == nil {
		return "", fmt.Errorf("failed find the summary element: %v", err)
	}

	summary, err := truncEl.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract summary text: %v", err)
	}

	summary = strings.ReplaceAll(summary, "\n", " ")

	return summary, nil
}

func (s *Scraper) getGenres() ([]string, error) {
	var genres []string
	moreButton, err := s.Page.Timeout(5*time.Second).ElementR(".Button.Button--tag.Button--medium", `\.\.\.more`)
	if err == nil && moreButton != nil {
		err = moreButton.Click("left", 1)
		if err != nil {
			return nil, fmt.Errorf("failed to click '...more' button: %v", err)
		}
		s.Page.MustWaitIdle()
	} else {
		return nil, fmt.Errorf("no '...more' button found or it took too long")
	}

	genreButtons, err := s.Page.Timeout(25 * time.Second).Elements(".BookPageMetadataSection__genreButton")
	if err != nil || len(genreButtons) == 0 {
		return genres, fmt.Errorf("no genre buttons found: %v", err)
	}

	for _, el := range genreButtons {
		text, _ := el.Text()
		genres = append(genres, text)
	}

	return genres, nil
}

func (s *Scraper) getStars() (float64, error) {
	stars, err := s.getTextFromSelector(`div.RatingStatistics__rating`)

	if err != nil {
		return -1, err
	}

	floatStars, err := strconv.ParseFloat(stars, 64)

	return floatStars, err
}

func (s *Scraper) getRatings() (int, error) {
	ratings, err := s.getTextFromSelector(`span[data-testid="ratingsCount"]`)

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

func (s *Scraper) getReviews() (int, error) {
	reviews, err := s.getTextFromSelector(`span[data-testid="reviewsCount"]`)

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

func (s *Scraper) getFormat() (models.Format, error) {
	format, err := s.getTextFromSelector(`p[data-testid="pagesFormat"]`)

	if err != nil {
		return models.Format{}, nil
	}

	re := regexp.MustCompile(`(\d+).*?,\s*(.+)`)
	matches := re.FindStringSubmatch(format)

	if len(matches) < 3 {
		return models.Format{}, fmt.Errorf("no matches found")
	}

	pageNo, err := strconv.Atoi(matches[1])
	if err != nil {
		return models.Format{}, fmt.Errorf("page no. not extracted")
	}

	bookType := strings.TrimSpace(matches[2])

	return models.Format{PageNo: pageNo, Type: bookType}, nil
}

func (s *Scraper) getDate() (time.Time, error) {

	published, err := s.getTextFromSelector(`p[data-testid="publicationInfo"]`)

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
		return time.Time{}, err
	}

	return date, nil
}

func (s *Scraper) getBookLinks() ([]string, error) {
	if err := s.Page.WaitLoad(); err != nil {
		return nil, err
	}

	var links []string

	res, err := s.Page.Timeout(25 * time.Second).Elements("a.bookTitle")

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
