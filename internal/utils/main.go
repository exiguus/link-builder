package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
)

// Exported HandleError function
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
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON from file %s: %w", filePath, err)
	}
	return nil
}

func CreateDirectoryIfNotExists(dirPath string) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
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
	if err := CreateDirectoryIfNotExists(dir); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
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
		return nil, nil
	}
	return regexp.Compile(ignorePattern)
}
