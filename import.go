package main

import (
	"log"
	"os"
)

func processImport(importInputFilePath, importOutputFilePath string, validateHead bool) {
	var input struct {
		Messages []struct {
			Date         string `json:"date"`
			TextEntities []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"text_entities"`
		} `json:"messages"`
	}
	err := readJSONFile(importInputFilePath, &input)
	handleError(err, "Failed to read and parse input JSON file")

	allURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{}

	idCounter := 1
	for _, message := range input.Messages {
		for _, entity := range message.TextEntities {
			if os.Getenv("DEBUG") == "true" {
				log.Printf("Processing entity: %+v", entity)
			}
			if entity.Type == "link" {
				allURLs = append(allURLs, struct {
					ID   int    `json:"id"`
					Date string `json:"date"`
					URL  string `json:"url"`
				}{
					ID:   idCounter,
					Date: message.Date,
					URL:  entity.Text,
				})
				idCounter++
			}
		}
	}

	ignoreRegex, err := compileIgnoreRegex()
	handleError(err, "Failed to compile ignore regex")

	validURLs, ignoredCount := validateURLsConcurrently(
		func() []string {
			urls := make([]string, len(allURLs))
			for i, urlObj := range allURLs {
				urls[i] = urlObj.URL
			}
			return urls
		}(),
		validateHead,
		ignoreRegex,
	)

	validURLs = removeSessionQueryStrings(validURLs)
	validURLs = ensureUniqueURLs(validURLs, allURLs)

	// Log statistics
	totalURLs := len(allURLs)
	invalidURLs := totalURLs - len(validURLs) - ignoredCount
	log.Printf("Total URLs read: %d", totalURLs)
	log.Printf("Valid URLs: %d", len(validURLs))
	log.Printf("Invalid URLs: %d", invalidURLs)
	log.Printf("Ignored URLs: %d", ignoredCount)

	filteredURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{}
	for _, urlObj := range allURLs {
		if validURLs[urlObj.URL] {
			filteredURLs = append(filteredURLs, urlObj)
		}
	}

	err = writeJSONFile(importOutputFilePath, filteredURLs)
	handleError(err, "Failed to write output JSON file")

	log.Printf("URLs successfully processed and saved to %s", importOutputFilePath)
}
