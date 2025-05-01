# URL Processor

This Go program processes a JSON file containing Telegram messages with URLs, validates the URLs, and outputs a sorted list of unique, valid URLs to a JSON file. It also generates link previews for the URLs and saves them to a separate JSON file. The program can be configured via command-line arguments or environment variables.

## Features

- Reads a JSON file with messages and extracts URLs.
- Validates URLs for proper structure.
- Removes session-related query strings from URLs for cleaner output.
- Warns if URLs contain session-related identifiers.
- Ensures uniqueness of URLs in the output.
- Outputs a sorted list of unique, valid URLs to a JSON file.
- Generates link previews for URLs and saves them to a JSON file.
- Supports concurrent URL validation for improved performance.
- Configurable via command-line arguments or environment variables.

## Usage

### Command-Line Arguments

#### Import/Export Functionality

- `-import-urls`: Import URLs from a JSON file (default: `imports/export.json`) and create a JSON file that can be used to generate link previews (default: `dist/urls.json`). The `export.json` file should contain an array of objects with a `message` field containing the URLs. It can be generated using the Telegram export Chat functionality.
- `-import-input`: Path to the input JSON file for import/export (default: `import/export.json`).
- `-import-output`: Path to the output JSON file for import/export (default: `dist/urls.json`).

#### Link Preview Functionality

- `-generate-preview`: Generate link previews from URLs.
- `-preview-input`: Path to the input JSON file containing URLs for previews (default: `dist/urls.json`).
- `-preview-output`: Path to the output JSON file for link previews (default: `dist/previews.json`).

### Examples

#### Process URLs (Import/Export) with Session Query String Removal

```bash
go run . -import-urls -import-input=import/export.json -import-output=dist/urls.json
```

#### Generate Link Previews

```bash
go run . -generate-previews -preview-input=dist/urls.json -preview-output=dist/previews.json
```

### Environment Variables

The following environment variables can be used to configure the application:

- `INPUT_FILE`: Path to the input JSON file.
- `OUTPUT_FILE`: Path to the output JSON file.
- `IMPORT_IGNORE`: A regex pattern to ignore specific URLs during processing. If not set, no URLs are ignored.
- `DEBUG`: Set to `true` to enable debug logging for additional output during processing.

Example usage:

```bash
export INPUT_FILE=import/export.json
export OUTPUT_FILE=dist/urls.json
export IMPORT_IGNORE=".*example.com.*"
export DEBUG=true
```

## Development

### Prerequisites

- Go 1.16 or later.

### Running Tests

Run the following command to execute all tests:

```bash
go test -v ./...
```

### Project Structure

- `main.go`: Main program logic.
- `main_test.go`: Unit tests for the program.
- `internal/`: Internal package for URL processing and validation.
- `internal/preview/`: Package for generating link previews.
- `internal/imports/`: Package for importing/exporting URLs.
- `internal/utils/`: Utility functions for the program.
- `internal/validation/`: Package for URL validation.
- `imports/export.json`: Example input file.
- `dist/urls.json`: Output file (generated).
- `dist/previews.json`: Output file for link previews (generated).

## License

This project is licensed under the MIT License.
