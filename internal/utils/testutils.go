package utils

import (
	"io/ioutil"
	"testing"
)

// CreateTempFile creates a temporary file with the given content and returns its path.
func CreateTempFile(t *testing.T, content string, fileName string) string {
	tempFile, err := ioutil.TempFile("", fileName)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	return tempFile.Name()
}
