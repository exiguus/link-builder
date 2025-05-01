package main

import (
	"context"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/time/rate"
)

func validateURLsConcurrently(urls []string, validateHead bool, ignoreRegex *regexp.Regexp) (map[string]bool, int) {
	validURLs := make(map[string]bool)
	var mutex sync.Mutex
	var ignoredCount int32

	rateLimiter := rate.NewLimiter(rate.Limit(10), 1)
	processURL := func(rawURL string) {
		if ignoreRegex != nil && ignoreRegex.MatchString(rawURL) {
			atomic.AddInt32(&ignoredCount, 1)
			return
		}
		if validateHead {
			if err := rateLimiter.Wait(context.Background()); err != nil {
				log.Printf("Rate limiter error: %v", err)
				return
			}
		}
		if isValidURL(rawURL, validateHead) {
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

func removeSessionQueryStrings(validURLs map[string]bool) map[string]bool {
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

func warnIfURLsContainSession(validURLs map[string]bool) {
	for url := range validURLs {
		if strings.Contains(strings.ToLower(url), "session") {
			log.Printf("Warning: URL contains 'session': %s", url)
		}
	}
}

func ensureUniqueURLs(validURLs map[string]bool, allURLs []struct {
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
