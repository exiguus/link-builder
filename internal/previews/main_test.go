package previews_test

import (
	"os"
	"testing"

	"link-builder/internal/previews"
	"link-builder/internal/utils"
)

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
