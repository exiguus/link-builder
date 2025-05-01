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

func TestReadJSONFile_Errors(t *testing.T) {
	// Test with a non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		var result map[string]string
		err := utils.ReadJSONFile("non_existent_file.json", &result)
		if err == nil {
			t.Errorf("Expected error for non-existent file, got nil")
		}
	})

	// Test with invalid JSON content
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidContent := `{"key": "value"`
		tempFile := utils.CreateTempFile(t, invalidContent, "invalid_json.json")
		defer os.Remove(tempFile)

		var result map[string]string
		err := utils.ReadJSONFile(tempFile, &result)
		if err == nil {
			t.Errorf("Expected error for invalid JSON, got nil")
		}
	})
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

func TestWriteJSONFile_Errors(t *testing.T) {
	// Test with a directory path instead of a file path
	t.Run("DirectoryPath", func(t *testing.T) {
		dirPath := t.TempDir()
		data := map[string]string{"key": "value"}
		err := utils.WriteJSONFile(dirPath, data)
		if err == nil {
			t.Errorf("Expected error for directory path, got nil")
		}
	})

	// Test with data that cannot be marshaled into JSON
	t.Run("UnmarshalableData", func(t *testing.T) {
		mockOutputFile := filepath.Join(t.TempDir(), "mock_output.json")
		data := map[string]interface{}{"key": make(chan int)} // Channels cannot be marshaled into JSON
		err := utils.WriteJSONFile(mockOutputFile, data)
		if err == nil {
			t.Errorf("Expected error for unmarshalable data, got nil")
		}
	})
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
