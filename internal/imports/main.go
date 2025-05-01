package imports

import (
	"fmt"
	"link-builder/internal/utils"
	"link-builder/internal/validation"
	"log"
	"os"
)

func ProcessImport(importInputFilePath, importOutputFilePath string) error {
	var input struct {
		Messages []struct {
			Date         string `json:"date"`
			TextEntities []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"text_entities"`
		} `json:"messages"`
	}
	err := utils.ReadJSONFile(importInputFilePath, &input)
	if err != nil {
		return fmt.Errorf("reading and parsing input JSON file: %w", err)
	}

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

	ignoreRegex, err := utils.CompileIgnoreRegex()
	if err != nil {
		ignoreRegex = nil
	}

	validURLs, ignoredCount := validation.ValidateURLsConcurrently(
		func() []string {
			urls := make([]string, len(allURLs))
			for i, urlObj := range allURLs {
				urls[i] = urlObj.URL
			}
			return urls
		}(),
		ignoreRegex,
	)

	validURLs = validation.ProcessURLs(validURLs)
	validURLs = validation.EnsureUniqueURLs(validURLs, allURLs)

	// Log statistics
	totalURLs := len(allURLs)
	invalidURLs := totalURLs - len(validURLs) - ignoredCount
	logStatistics(totalURLs, len(validURLs), invalidURLs, ignoredCount)

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

	err = utils.WriteJSONFile(importOutputFilePath, filteredURLs)
	if err != nil {
		return fmt.Errorf("writing output JSON file: %w", err)
	}

	log.Printf("URLs successfully processed and saved to %s", importOutputFilePath)
	return nil
}

func logStatistics(totalURLs, validURLsCount, invalidURLs, ignoredCount int) {
	log.Printf("Total URLs read: %d", totalURLs)
	log.Printf("Valid URLs: %d", validURLsCount)
	log.Printf("Invalid URLs: %d", invalidURLs)
	log.Printf("Ignored URLs: %d", ignoredCount)
}
