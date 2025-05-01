package main

import (
	"flag"
	"link-builder/internal/imports"
	"link-builder/internal/previews"
	"log"
	"os"
)

type Config struct {
	ImportInputFilePath   string
	ImportOutputFilePath  string
	ProcessImports        bool
	PreviewInputFilePath  string
	PreviewOutputFilePath string
	GeneratePreviews      bool
	Debug                 bool
}

func loadConfig() Config {
	config := Config{
		ImportInputFilePath:   "imports/export.json",
		ImportOutputFilePath:  "dist/urls.json",
		ProcessImports:        false,
		PreviewInputFilePath:  "dist/urls.json",
		PreviewOutputFilePath: "dist/previews.json",
		GeneratePreviews:      false,
		Debug:                 false,
	}

	flag.StringVar(
		&config.ImportInputFilePath,
		"import-input",
		"imports/export.json",
		"Path to the input JSON file for import/export",
	)
	flag.StringVar(
		&config.ImportOutputFilePath,
		"import-output",
		"dist/urls.json",
		"Path to the output JSON file for import/export",
	)
	flag.BoolVar(&config.ProcessImports, "import-urls", false, "Import URLs from import/export JSON file")

	flag.StringVar(
		&config.PreviewInputFilePath,
		"preview-input",
		"dist/urls.json",
		"Path to the input JSON file containing URLs for previews",
	)
	flag.StringVar(
		&config.PreviewOutputFilePath,
		"preview-output",
		"dist/previews.json",
		"Path to the output JSON file for link previews",
	)
	flag.BoolVar(&config.GeneratePreviews, "generate-preview", false, "Generate link previews from URLs")

	flag.Parse()

	if os.Getenv("DEBUG") == "true" {
		config.Debug = true
	}

	return config
}

func main() {
	log.Println("Starting the URL Processor program")

	config := loadConfig()

	if config.Debug {
		log.Println("Debug mode enabled")
	}

	if config.GeneratePreviews {
		if err := previews.GenerateLinkPreviews(
			config.PreviewInputFilePath,
			config.PreviewOutputFilePath,
			previews.DefaultLinkPreviewer{},
		); err != nil {
			log.Printf("Error generating link previews: %v", err)
			os.Exit(1)
		}
		log.Println("URL Processor program completed successfully")
		return
	}

	if config.ProcessImports {
		if err := imports.ProcessImport(
			config.ImportInputFilePath,
			config.ImportOutputFilePath,
		); err != nil {
			log.Printf("Error processing imports: %v", err)
			os.Exit(1)
		}
		log.Println("URL Processor program completed successfully")
		return
	}

	log.Println("No valid flags provided. Use -import-urls or -generate-preview to run the program.")
}
