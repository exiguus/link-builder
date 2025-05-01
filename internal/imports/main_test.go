package imports_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"urls-processor/internal/imports"
	"urls-processor/internal/utils"
)

func setupMockImportFiles(t *testing.T, inputData string) (string, string) {
	inputFilePath := filepath.Join(t.TempDir(), "mock_input.json")
	outputFilePath := filepath.Join(t.TempDir(), "mock_output.json")

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	return inputFilePath, outputFilePath
}

func validateOutput(t *testing.T, outputFilePath string, expectedURLs map[string]bool) {
	var output []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	if err := utils.ReadJSONFile(outputFilePath, &output); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	for _, entry := range output {
		if entry.URL == "" {
			t.Errorf("Empty URL found in output: %+v", entry)
			continue
		}
		if !expectedURLs[entry.URL] {
			t.Errorf("Unexpected URL in output: %s", entry.URL)
		}
		delete(expectedURLs, entry.URL)
	}

	if len(expectedURLs) > 0 {
		t.Errorf("Missing expected URLs in output: %+v", expectedURLs)
	}
}

func TestProcessImportWithMocks(t *testing.T) {
	mockInput := `{
		"messages": [
			{
				"date": "2025-05-01",
				"text_entities": [
					{"type": "link", "text": "http://example.com"},
					{"type": "link", "text": "https://example.org"}
				]
			}
		]
	}`

	inputFilePath := utils.CreateTempFile(t, mockInput, "mock_import_input.json")
	defer os.Remove(inputFilePath)

	outputFilePath := filepath.Join(t.TempDir(), "mock_import_output.json")
	defer os.Remove(outputFilePath)

	imports.ProcessImport(inputFilePath, outputFilePath)

	expectedURLs := map[string]bool{
		"http://example.com":  true,
		"https://example.org": true,
	}

	validateOutput(t, outputFilePath, expectedURLs)
}
