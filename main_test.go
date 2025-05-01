package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tiendc/go-linkpreview"
)

// TestMainFunctionWithMocks tests the main function with mock input and output files.
func TestMainFunctionWithMocks(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_main_function")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFilePath := filepath.Join(tempDir, "mock_input.json")
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	inputData := `{
		"messages": [
			{
				"text_entities": [
					{"type": "link", "text": "http://example.com"},
					{"type": "link", "text": "https://example.org"},
					{"type": "text", "text": "Not a URL"}
				]
			},
			{
				"text_entities": [
					{"type": "link", "text": "http://another-example.com"}
				]
			}
		]
	}`

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	cmd := exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "-import-input="+inputFilePath, "-import-output="+outputFilePath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run main program: %v", err)
	}

	outputData, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var output []struct {
		ID      int         `json:"id"`
		Date    string      `json:"date"`
		URL     string      `json:"url"`
		Preview interface{} `json:"preview"`
	}
	if err := json.Unmarshal(outputData, &output); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	expectedURLs := map[string]bool{
		"http://example.com":         true,
		"https://example.org":        true,
		"http://another-example.com": true,
	}

	for _, entry := range output {
		if entry.URL == "" {
			t.Errorf("Empty URL found in output: %+v", entry)
			continue
		}
		if !expectedURLs[entry.URL] {
			t.Errorf("Unexpected URL in output: %s", entry.URL)
		}
		delete(expectedURLs, entry.URL)
	}

	if len(expectedURLs) > 0 {
		t.Errorf("Missing expected URLs in output: %+v", expectedURLs)
	}
}

// TestProcessImport tests the processImport function.
func TestProcessImport(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_process_import")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFilePath := filepath.Join(tempDir, "mock_input.json")
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	inputData := `{
		"messages": [
			{
				"text_entities": [
					{"type": "link", "text": "http://example.com"},
					{"type": "link", "text": "https://example.org"},
					{"type": "text", "text": "Not a URL"}
				]
			},
			{
				"text_entities": [
					{"type": "link", "text": "http://another-example.com"}
				]
			}
		]
	}`

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	processImport(inputFilePath, outputFilePath, false)

	outputData, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var output []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}
	if err := json.Unmarshal(outputData, &output); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	expectedURLs := map[string]bool{
		"http://example.com":         true,
		"https://example.org":        true,
		"http://another-example.com": true,
	}

	for _, entry := range output {
		if entry.URL == "" {
			t.Errorf("Empty URL found in output: %+v", entry)
			continue
		}
		if !expectedURLs[entry.URL] {
			t.Errorf("Unexpected URL in output: %s", entry.URL)
		}
		delete(expectedURLs, entry.URL)
	}

	if len(expectedURLs) > 0 {
		t.Errorf("Missing expected URLs in output: %+v", expectedURLs)
	}
}

// TestHelpCommand tests the help command for expected output.
func TestHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	helpOutput := string(output)
	if !strings.Contains(helpOutput, "-import-input") || !strings.Contains(helpOutput, "-import-output") || !strings.Contains(helpOutput, "-validate-head") {
		t.Errorf("Help output does not contain expected flags: %s", helpOutput)
	}
}

// TestInvalidCommand tests the behavior for invalid command-line flags.
func TestInvalidCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "--invalid-flag")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected error for invalid flag, but got none")
	}

	errorOutput := string(output)
	if !strings.Contains(errorOutput, "flag provided but not defined") {
		t.Errorf("Unexpected error output for invalid flag: %s", errorOutput)
	}
}

// TestReadJSONFile tests the readJSONFile function.
func TestReadJSONFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test_input_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	mockData := `{"messages": [{"text_entities": [{"type": "link", "text": "http://example.com"}]}]}`
	if _, err := tempFile.Write([]byte(mockData)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	defer tempFile.Close()

	var input struct {
		Messages []struct {
			TextEntities []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"text_entities"`
		} `json:"messages"`
	}
	if err := readJSONFile(tempFile.Name(), &input); err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	if len(input.Messages) != 1 || len(input.Messages[0].TextEntities) != 1 {
		t.Errorf("Unexpected data in JSON file: %+v", input)
	}
}

// TestWriteJSONFile tests the writeJSONFile function.
func TestWriteJSONFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.Close()

	mockData := []string{"http://example.com", "https://example.org"}

	if err := writeJSONFile(tempFile.Name(), mockData); err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	data, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read back JSON file: %v", err)
	}

	var output []string
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("Failed to parse JSON file: %v", err)
	}

	if len(output) != 2 || output[0] != "http://example.com" || output[1] != "https://example.org" {
		t.Errorf("Unexpected data in JSON file: %+v", output)
	}
}

// TestWriteJSONFileError tests the writeJSONFile function for error handling.
func TestWriteJSONFileError(t *testing.T) {
	// Attempt to write to an invalid directory
	invalidFilePath := "/invalid_directory/output.json"
	mockData := []string{"http://example.com", "https://example.org"}

	err := writeJSONFile(invalidFilePath, mockData)
	if err == nil {
		t.Fatalf("Expected error when writing to an invalid directory, but got none")
	}

	// Verify the error message contains the file path
	if !strings.Contains(err.Error(), invalidFilePath) {
		t.Errorf("Error message does not contain file path: %v", err)
	}
}

// TestIsValidURL tests the isValidURL function.
func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		url          string
		validateHead bool
		expected     bool
	}{
		{"http://example.com", false, true},
		{"https://example.org", false, true},
		{"ftp://example.com", false, false},
		{"invalid-url", false, false},
	}

	for _, tc := range testCases {
		result := isValidURL(tc.url, tc.validateHead)
		if result != tc.expected {
			t.Errorf("isValidURL(%q, %v) = %v; want %v", tc.url, tc.validateHead, result, tc.expected)
		}
	}
}

// TestLinkPreviewParse tests the linkpreview.Parse function.
func TestLinkPreviewParse(t *testing.T) {
	url := "http://example.com"
	preview, err := linkpreview.Parse(url)
	if err != nil {
		t.Fatalf("Failed to parse link preview for %s: %v", url, err)
	}

	if preview.Title == "" && preview.Description == "" && preview.OGMeta == nil && preview.TwitterMeta == nil {
		t.Errorf("Invalid preview data for %s", url)
	}
}

// TestInputOutputFlags tests the input and output flags.
func TestInputOutputFlags(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_input_output_flags")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFilePath := filepath.Join(tempDir, "mock_input.json")
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	inputData := `{
		"messages": [
			{
				"text_entities": [
					{"type": "link", "text": "http://example.com"}
				]
			}
		]
	}`

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	cmd := exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "-import-input="+inputFilePath, "-import-output="+outputFilePath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run main program with valid flags: %v", err)
	}

	cmd = exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "-import-input=invalid_input.json", "-import-output="+outputFilePath)
	if err := cmd.Run(); err == nil {
		t.Fatalf("Expected error for invalid input file, but got none")
	}
}

// TestGeneratePreviewsFlag tests the -generate-previews flag.
func TestGeneratePreviewsFlag(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_generate_previews")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFilePath := filepath.Join(tempDir, "mock_input.json")
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	inputData := `[
		{"id": 1, "date": "2025-05-01", "url": "http://example.com"},
		{"id": 2, "date": "2025-05-01", "url": "https://example.org"}
	]`

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		t.Fatalf("Input file does not exist: %s", inputFilePath)
	}

	cmd := exec.Command("go", "run", "main.go", "previews.go", "import.go", "utils.go", "types.go", "validation.go", "-generate-previews", "-preview-input="+inputFilePath, "-preview-output="+outputFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run main program with -generate-previews flag: %v\nOutput: %s", err, string(output))
	}

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFilePath)
	}

	outputData, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var previews []struct {
		ID      int         `json:"id"`
		Date    string      `json:"date"`
		URL     string      `json:"url"`
		Preview interface{} `json:"preview"`
	}
	if err := json.Unmarshal(outputData, &previews); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	if len(previews) != 2 {
		t.Errorf("Expected 2 previews, got %d", len(previews))
	}
}

// TestRemoveSessionQueryStrings tests the removal of session query strings.
func TestRemoveSessionQueryStrings(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com;jsessionid=12345":           true,
		"http://example.org/path?sessionid=abc":         true,
		"http://example.net/path?key=value&session=xyz": true,
		"http://example.edu/path?key=value":             true,
	}

	expectedURLs := map[string]bool{
		"http://example.com":                true,
		"http://example.org/path":           true,
		"http://example.net/path?key=value": true,
		"http://example.edu/path?key=value": true,
	}

	result := removeSessionQueryStrings(validURLs)

	if len(result) != len(expectedURLs) {
		t.Fatalf("Expected %d URLs, got %d", len(expectedURLs), len(result))
	}

	for url := range expectedURLs {
		if !result[url] {
			t.Errorf("Expected URL %s not found in result", url)
		}
	}
}

// TestEnsureUniqueURLs tests ensuring unique valid URLs.
func TestEnsureUniqueURLs(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
		"http://example.net": true,
	}

	allURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{
		{ID: 1, Date: "2025-05-01", URL: "http://example.com"},
		{ID: 2, Date: "2025-05-01", URL: "http://example.org"},
		{ID: 3, Date: "2025-05-01", URL: "http://example.com"},
		{ID: 4, Date: "2025-05-01", URL: "http://example.net"},
	}

	expectedUniqueURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
		"http://example.net": true,
	}

	result := ensureUniqueURLs(validURLs, allURLs)

	if len(result) != len(expectedUniqueURLs) {
		t.Fatalf("Expected %d unique URLs, got %d", len(expectedUniqueURLs), len(result))
	}

	for url := range expectedUniqueURLs {
		if !result[url] {
			t.Errorf("Expected URL %s not found in result", url)
		}
	}
}

// TestWarnIfURLsContainSession tests warnings for URLs containing session identifiers.
func TestWarnIfURLsContainSession(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com;jsessionid=12345":           true,
		"http://example.org/path?sessionid=abc":         true,
		"http://example.net/path?key=value&session=xyz": true,
		"http://example.edu/path?key=value":             true,
	}

	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	warnIfURLsContainSession(validURLs)

	logContent := logOutput.String()

	expectedWarnings := []string{
		"http://example.com;jsessionid=12345",
		"http://example.org/path?sessionid=abc",
		"http://example.net/path?key=value&session=xyz",
	}

	for _, url := range expectedWarnings {
		if !strings.Contains(logContent, url) {
			t.Errorf("Expected warning for URL %s not found in log output", url)
		}
	}

	unexpectedWarnings := []string{
		"http://example.edu/path?key=value",
	}

	for _, url := range unexpectedWarnings {
		if strings.Contains(logContent, url) {
			t.Errorf("Unexpected warning for URL %s found in log output", url)
		}
	}
}

// TestEnsureUniqueValidURLs tests ensuring unique valid URLs.
func TestEnsureUniqueValidURLs(t *testing.T) {
	validURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
		"http://example.net": true,
	}

	allURLs := []struct {
		ID   int    `json:"id"`
		Date string `json:"date"`
		URL  string `json:"url"`
	}{
		{ID: 1, Date: "2025-05-01", URL: "http://example.com"},
		{ID: 2, Date: "2025-05-01", URL: "http://example.org"},
		{ID: 3, Date: "2025-05-01", URL: "http://example.com"},
		{ID: 4, Date: "2025-05-01", URL: "http://example.net"},
	}

	expectedUniqueURLs := map[string]bool{
		"http://example.com": true,
		"http://example.org": true,
		"http://example.net": true,
	}

	result := ensureUniqueURLs(validURLs, allURLs)

	if len(result) != len(expectedUniqueURLs) {
		t.Fatalf("Expected %d unique URLs, got %d", len(expectedUniqueURLs), len(result))
	}

	for url := range expectedUniqueURLs {
		if !result[url] {
			t.Errorf("Expected URL %s not found in result", url)
		}
	}
}

// TestValidateURLsConcurrently tests the validateURLsConcurrently function.
func TestValidateURLsConcurrently(t *testing.T) {
	// Test with an empty URL list
	validURLs, ignoredCount := validateURLsConcurrently([]string{}, false, nil)
	if len(validURLs) != 0 {
		t.Errorf("Expected 0 valid URLs, got %d", len(validURLs))
	}
	if ignoredCount != 0 {
		t.Errorf("Expected 0 ignored URLs, got %d", ignoredCount)
	}

	// Test with invalid URLs
	urls := []string{"invalid-url", "ftp://example.com"}
	validURLs, ignoredCount = validateURLsConcurrently(urls, false, nil)
	if len(validURLs) != 0 {
		t.Errorf("Expected 0 valid URLs, got %d", len(validURLs))
	}
	if ignoredCount != 0 {
		t.Errorf("Expected 0 ignored URLs, got %d", ignoredCount)
	}

	// Test with valid URLs
	urls = []string{"http://example.com", "https://example.org"}
	validURLs, ignoredCount = validateURLsConcurrently(urls, false, nil)
	if len(validURLs) != 2 {
		t.Errorf("Expected 2 valid URLs, got %d", len(validURLs))
	}
	if ignoredCount != 0 {
		t.Errorf("Expected 0 ignored URLs, got %d", ignoredCount)
	}
}

// TestCompileIgnoreRegex tests the compileIgnoreRegex function.
func TestCompileIgnoreRegex(t *testing.T) {
	os.Setenv("IMPORT_IGNORE", "example\\.com")
	defer os.Unsetenv("IMPORT_IGNORE")

	regex, err := compileIgnoreRegex()
	if err != nil {
		t.Fatalf("Failed to compile valid regex: %v", err)
	}
	if regex == nil || !regex.MatchString("example.com") {
		t.Errorf("Expected regex to match 'example.com', but it did not")
	}

	os.Setenv("IMPORT_IGNORE", "[invalid-regex")
	_, err = compileIgnoreRegex()
	if err == nil {
		t.Fatalf("Expected error for invalid regex, but got none")
	}
}

// TestGenerateLinkPreviews tests the generateLinkPreviews function.
func TestGenerateLinkPreviews(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test_generate_link_previews")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputFilePath := filepath.Join(tempDir, "mock_input.json")
	outputFilePath := filepath.Join(tempDir, "mock_output.json")

	inputData := `[
		{"id": 1, "date": "2025-05-01", "url": "http://example.com"},
		{"id": 2, "date": "2025-05-01", "url": "invalid-url"}
	]`

	if err := ioutil.WriteFile(inputFilePath, []byte(inputData), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	generateLinkPreviews(inputFilePath, outputFilePath)

	outputData, err := ioutil.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var previews []struct {
		ID      int         `json:"id"`
		Date    string      `json:"date"`
		URL     string      `json:"url"`
		Preview interface{} `json:"preview"`
	}
	if err := json.Unmarshal(outputData, &previews); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	if len(previews) != 1 {
		t.Errorf("Expected 1 valid preview, got %d", len(previews))
	}
}
