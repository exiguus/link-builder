package utils_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"urls-processor/internal/utils"
)

func TestUtils(t *testing.T) {
	t.Run("HandleError", func(t *testing.T) {
		t.Run("NoError", func(t *testing.T) {
			utils.HandleError(nil, "This should not fail")
		})
	})

	t.Run("ReadJSONFile", func(t *testing.T) {
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
	})

	t.Run("WriteJSONFile", func(t *testing.T) {
		mockOutputFile := filepath.Join(t.TempDir(), "mock_output.json")
		defer os.Remove(mockOutputFile)

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
	})

	t.Run("IsValidURL", func(t *testing.T) {
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
	})

	t.Run("CompileIgnoreRegex", func(t *testing.T) {
		t.Run("NoPattern", func(t *testing.T) {
			os.Setenv("IMPORT_IGNORE", "")
			regex, err := utils.CompileIgnoreRegex()
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if regex != nil {
				t.Errorf("Expected nil regex, got: %v", regex)
			}
		})

		t.Run("ValidPattern", func(t *testing.T) {
			os.Setenv("IMPORT_IGNORE", "^test.*$")
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
	})
}

func setupMockFiles(t *testing.T, inputData string, inputFileName string) string {
	tempFile, err := ioutil.TempFile("", inputFileName)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tempFile.Write([]byte(inputData)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	return tempFile.Name()
}
