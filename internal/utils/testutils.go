package utils

import (
	"io/ioutil"
	"testing"
)

// CreateTempFile creates a temporary file with the given content and name for testing purposes.
func CreateTempFile(t *testing.T, content, name string) string {
	t.Helper()
	tempFile, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	return tempFile.Name()
}
