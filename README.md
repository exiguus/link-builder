# URL Processor

This Go program processes a JSON file containing messages with URLs, validates the URLs, and outputs a sorted list of unique, valid URLs to a JSON file. It supports optional `HEAD` request validation and can be configured via command-line arguments or environment variables.

## Features

- Reads a JSON file with messages and extracts URLs.
- Validates URLs for proper structure and optional reachability using `HEAD` requests.
- Outputs a sorted list of unique, valid URLs to a JSON file.
- Generates link previews for URLs and saves them to a JSON file.
- Supports concurrent URL validation for improved performance.
- Configurable via command-line arguments or environment variables.

## Usage

### Command-Line Arguments

#### Import/Export Functionality

- `-import-input`: Path to the input JSON file for import/export (default: `import/export.json`).
- `-import-output`: Path to the output JSON file for import/export (default: `dist/urls.json`).
- `-validate-head`: Enable `HEAD` requests to validate URLs (default: `false`).

#### Link Preview Functionality

- `-generate-previews`: Generate link previews from URLs.
- `-preview-input`: Path to the input JSON file containing URLs for previews (default: `dist/urls.json`).
- `-preview-output`: Path to the output JSON file for link previews (default: `dist/previews.json`).

### Examples

#### Process URLs (Import/Export)

```bash
go run . -import-input=import/export.json -import-output=dist/urls.json
```

#### Generate Link Previews

```bash
go run . -generate-previews -preview-input=dist/urls.json -preview-output=dist/previews.json
```

### Environment Variables

The following environment variables can be used to configure the application:

- `INPUT_FILE`: Path to the input JSON file.
- `OUTPUT_FILE`: Path to the output JSON file.
- `VALIDATE_HEAD`: Set to `true` to enable `HEAD` requests for URL validation.
- `IMPORT_IGNORE`: A regex pattern to ignore specific URLs during processing. If not set, no URLs are ignored.
- `DEBUG`: Set to `true` to enable debug logging for additional output during processing.

Example usage:

```bash
export INPUT_FILE=import/export.json
export OUTPUT_FILE=dist/urls.json
export VALIDATE_HEAD=true
export IMPORT_IGNORE=".*example.com.*"
export DEBUG=true
```

## Development

### Prerequisites

- Go 1.16 or later.

### Running Tests

Run the following command to execute all tests:

```bash
go test -v .
```

### Project Structure

- `main.go`: Main program logic.
- `main_test.go`: Unit tests for the program.
- `utils.go`: Utility functions for file operations, error handling, and URL validation.
- `validation.go`: Functions for validating and processing URLs, including concurrency support.
- `previews.go`: Logic for generating link previews from URLs.
- `types.go`: Shared data structures and types used across the project.
- `import/export.json`: Example input file.
- `dist/urls.json`: Output file (generated).
- `dist/previews.json`: Output file for link previews (generated).

## License

This project is licensed under the MIT License.
