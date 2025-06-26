package scrape

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/nchandur/go-reads/internal/models"
)

// fetch format of the book.
// format contains the number of pages in the book and the type of book
func fetchFormat(page *rod.Page) (models.Format, error) {
	format, err := getTextFromSelector(page, `p[data-testid="pagesFormat"]`)

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
