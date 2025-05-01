package previews_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"urls-processor/internal/previews"
	"urls-processor/internal/utils"
)

func setupPreviewTestFiles(t *testing.T, inputData string, inputFileName string) (string, string) {
	tempDir := t.TempDir()
	inputFilePath := filepath.Join(tempDir, inputFileName)
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	return inputFilePath, outputFilePath
}

func validateJSON(t *testing.T, jsonData string) {
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonData), &temp); err != nil {
		t.Fatalf("Invalid JSON provided: %v", err)
	}
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
	if err := utils.ReadJSONFile(tempOutputFile, &result); err != nil {
		t.Errorf("Failed to read output JSON file: %v", err)
	}

	if len(result) != 1 || result[0].URL != "http://example.com" {
		t.Errorf("Unexpected result: %+v", result)
	}
}
