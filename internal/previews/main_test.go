package previews_test

import (
	"errors"
	"os"
	"testing"

	"link-builder/internal/previews"
	"link-builder/internal/types"
	"link-builder/internal/utils"
)

type MockLinkPreviewer struct{}

var ErrNoValidPreview = errors.New("no valid preview generated")

const exampleComURL = "http://example.com"

func (m MockLinkPreviewer) Parse(_ string) (*previews.Preview, error) {
	return nil, ErrNoValidPreview // Return a sentinel error instead of nil, nil
}

func TestGenerateLinkPreviews(t *testing.T) {
	mockInput := `[{"id": 1, "date": "2025-05-01", "url": "` + exampleComURL + `"}]`
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

	if len(result) != 1 || result[0].URL != exampleComURL {
		t.Errorf("Unexpected result: %+v", result)
	}
}

func TestGenerateLinkPreviewsEdgeCases(t *testing.T) {
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

func TestParseInputFile(t *testing.T) {
	// Test with valid input
	validInput := `[{"id": 1, "date": "2025-05-01", "url": "` + exampleComURL + `"}]`
	tempFile := utils.CreateTempFile(t, validInput, "valid_input.json")
	defer os.Remove(tempFile)

	result, err := previews.ParseInputFile(tempFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result) != 1 || result[0].URL != exampleComURL {
		t.Errorf("Unexpected result: %+v", result)
	}

	// Test with invalid JSON
	invalidInput := `[{"id": 1, "date": "2025-05-01", "url": "` + exampleComURL + `"`
	tempFile = utils.CreateTempFile(t, invalidInput, "invalid_input.json")
	defer os.Remove(tempFile)

	_, err = previews.ParseInputFile(tempFile)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}

	// Test with missing fields
	missingFieldsInput := `[{"id": 1, "url": "` + exampleComURL + `"}]`
	tempFile = utils.CreateTempFile(t, missingFieldsInput, "missing_fields_input.json")
	defer os.Remove(tempFile)

	_, err = previews.ParseInputFile(tempFile)
	if err == nil {
		t.Errorf("Expected error for missing fields, got nil")
	}
}

func TestLoadCache(t *testing.T) {
	// Test with non-existent file
	nonExistentFile := "non_existent.json"
	t.Logf("Testing non-existent file: %s", nonExistentFile)

	// Ensure the file does not exist before the test
	if _, err := os.Stat(nonExistentFile); err == nil {
		os.Remove(nonExistentFile)
	}

	_, err := previews.LoadCache(nonExistentFile)
	if err != nil {
		t.Errorf("Expected no error for non-existent file, got: %v", err)
	}

	// Verify the file was created
	_, statErr := os.Stat(nonExistentFile)
	if os.IsNotExist(statErr) {
		t.Errorf("Expected file to be created, but it does not exist: %s", nonExistentFile)
	}
	defer os.Remove(nonExistentFile)

	// Test with empty file
	tempFile := utils.CreateTempFile(t, "", "empty_cache.json")
	defer os.Remove(tempFile)

	cache, err := previews.LoadCache(tempFile)
	if err != nil {
		t.Errorf("Expected no error for empty file, got: %v", err)
	}
	if len(cache) != 0 {
		t.Errorf("Expected empty cache, got: %+v", cache)
	}

	// Test with invalid JSON
	invalidContent := "invalid-json"
	tempFile = utils.CreateTempFile(t, invalidContent, "invalid_cache.json")
	defer os.Remove(tempFile)

	_, err = previews.LoadCache(tempFile)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}
}

func TestSaveOutput(t *testing.T) {
	// Test with valid output
	output := []types.LinkPreviewOutput{
		{ID: 1, Date: "2025-05-01", URL: exampleComURL, Preview: map[string]interface{}{"title": "Example"}},
	}
	tempFile := utils.CreateTempFile(t, "", "valid_output.json")
	defer os.Remove(tempFile)

	err := previews.SaveOutput(tempFile, output)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Test with invalid file path
	invalidFilePath := "/invalid/path/output.json"
	err = previews.SaveOutput(invalidFilePath, output)
	if err == nil {
		t.Errorf("Expected error for invalid file path, got nil")
	}
}
