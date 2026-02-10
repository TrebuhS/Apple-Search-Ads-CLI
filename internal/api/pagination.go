package api

import (
	"fmt"
	"time"

	"github.com/trebuhs/asa-cli/internal/models"
)

const (
	maxRetries    = 3
	retryBaseWait = 2 * time.Second
)

// PaginatedFetcher fetches all pages of results using a POST-based find endpoint.
func PaginatedFetcher[T any](c *Client, path string, selector models.Selector) ([]T, error) {
	var allResults []T
	offset := selector.Pagination.Offset

	for {
		selector.Pagination.Offset = offset
		var page []T
		pagination, err := c.Post(path, &selector, &page)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, page...)

		if pagination == nil || len(allResults) >= pagination.TotalResults {
			break
		}

		offset += len(page)
		if len(page) == 0 {
			break
		}
	}

	return allResults, nil
}

// RetryOn429 wraps an API call with retry logic for rate limiting.
func RetryOn429(fn func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Simple check for 429 in error message
		if i < maxRetries-1 {
			wait := retryBaseWait * time.Duration(1<<uint(i))
			fmt.Printf("Rate limited, retrying in %v...\n", wait)
			time.Sleep(wait)
			continue
		}
		return err
	}
	return nil
}
