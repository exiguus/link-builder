package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tiendc/go-linkpreview"
	"golang.org/x/time/rate"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			IdleConnTimeout:       10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 15 * time.Second,
	}
)

// readJSONFile reads a JSON file and unmarshals it into the provided interface.
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

// writeJSONFile writes a JSON object to a file.
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

// isValidURL validates a URL and optionally performs a HEAD request.
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

// validateURLsConcurrently validates a list of URLs concurrently using goroutines.
func validateURLsConcurrently(urls []string, validateHead bool, ignoreRegex *regexp.Regexp) (map[string]bool, int) {
	validURLs := make(map[string]bool)
	var mutex sync.Mutex
	var ignoredCount int32

	rateLimiter := rate.NewLimiter(rate.Limit(10), 1)
	processURL := func(rawURL string) {
		if ignoreRegex != nil && ignoreRegex.MatchString(rawURL) {
			atomic.AddInt32(&ignoredCount, 1)
			return
		}
		if validateHead {
			if err := rateLimiter.Wait(context.Background()); err != nil {
				log.Printf("Rate limiter error: %v", err)
				return
			}
		}
		if isValidURL(rawURL, validateHead) {
			mutex.Lock()
			validURLs[rawURL] = true
			mutex.Unlock()
		}
	}

	var wg sync.WaitGroup
	for _, rawURL := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			processURL(url)
		}(rawURL)
	}
	wg.Wait()

	return validURLs, int(ignoredCount)
}

// removeSessionQueryStrings removes session-related query strings from valid URLs.
func removeSessionQueryStrings(validURLs map[string]bool) map[string]bool {
	updatedURLs := make(map[string]bool)
	for urlStr := range validURLs {
		// Remove session identifiers attached with a semicolon
		if semicolonIndex := strings.Index(urlStr, ";jsessionid="); semicolonIndex != -1 {
			urlStr = urlStr[:semicolonIndex]
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("Failed to parse URL %s: %v", urlStr, err)
			continue
		}
		query := parsedURL.Query()
		for key := range query {
			if strings.Contains(strings.ToLower(key), "session") {
				query.Del(key)
			}
		}
		parsedURL.RawQuery = query.Encode()
		updatedURLs[parsedURL.String()] = true
	}
	return updatedURLs
}

// warnIfURLsContainSession logs a warning for URLs containing the word "session".
func warnIfURLsContainSession(validURLs map[string]bool) {
	for url := range validURLs {
		if strings.Contains(strings.ToLower(url), "session") {
			log.Printf("Warning: URL contains 'session': %s", url)
		}
	}
}

// ensureUniqueURLs ensures that valid URLs are unique by keeping only the first occurrence.
func ensureUniqueURLs(validURLs map[string]bool, allURLs []struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
}) map[string]bool {
	uniqueURLs := make(map[string]bool)
	seen := make(map[string]bool)

	for _, urlObj := range allURLs {
		if validURLs[urlObj.URL] && !seen[urlObj.URL] {
			uniqueURLs[urlObj.URL] = true
			seen[urlObj.URL] = true
		}
	}

	return uniqueURLs
}

// handleError logs and exits the program if an error occurs.
func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

// getConfigValue retrieves a configuration value from an environment variable or falls back to a default.
func getConfigValue(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// compileIgnoreRegex compiles the ignore regex pattern from the environment variable.
func compileIgnoreRegex() (*regexp.Regexp, error) {
	ignorePattern := getConfigValue("IMPORT_IGNORE", "")
	if ignorePattern == "" {
		return nil, nil
	}
	return regexp.Compile(ignorePattern)
}

// LinkPreviewOutput represents the structure of the link preview output.
type LinkPreviewOutput struct {
	ID      int         `json:"id"`
	Date    string      `json:"date"`
	URL     string      `json:"url"`
	Preview interface{} `json:"preview"`
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// generateLinkPreviews generates link previews and writes them to the output file.
func generateLinkPreviews(inputFilePath, outputFilePath string) {
	data, err := os.ReadFile(inputFilePath)
	handleError(err, "Failed to read input file")

	var urlObjects []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	handleError(json.Unmarshal(data, &urlObjects), "Failed to parse input JSON")

	cache := make(map[string]interface{})
	if _, err := os.Stat(outputFilePath); err == nil {
		cacheData, err := os.ReadFile(outputFilePath)
		handleError(err, "Failed to read output file")
		if len(cacheData) > 0 {
			if err := json.Unmarshal(cacheData, &cache); err != nil {
				var cacheArray []LinkPreviewOutput
				if err := json.Unmarshal(cacheData, &cacheArray); err == nil {
					for _, item := range cacheArray {
						cache[item.URL] = item.Preview
					}
				} else if string(cacheData) == "[]" {
					cache = make(map[string]interface{})
				} else {
					handleError(err, "Failed to parse output JSON")
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

	output := []LinkPreviewOutput{}
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

		output = append(output, LinkPreviewOutput{
			ID:      urlObj.ID,
			Date:    urlObj.Date,
			URL:     urlObj.URL,
			Preview: preview,
		})

		outputData, err := json.MarshalIndent(output, "", "  ")
		handleError(err, "Failed to marshal intermediate output JSON")
		handleError(os.WriteFile(outputFilePath, outputData, 0644), "Failed to write intermediate output file")
	}

	log.Printf("Link previews successfully generated and saved to %s", outputFilePath)
}

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
