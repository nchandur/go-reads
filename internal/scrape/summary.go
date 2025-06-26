package scrape

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// fetches the summary for the book
func extractSummary(page *rod.Page) (string, error) {

	firstButton, err := page.Timeout(25 * time.Second).Element(".Button.Button--tertiary.Button--medium")
	if err != nil || firstButton == nil {
		return "", fmt.Errorf("failed find the button: %v", err)
	}

	err = firstButton.Click("left", 1)
	if err != nil {
		return "", fmt.Errorf("failed to click the button: %v", err)
	}

	page.MustWaitIdle()

	truncEl, err := page.Element(".TruncatedContent")
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
