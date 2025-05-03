# Link Preview Generator

This Go program processes Telegram messages with URLs, validates them, and generates link previews. It outputs cleaned, unique URLs and their previews to JSON files.

## Features

- Extracts and validates URLs from Telegram messages.
- Removes session-related query strings.
- Ensures unique, valid URLs.
- Generates link previews.
- Configurable via command-line arguments or environment variables.

## Usage

### Command-Line Arguments

#### Import/Export

- `-import-urls`: Import URLs from a JSON file (default: `imports/export.json`) and output cleaned URLs (default: `dist/urls.json`).
- `-import-input`: Input JSON file path (default: `import/export.json`).
- `-import-output`: Output JSON file path (default: `dist/urls.json`).

#### Link Previews

- `-generate-preview`: Generate link previews.
- `-preview-input`: Input JSON file for URLs (default: `dist/urls.json`).
- `-preview-output`: Output JSON file for previews (default: `dist/previews.json`).

### Examples

#### Import/Export URLs

```bash
export IMPORT_IGNORE=".*example.com.*"
go run . -import-urls -import-input=import/export.json -import-output=dist/urls.json
```

#### Generate Previews

```bash
export DEBUG=true
go run . -generate-preview -preview-input=dist/urls.json -preview-output=dist/previews.json
```

### Environment Variables

- `IMPORT_IGNORE`: Regex to ignore URLs.
  - **Example**: `export IMPORT_IGNORE=".*example.com.*"`
  - **Default**: No URLs ignored.
- `DEBUG`: Enable debug logging.
  - **Example**: `export DEBUG=true`
  - **Default**: `false`.

## Development

### Prerequisites

- Go 1.24.2 or later.

### Tests

Run all tests:

```bash
go test -v ./...
```

### Project Structure

```plain
.
├── bin
├── dist
├── imports
└── internal
    ├── imports
    ├── previews
    ├── types
    ├── utils
    └── validation
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
