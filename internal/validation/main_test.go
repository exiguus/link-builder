package validation_test

import (
	"regexp"
	"testing"

	"link-builder/internal/validation"
)

const (
	exampleCom  = "http://example.com"
	exampleOrg  = "http://example.org"
	exampleDate = "2025-05-01"
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
		exampleCom,
		"https://example.com",
		"invalid-url",
	}

	ignoreRegex := regexp.MustCompile("^https://.*$")
	validURLs, ignoredCount := validation.ValidateURLsConcurrently(urls, ignoreRegex)

	if len(validURLs) != 1 || !validURLs[exampleCom] {
		t.Errorf("Unexpected valid URLs: %+v", validURLs)
	}

	if ignoredCount != 1 {
		t.Errorf("Expected 1 ignored URL, got %d", ignoredCount)
	}
}

func TestProcessURLs(t *testing.T) {
	validURLs := map[string]bool{
		exampleCom + ";jsessionid=12345": true,
		exampleCom + "?sessionid=abc":    true,
		exampleCom:                       true,
	}

	processedURLs := validation.ProcessURLs(validURLs)

	if len(processedURLs) != 1 {
		t.Errorf("Expected 1 processed URL, got %d", len(processedURLs))
	}

	if _, exists := processedURLs[exampleCom]; !exists {
		t.Errorf("Expected '%s' in processed URLs", exampleCom)
	}
}

func TestEnsureUniqueURLs(t *testing.T) {
	validURLs := map[string]bool{
		exampleCom: true,
		exampleOrg: true,
	}

	allURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{
		{ID: 1, Date: exampleDate, URL: exampleCom},
		{ID: 2, Date: exampleDate, URL: exampleOrg},
		{ID: 3, Date: exampleDate, URL: "http://example.net"},
	}

	uniqueURLs := validation.EnsureUniqueURLs(validURLs, allURLs)

	if len(uniqueURLs) != 2 {
		t.Errorf("Expected 2 unique URLs, got %d", len(uniqueURLs))
	}

	if !uniqueURLs[exampleCom] || !uniqueURLs[exampleOrg] {
		t.Errorf("Expected '%s' and '%s' in unique URLs", exampleCom, exampleOrg)
	}
}
