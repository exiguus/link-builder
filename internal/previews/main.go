package previews

import (
	"encoding/json"
	"log"
	"os"
	"urls-processor/internal/types"
	"urls-processor/internal/utils"

	"github.com/tiendc/go-linkpreview"
)

func GenerateLinkPreviews(inputFilePath, outputFilePath string) {
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		utils.HandleError(err, "Failed to read input file")
	}

	if !json.Valid(data) {
		log.Fatalf("Input JSON is invalid: %s", inputFilePath)
	}

	var urlObjects []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	utils.HandleError(json.Unmarshal(data, &urlObjects), "Failed to parse input JSON")

	cache := make(map[string]interface{})
	if _, err := os.Stat(outputFilePath); err == nil {
		cacheData, err := os.ReadFile(outputFilePath)
		utils.HandleError(err, "Failed to read output file")
		if len(cacheData) > 0 {
			if err := json.Unmarshal(cacheData, &cache); err != nil {
				var cacheArray []types.LinkPreviewOutput
				if err := json.Unmarshal(cacheData, &cacheArray); err == nil {
					for _, item := range cacheArray {
						cache[item.URL] = item.Preview
					}
				} else if string(cacheData) == "[]" {
					cache = make(map[string]interface{})
				} else {
					utils.HandleError(err, "Failed to parse output JSON")
				}
			}
		}
	}

	totalURLs := len(urlObjects)
	cachedCount := 0
	for _, urlObj := range urlObjects {
		if _, exists := cache[urlObj.URL]; exists {
			cachedCount++
		}
	}
	toProcessCount := totalURLs - cachedCount
	log.Printf("Total URLs: %d, Cached: %d, To Process: %d", totalURLs, cachedCount, toProcessCount)

	output := []types.LinkPreviewOutput{}
	currentCount := 0
	for _, urlObj := range urlObjects {
		currentCount++
		log.Printf("Processing URL %d/%d: %s", currentCount, totalURLs, urlObj.URL)

		preview, exists := cache[urlObj.URL]
		if !exists {
			parsedPreview, err := linkpreview.Parse(urlObj.URL)
			if err != nil {
				log.Printf("Failed to generate preview for %s: %v", urlObj.URL, err)
				continue
			}

			if parsedPreview.Title == "" && parsedPreview.Description == "" && (parsedPreview.OGMeta == nil && parsedPreview.TwitterMeta == nil) {
				log.Printf("Skipping invalid preview for %s", urlObj.URL)
				continue
			}

			preview = map[string]interface{}{
				"title":        parsedPreview.Title,
				"description":  parsedPreview.Description,
				"og_meta":      parsedPreview.OGMeta,
				"twitter_meta": parsedPreview.TwitterMeta,
			}
			cache[urlObj.URL] = preview
		}

		output = append(output, types.LinkPreviewOutput{
			ID:      urlObj.ID,
			Date:    urlObj.Date,
			URL:     urlObj.URL,
			Preview: preview,
		})

		outputData, err := json.MarshalIndent(output, "", "  ")
		utils.HandleError(err, "Failed to marshal intermediate output JSON")
		utils.HandleError(os.WriteFile(outputFilePath, outputData, 0644), "Failed to write intermediate output file")
	}

	log.Printf("Link previews successfully generated and saved to %s", outputFilePath)
}
