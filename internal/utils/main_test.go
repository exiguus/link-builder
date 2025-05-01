package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"link-builder/internal/utils"
)

func TestHandleError(t *testing.T) {
	t.Run("NoError", func(_ *testing.T) {
		utils.HandleError(nil, "This should not fail")
	})
}

func TestReadJSONFile(t *testing.T) {
	mockContent := `{"key": "value"}`
	tempFile := utils.CreateTempFile(t, mockContent, "mock_read_json.json")
	defer os.Remove(tempFile)

	var result map[string]string
	if err := utils.ReadJSONFile(tempFile, &result); err != nil {
		t.Errorf("Failed to read JSON file: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("Expected 'value', got '%s'", result["key"])
	}
}

func TestWriteJSONFile(t *testing.T) {
	mockOutputFile := filepath.Join(t.TempDir(), "mock_output.json")
	data := map[string]string{"key": "value"}
	if err := utils.WriteJSONFile(mockOutputFile, data); err != nil {
		t.Errorf("Failed to write JSON file: %v", err)
	}

	var result map[string]string
	if err := utils.ReadJSONFile(mockOutputFile, &result); err != nil {
		t.Errorf("Failed to read back JSON file: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("Expected 'value', got '%s'", result["key"])
	}
}

func TestIsValidURL(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com",
	}
	invalidURLs := []string{
		"ftp://example.com",
		"example.com",
		"http://",
	}

	for _, url := range validURLs {
		if !utils.IsValidURL(url) {
			t.Errorf("Expected valid URL, got invalid: %s", url)
		}
	}

	for _, url := range invalidURLs {
		if utils.IsValidURL(url) {
			t.Errorf("Expected invalid URL, got valid: %s", url)
		}
	}
}

func TestCompileIgnoreRegex(t *testing.T) {
	t.Run("NoPattern", func(t *testing.T) {
		t.Setenv("IMPORT_IGNORE", "")
		regex, err := utils.CompileIgnoreRegex()
		if err == nil {
			t.Errorf("Expected an error, got nil")
		}
		if regex != nil {
			t.Errorf("Expected nil regex, got: %v", regex)
		}
	})

	t.Run("ValidPattern", func(t *testing.T) {
		t.Setenv("IMPORT_IGNORE", "^test.*$")
		regex, err := utils.CompileIgnoreRegex()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if regex == nil {
			t.Errorf("Expected valid regex, got nil")
		}
		if !regex.MatchString("test123") {
			t.Errorf("Expected regex to match 'test123', but it did not")
		}
	})
}
