package utils

import (
	"os"
	"testing"
)

// CreateTempFile creates a temporary file with the given content and name for testing purposes.
func CreateTempFile(t *testing.T, content, name string) string {
	t.Helper()
	tempFile, err := os.CreateTemp(t.TempDir(), name)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, writeErr := tempFile.WriteString(content); writeErr != nil {
		t.Fatalf("Failed to write to temporary file: %v", writeErr)
	}
	if closeErr := tempFile.Close(); closeErr != nil {
		t.Fatalf("Failed to close temporary file: %v", closeErr)
	}

	return tempFile.Name()
}
