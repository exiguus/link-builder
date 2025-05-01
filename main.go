package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	log.Println("Starting the URL Processor program")

	importInputFilePath := flag.String("import-input", "import/export.json", "Path to the input JSON file for import/export")
	importOutputFilePath := flag.String("import-output", "dist/urls.json", "Path to the output JSON file for import/export")
	validateHead := flag.Bool("validate-head", false, "Enable HEAD requests to validate URLs")

	previewInputFilePath := flag.String("preview-input", "dist/urls.json", "Path to the input JSON file containing URLs for previews")
	previewOutputFilePath := flag.String("preview-output", "dist/previews.json", "Path to the output JSON file for link previews")
	generatePreviews := flag.Bool("generate-previews", false, "Generate link previews from URLs")

	flag.Parse()

	if *generatePreviews {
		generateLinkPreviews(*previewInputFilePath, *previewOutputFilePath)
		log.Println("URL Processor program completed successfully")
		return
	}

	var input struct {
		Messages []struct {
			Date         string `json:"date"`
			TextEntities []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"text_entities"`
		} `json:"messages"`
	}
	err := readJSONFile(*importInputFilePath, &input)
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
		*validateHead,
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

	err = writeJSONFile(*importOutputFilePath, filteredURLs)
	handleError(err, "Failed to write output JSON file")

	log.Printf("URLs successfully processed and saved to %s", *importOutputFilePath)
	log.Println("URL Processor program completed successfully")
}
