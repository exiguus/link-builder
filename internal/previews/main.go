package previews

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tiendc/go-linkpreview"

	"link-builder/internal/types"
)

// Preview represents the metadata extracted from a URL.
type Preview struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	OGMeta      map[string]string `json:"og_meta,omitempty"`
	TwitterMeta map[string]string `json:"twitter_meta,omitempty"`
}

type LinkPreviewer interface {
	Parse(url string) (*Preview, error)
}

type DefaultLinkPreviewer struct{}

func (d DefaultLinkPreviewer) Parse(url string) (*Preview, error) {
	result, err := linkpreview.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing link preview: %w", err)
	}
	twitterMeta := make(map[string]string)
	if result.TwitterMeta != nil {
		twitterMeta["title"] = result.TwitterMeta.Title
		twitterMeta["description"] = result.TwitterMeta.Description
		twitterMeta["card"] = result.TwitterMeta.Card
		twitterMeta["site"] = result.TwitterMeta.Site
		twitterMeta["creator"] = result.TwitterMeta.Creator
		twitterMeta["image"] = result.TwitterMeta.Image
	}
	ogMeta := make(map[string]string)
	if result.OGMeta != nil {
		ogMeta["title"] = result.OGMeta.Title
		ogMeta["type"] = result.OGMeta.Type
		ogMeta["description"] = result.OGMeta.Description
		ogMeta["url"] = result.OGMeta.URL
		if len(result.OGMeta.Images) > 0 {
			ogMeta["image"] = result.OGMeta.Images[0].URL
		}
		ogMeta["site_name"] = result.OGMeta.SiteName
	}
	return &Preview{
		Title:       result.Title,
		Description: result.Description,
		OGMeta:      ogMeta,
		TwitterMeta: twitterMeta,
	}, nil
}

// GenerateLinkPreviews generates link previews for a list of URLs and saves the results to a file.
func GenerateLinkPreviews(inputFilePath, outputFilePath string, previewer LinkPreviewer) error {
	// Break down logic into smaller helper functions.
	urlObjects, err := parseInputFile(inputFilePath)
	if err != nil {
		return err
	}

	cache, err := loadCache(outputFilePath)
	if err != nil {
		return err
	}

	output, err := generatePreviews(urlObjects, cache, previewer, outputFilePath)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Program interrupted. Writing current output to file.")
			if writeErr := saveOutput(outputFilePath, output); writeErr != nil {
				log.Printf("Failed to write output file on interrupt: %v", writeErr)
			}
			panic(r) // Re-throw the panic after handling
		}
	}()

	if len(output) == 0 {
		log.Println("No valid previews generated. Writing empty output file.")
		if writeErr := saveOutput(outputFilePath, output); writeErr != nil {
			return fmt.Errorf("failed to write empty output file: %w", writeErr)
		}
		return errors.New("no valid previews generated")
	}

	return saveOutput(outputFilePath, output)
}

// Helper functions for GenerateLinkPreviews.
func parseInputFile(inputFilePath string) ([]struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
}, error) {
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading input file: %w", err)
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("input JSON is invalid: %s", inputFilePath)
	}

	var urlObjects []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	if err = json.Unmarshal(data, &urlObjects); err != nil {
		return nil, fmt.Errorf("parsing input JSON: %w", err)
	}

	for _, obj := range urlObjects {
		if obj.Date == "" {
			return nil, errors.New("missing required field: Date")
		}
		if obj.URL == "" {
			return nil, errors.New("missing required field: URL")
		}
	}

	return urlObjects, nil
}

func ParseInputFile(inputFilePath string) ([]struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
}, error) {
	return parseInputFile(inputFilePath)
}

func loadCache(outputFilePath string) (map[string]interface{}, error) {
	cache := make(map[string]interface{})
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		log.Printf("Output file %s does not exist. Creating it.", outputFilePath)
		emptyFile, createErr := os.Create(outputFilePath)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create output file: %w", createErr)
		}
		if closeErr := emptyFile.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to close output file: %w", closeErr)
		}
	} else if err != nil {
		log.Printf("Error checking file: %s, error: %v", outputFilePath, err)
		return nil, fmt.Errorf("error checking output file: %w", err)
	}

	cacheData, cacheReadErr := os.ReadFile(outputFilePath)
	if cacheReadErr != nil {
		return nil, fmt.Errorf("reading output file: %w", cacheReadErr)
	}

	if len(cacheData) == 0 {
		return cache, nil
	}

	if err := json.Unmarshal(cacheData, &cache); err != nil {
		var cacheArray []types.LinkPreviewOutput
		if err = json.Unmarshal(cacheData, &cacheArray); err == nil {
			for _, item := range cacheArray {
				cache[item.URL] = item.Preview
			}
		} else if string(cacheData) == "[]" {
			cache = make(map[string]interface{})
		} else {
			return nil, fmt.Errorf("parsing output JSON: %w", err)
		}
	}

	return cache, nil
}

func LoadCache(outputFilePath string) (map[string]interface{}, error) {
	return loadCache(outputFilePath)
}

func generatePreviews(urlObjects []struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
}, cache map[string]interface{}, previewer LinkPreviewer, outputFilePath string) ([]types.LinkPreviewOutput, error) {
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
			parsedPreview, err := previewer.Parse(urlObj.URL)
			if err != nil {
				log.Printf("Failed to generate preview for %s: %v", urlObj.URL, err)
				continue
			}

			if parsedPreview == nil {
				log.Printf("Skipping nil preview for %s", urlObj.URL)
				continue
			}

			if parsedPreview.Title == "" &&
				parsedPreview.Description == "" &&
				(parsedPreview.OGMeta == nil && parsedPreview.TwitterMeta == nil) {
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

		// Write the current state of the output to the file after processing each URL
		if writeErr := saveOutput(outputFilePath, output); writeErr != nil {
			log.Printf("Failed to write output file after processing URL %s: %v", urlObj.URL, writeErr)
		}
	}

	if len(output) == 0 {
		log.Println("No valid previews generated.")
		return nil, errors.New("no valid previews generated")
	}

	return output, nil
}

func saveOutput(outputFilePath string, output []types.LinkPreviewOutput) error {
	outputData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling output JSON: %w", err)
	}
	if err = os.WriteFile(outputFilePath, outputData, 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}
	return nil
}

func SaveOutput(outputFilePath string, output []types.LinkPreviewOutput) error {
	return saveOutput(outputFilePath, output)
}
