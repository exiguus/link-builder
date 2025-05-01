// Change package name to `validation`
package validation

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"urls-processor/internal/utils"
)

func ValidateURLsConcurrently(urls []string, ignoreRegex *regexp.Regexp) (map[string]bool, int) {
	validURLs := make(map[string]bool)
	var mutex sync.Mutex
	var ignoredCount int32

	processURL := func(rawURL string) {
		if ignoreRegex != nil && ignoreRegex.MatchString(rawURL) {
			atomic.AddInt32(&ignoredCount, 1)
			return
		}
		if validateURL(rawURL) {
			mutex.Lock()
			validURLs[rawURL] = true
			mutex.Unlock()
		}
	}

	var wg sync.WaitGroup
	for _, rawURL := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			processURL(url)
		}(rawURL)
	}
	wg.Wait()

	return validURLs, int(ignoredCount)
}

func validateURL(url string) bool {
	return utils.IsValidURL(url)
}

func RemoveSessionQueryStrings(validURLs map[string]bool) map[string]bool {
	updatedURLs := make(map[string]bool)
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
			}
		}
		parsedURL.RawQuery = query.Encode()
		updatedURLs[parsedURL.String()] = true
	}
	return updatedURLs
}

func WarnIfURLsContainSession(validURLs map[string]bool) {
	for url := range validURLs {
		if strings.Contains(strings.ToLower(url), "session") {
			log.Printf("Warning: URL contains 'session': %s", url)
		}
	}
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
