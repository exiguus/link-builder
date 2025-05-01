package imports_test

import (
	"os"
	"testing"

	"link-builder/internal/imports"
	"link-builder/internal/utils"
)

func TestProcessImport(t *testing.T) {
	t.Setenv("IMPORT_IGNORE", "^https?://excluded\\.com$") // Update regex to exclude a different URL

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
	if err = utils.ReadJSONFile(tempOutputFile, &result); err != nil {
		t.Errorf("Failed to read output JSON file: %v", err)
	}

	if len(result) != 1 || result[0].URL != "http://example.com" {
		t.Errorf("Unexpected result: %+v", result)
	}
}

func TestProcessImport_EdgeCases(t *testing.T) {
	// Test with invalid JSON
	invalidJSON := `{"messages": [ { "date": "2025-05-01", "text_entities": [ { "type": "link", "text": "http://example.com" } ] }`
	tempInputFile := utils.CreateTempFile(t, invalidJSON, "invalid_import_input.json")
	defer os.Remove(tempInputFile)

	tempOutputFile := utils.CreateTempFile(t, "", "invalid_import_output.json")
	defer os.Remove(tempOutputFile)

	err := imports.ProcessImport(tempInputFile, tempOutputFile)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}

	// Test with missing fields
	missingFieldsJSON := `{"messages": [ { "text_entities": [ { "type": "link", "text": "http://example.com" } ] } ]}`
	tempInputFile = utils.CreateTempFile(t, missingFieldsJSON, "missing_fields_import_input.json")
	defer os.Remove(tempInputFile)

	err = imports.ProcessImport(tempInputFile, tempOutputFile)
	if err == nil {
		t.Errorf("Expected error for missing fields, got nil")
	}

	// Test with empty input file
	tempInputFile = utils.CreateTempFile(t, "", "empty_import_input.json")
	defer os.Remove(tempInputFile)

	err = imports.ProcessImport(tempInputFile, tempOutputFile)
	if err == nil {
		t.Errorf("Expected error for empty input file, got nil")
	}
}
