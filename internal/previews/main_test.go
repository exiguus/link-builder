package previews_test

import (
	"errors"
	"os"
	"testing"

	"link-builder/internal/previews"
	"link-builder/internal/utils"
)

type MockLinkPreviewer struct{}

var ErrNoValidPreview = errors.New("no valid preview generated")

func (m MockLinkPreviewer) Parse(_ string) (*previews.Preview, error) {
	return nil, ErrNoValidPreview // Return a sentinel error instead of nil, nil
}

func TestGenerateLinkPreviews(t *testing.T) {
	mockInput := `[{"id": 1, "date": "2025-05-01", "url": "http://example.com"}]`
	tempInputFile := utils.CreateTempFile(t, mockInput, "mock_preview_input.json")
	defer os.Remove(tempInputFile)

	tempOutputFile := utils.CreateTempFile(t, "", "mock_preview_output.json")
	defer os.Remove(tempOutputFile)

	err := previews.GenerateLinkPreviews(tempInputFile, tempOutputFile, previews.DefaultLinkPreviewer{})
	if err != nil {
		t.Errorf("GenerateLinkPreviews failed: %v", err)
	}

	var result []struct {
		ID      int         `json:"id"`
		Date    string      `json:"date"`
		URL     string      `json:"url"`
		Preview interface{} `json:"preview"`
	}
	if readErr := utils.ReadJSONFile(tempOutputFile, &result); readErr != nil {
		t.Errorf("Failed to read output JSON file: %v", readErr)
	}

	if len(result) != 1 || result[0].URL != "http://example.com" {
		t.Errorf("Unexpected result: %+v", result)
	}
}

func TestGenerateLinkPreviews_EdgeCases(t *testing.T) {
	mockPreviewer := previews.DefaultLinkPreviewer{}

	// Test with empty input file
	t.Run("EmptyInputFile", func(t *testing.T) {
		mockInputFile := utils.CreateTempFile(t, "", "empty_preview_input.json")
		defer os.Remove(mockInputFile)

		mockOutputFile := utils.CreateTempFile(t, "", "empty_preview_output.json")
		defer os.Remove(mockOutputFile)

		err := previews.GenerateLinkPreviews(mockInputFile, mockOutputFile, mockPreviewer)
		if err == nil {
			t.Errorf("Expected error for empty input file, got nil")
		}
	})

	// Test with invalid URLs
	t.Run("InvalidURLs", func(t *testing.T) {
		mockInput := `[
			{"id": 1, "date": "2025-05-01", "url": "ftp://example.com"},
			{"id": 2, "date": "2025-05-01", "url": "example.com"}
		]`
		mockInputFile := utils.CreateTempFile(t, mockInput, "invalid_urls_input.json")
		defer os.Remove(mockInputFile)

		mockOutputFile := utils.CreateTempFile(t, "", "invalid_urls_output.json")
		defer os.Remove(mockOutputFile)

		err := previews.GenerateLinkPreviews(mockInputFile, mockOutputFile, mockPreviewer)
		if err == nil {
			t.Errorf("Expected error for invalid URLs, got nil")
		}
	})

	// Test with no valid previews generated
	t.Run("NoValidPreviews", func(t *testing.T) {
		mockInput := `[
			{"id": 1, "date": "2025-05-01", "url": "http://invalid-url.com"}
		]`
		mockInputFile := utils.CreateTempFile(t, mockInput, "no_valid_previews_input.json")
		defer os.Remove(mockInputFile)

		mockOutputFile := utils.CreateTempFile(t, "", "no_valid_previews_output.json")
		defer os.Remove(mockOutputFile)

		mockPreviewer := MockLinkPreviewer{}
		err := previews.GenerateLinkPreviews(mockInputFile, mockOutputFile, mockPreviewer)
		if err == nil {
			t.Errorf("Expected error for no valid previews, got nil")
		}
	})
}
