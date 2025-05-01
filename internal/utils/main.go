package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
)

// HandleError handles errors by logging them with context.
func HandleError(err error, context string) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", context, err)
	}
}

func ReadJSONFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if readErr := json.Unmarshal(data, v); readErr != nil {
		return fmt.Errorf("failed to parse JSON from file %s: %w", filePath, readErr)
	}
	return nil
}

func CreateDirectoryIfNotExists(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	return nil
}

func WriteJSONFile(filePath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON to file %s: %w", filePath, err)
	}
	dir := "dist"
	if dirErr := CreateDirectoryIfNotExists(dir); dirErr != nil {
		return dirErr
	}

	if writeErr := os.WriteFile(filePath, data, 0600); writeErr != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, writeErr)
	}
	return nil
}

func IsValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return false
	}
	return true
}

func CompileIgnoreRegex() (*regexp.Regexp, error) {
	ignorePattern := os.Getenv("IMPORT_IGNORE")
	if ignorePattern == "" {
		return nil, errors.New("no valid data found")
	}
	compiledRegex, err := regexp.Compile(ignorePattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	return compiledRegex, nil
}
