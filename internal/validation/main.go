// Package validation provides utilities for validating URLs and processing them concurrently.
package validation

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"

	"link-builder/internal/utils"
)

func ValidateURLsConcurrently(urls []string, ignoreRegex *regexp.Regexp) (map[string]bool, int) {
	validURLs := make(map[string]bool)
	ignoredCount := int32(0)
	urlChan := make(chan string, len(urls))
	resultChan := make(chan struct {
		url     string
		isValid bool
		ignored bool
	}, len(urls))

	for _, rawURL := range urls {
		urlChan <- rawURL
	}
	close(urlChan)

	worker := func() {
		for rawURL := range urlChan {
			if ignoreRegex != nil && ignoreRegex.MatchString(rawURL) {
				resultChan <- struct {
					url     string
					isValid bool
					ignored bool
				}{url: rawURL, isValid: false, ignored: true}
				continue
			}
			isValid := validateURL(rawURL)
			resultChan <- struct {
				url     string
				isValid bool
				ignored bool
			}{url: rawURL, isValid: isValid, ignored: false}
		}
	}

	workerCount := 10
	for range make([]struct{}, workerCount) {
		go worker()
	}

	for range urls {
		result := <-resultChan
		if result.ignored {
			atomic.AddInt32(&ignoredCount, 1)
		} else if result.isValid {
			validURLs[result.url] = true
		}
	}

	return validURLs, int(ignoredCount)
}

func validateURL(url string) bool {
	return utils.IsValidURL(url)
}

func ProcessURLs(validURLs map[string]bool) map[string]bool {
	processedURLs := make(map[string]bool)
	for urlStr := range validURLs {
		if semicolonIndex := strings.Index(urlStr, ";jsessionid="); semicolonIndex != -1 {
			urlStr = urlStr[:semicolonIndex]
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("Failed to parse URL %s: %v", urlStr, err)
			continue
		}
		query := parsedURL.Query()
		for key := range query {
			if strings.Contains(strings.ToLower(key), "session") {
				query.Del(key)
				log.Printf("Warning: URL contains 'session': %s", urlStr)
			}
		}
		parsedURL.RawQuery = query.Encode()
		processedURLs[parsedURL.String()] = true
	}
	return processedURLs
}

func EnsureUniqueURLs(validURLs map[string]bool, allURLs []struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
}) map[string]bool {
	uniqueURLs := make(map[string]bool)
	for _, urlObj := range allURLs {
		if validURLs[urlObj.URL] {
			uniqueURLs[urlObj.URL] = true
		}
	}
	return uniqueURLs
}
