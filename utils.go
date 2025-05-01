package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var httpClient = &http.Client{}

func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func readJSONFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON from file %s: %w", filePath, err)
	}
	return nil
}

func writeJSONFile(filePath string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON to file %s: %w", filePath, err)
	}
	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}

func isValidURL(rawURL string, validateHead bool) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return false
	}
	if validateHead {
		resp, err := httpClient.Head(rawURL)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return false
		}
	}
	return true
}

func compileIgnoreRegex() (*regexp.Regexp, error) {
	ignorePattern := os.Getenv("IMPORT_IGNORE")
	if ignorePattern == "" {
		return nil, nil
	}
	return regexp.Compile(ignorePattern)
}
