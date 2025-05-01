package main

import (
	"os"
	"os/exec"
	"testing"
)

func setupMockFiles(t *testing.T, inputData string, inputFileName string) string {
	tempFile, err := os.CreateTemp(t.TempDir(), inputFileName)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, writeErr := tempFile.WriteString(inputData); writeErr != nil {
		t.Fatalf("Failed to write to temporary file: %v", writeErr)
	}
	tempFile.Close()

	return tempFile.Name()
}

func TestMainProgram(t *testing.T) {
	t.Run("ImportURLs", func(t *testing.T) {
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
		mockInputFile := setupMockFiles(t, mockInput, "mock_import_input.json")
		defer os.Remove(mockInputFile)

		mockOutputFile := "mock_import_output.json"
		defer os.Remove(mockOutputFile)

		cmd := exec.Command(
			"go",
			"run",
			".",
			"-import-urls",
			"-import-input="+mockInputFile,
			"-import-output="+mockOutputFile,
		)
		cmd.Env = append(os.Environ(), "DEBUG=true")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run import-urls: %v\nOutput: %s", err, string(output))
		}
	})

	t.Run("GeneratePreviews", func(t *testing.T) {
		mockInput := `[
			{"id": 1, "date": "2025-05-01", "url": "http://example.com"},
			{"id": 2, "date": "2025-05-01", "url": "https://example.org"}
		]`
		mockInputFile := setupMockFiles(t, mockInput, "mock_preview_input.json")
		defer os.Remove(mockInputFile)

		mockOutputFile := "mock_preview_output.json"
		defer os.Remove(mockOutputFile)

		cmd := exec.Command(
			"go",
			"run",
			".",
			"-generate-preview",
			"-preview-input="+mockInputFile,
			"-preview-output="+mockOutputFile,
		)
		cmd.Env = append(os.Environ(), "DEBUG=true")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run generate-preview: %v\nOutput: %s", err, string(output))
		}
	})
}
