package scrape

import (
	"fmt"
	"regexp"
	"strconv"
)

// extract book ID from book URL (Goodreads ID)
func getBookID(str string) (int, error) {
	re := regexp.MustCompile(`[0-9]+`)

	match := re.FindStringSubmatch(str)

	if len(match) == 0 {
		return -1, fmt.Errorf("book id not found")
	}

	id, err := strconv.Atoi(match[0])

	if err != nil {
		return -1, err
	}

	return id, nil

}
