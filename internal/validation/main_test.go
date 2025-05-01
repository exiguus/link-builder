package validation_test

import (
	"link-builder/internal/validation"
	"regexp"
	"testing"
)

func TestIgnoreRegex(t *testing.T) {
	t.Run("ConstPattern", func(t *testing.T) {
		ignoreRegex := regexp.MustCompile("^http://ignored.com$")
		if !ignoreRegex.MatchString("http://ignored.com") {
			t.Errorf("Expected regex to match 'http://ignored.com'")
		}
	})

	t.Run("DynamicPattern", func(t *testing.T) {
		ignoreRegex := regexp.MustCompile("^https://.*$")
		if !ignoreRegex.MatchString("https://example.com") {
			t.Errorf("Expected regex to match 'https://example.com'")
		}
	})
}

func TestValidateURLsConcurrently(t *testing.T) {
	urls := []string{
		"http://example.com",
		"https://example.com",
		"invalid-url",
	}

	ignoreRegex := regexp.MustCompile("^https://.*$")
	validURLs, ignoredCount := validation.ValidateURLsConcurrently(urls, ignoreRegex)

	if len(validURLs) != 1 || !validURLs["http://example.com"] {
		t.Errorf("Unexpected valid URLs: %+v", validURLs)
	}

	if ignoredCount != 1 {
		t.Errorf("Expected 1 ignored URL, got %d", ignoredCount)
	}
}

func TestProcessURLs(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com;jsessionid=12345": true,
		"http://example.com?sessionid=abc":    true,
		"http://example.com":                  true,
	}

	processedURLs := validation.ProcessURLs(validURLs)

	if len(processedURLs) != 1 {
		t.Errorf("Expected 1 processed URL, got %d", len(processedURLs))
	}

	if _, exists := processedURLs["http://example.com"]; !exists {
		t.Errorf("Expected 'http://example.com' in processed URLs")
	}
}

func TestEnsureUniqueURLs(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
	}

	allURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{
		{ID: 1, Date: "2025-05-01", URL: "http://example.com"},
		{ID: 2, Date: "2025-05-01", URL: "http://example.org"},
		{ID: 3, Date: "2025-05-01", URL: "http://example.net"},
	}

	uniqueURLs := validation.EnsureUniqueURLs(validURLs, allURLs)

	if len(uniqueURLs) != 2 {
		t.Errorf("Expected 2 unique URLs, got %d", len(uniqueURLs))
	}

	if !uniqueURLs["http://example.com"] || !uniqueURLs["http://example.org"] {
		t.Errorf("Expected 'http://example.com' and 'http://example.org' in unique URLs")
	}
}
