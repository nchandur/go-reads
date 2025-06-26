package scrape

import "github.com/nchandur/go-reads/internal/models"

func Fetch(url string) (models.Book, error) {

	scraper, err := NewScraper(url)

	if err != nil {
		return models.Book{}, err
	}

	defer scraper.Page.Close()

	var book models.Book

	book.Work.Url = url

	book.Work.BookID, err = scraper.getBookID(url)

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Title, err = scraper.getTitle()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Author, err = scraper.getAuthor()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Summary, err = scraper.getSummary()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Genres, err = scraper.getGenres()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Stars, err = scraper.getStars()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Ratings, err = scraper.getRatings()

	if err != nil {
		return models.Book{}, err
	}

	book.Work.Reviews, err = scraper.getReviews()
	if err != nil {
		return models.Book{}, err
	}

	book.Work.Format, err = scraper.getFormat()
	if err != nil {
		return models.Book{}, err
	}

	book.Work.Published, err = scraper.getDate()

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func FetchBookLinks(url string) ([]string, error) {
	scraper, err := NewScraper(url)

	if err != nil {
		return nil, err
	}

	defer scraper.Page.Close()

	return scraper.getBookLinks()
}
