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
