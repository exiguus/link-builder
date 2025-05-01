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

func TestPreviews(t *testing.T) {
	t.Run("GenerateLinkPreviews", func(t *testing.T) {
		mockInput := `[
			{"id": 1, "date": "2025-05-01", "url": "http://example.com"},
			{"id": 2, "date": "2025-05-01", "url": "https://example.org"}
		]`

		inputFilePath := utils.CreateTempFile(t, mockInput, "mock_preview_input.json")
		defer os.Remove(inputFilePath)

		outputFilePath := filepath.Join(t.TempDir(), "mock_preview_output.json")
		defer os.Remove(outputFilePath)

		previews.GenerateLinkPreviews(inputFilePath, outputFilePath)

		outputData, err := ioutil.ReadFile(outputFilePath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		if len(outputData) == 0 {
			t.Errorf("Output file is empty")
		}

		var output []struct {
			ID      int         `json:"id"`
			Date    string      `json:"date"`
			URL     string      `json:"url"`
			Preview interface{} `json:"preview"`
		}
		if err := utils.ReadJSONFile(outputFilePath, &output); err != nil {
			t.Fatalf("Failed to parse output JSON: %v", err)
		}

		expectedCount := 2
		if len(output) != expectedCount {
			t.Errorf("Expected %d previews, got %d", expectedCount, len(output))
		}
	})
}
