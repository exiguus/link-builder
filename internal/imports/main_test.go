package imports_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"link-builder/internal/imports"
	"link-builder/internal/utils"
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

func TestProcessImport(t *testing.T) {
	mockInput := `{"messages": [{"date": "2025-05-01", "text_entities": [{"type": "link", "text": "http://example.com"}]}]}`
	tempInputFile := utils.CreateTempFile(t, mockInput, "mock_import_input.json")
	defer os.Remove(tempInputFile)

	tempOutputFile := utils.CreateTempFile(t, "", "mock_import_output.json")
	defer os.Remove(tempOutputFile)

	err := imports.ProcessImport(tempInputFile, tempOutputFile)
	if err != nil {
		t.Errorf("ProcessImport failed: %v", err)
	}

	var result []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	if err := utils.ReadJSONFile(tempOutputFile, &result); err != nil {
		t.Errorf("Failed to read output JSON file: %v", err)
	}

	if len(result) != 1 || result[0].URL != "http://example.com" {
		t.Errorf("Unexpected result: %+v", result)
	}
}
