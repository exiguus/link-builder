package validation_test

import (
	"regexp"
	"testing"

	"urls-processor/internal/validation"
)

func TestValidateURLsConcurrentlyWithMocks(t *testing.T) {
	mockURLs := []string{
		"http://example.com",
		"https://example.org",
		"http://ignored.com",
	}

	ignoreRegex, _ := regexp.Compile("^http://ignored.com$")
	validURLs, ignoredCount := validation.ValidateURLsConcurrently(mockURLs, ignoreRegex)

	expectedValid := map[string]bool{
		"http://example.com":  true,
		"https://example.org": true,
	}

	if len(validURLs) != len(expectedValid) {
		t.Errorf("Expected %d valid URLs, got %d", len(expectedValid), len(validURLs))
	}

	for url := range expectedValid {
		if !validURLs[url] {
			t.Errorf("Expected URL %s to be valid, but it was not", url)
		}
	}

	if ignoredCount != 1 {
		t.Errorf("Expected 1 ignored URL, got %d", ignoredCount)
	}
}

func TestRemoveSessionQueryStrings(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com;jsessionid=12345": true,
		"http://example.org?sessionid=67890":  true,
		"http://example.net":                  true,
	}

	updatedURLs := validation.RemoveSessionQueryStrings(validURLs)

	expectedURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
		"http://example.net": true,
	}

	if len(updatedURLs) != len(expectedURLs) {
		t.Errorf("Expected %d URLs after removing session query strings, got %d", len(expectedURLs), len(updatedURLs))
	}

	for url := range expectedURLs {
		if !updatedURLs[url] {
			t.Errorf("Expected URL %s to be present, but it was not", url)
		}
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

	expectedURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
	}

	if len(uniqueURLs) != len(expectedURLs) {
		t.Errorf("Expected %d unique URLs, got %d", len(expectedURLs), len(uniqueURLs))
	}

	for url := range expectedURLs {
		if !uniqueURLs[url] {
			t.Errorf("Expected URL %s to be unique, but it was not", url)
		}
	}
}
