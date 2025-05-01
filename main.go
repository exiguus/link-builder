package main

import (
	"flag"
	"log"
	"urls-processor/internal/imports"
	"urls-processor/internal/previews"
)

func main() {
	log.Println("Starting the URL Processor program")

	importInputFilePath := flag.String("import-input", "imports/export.json", "Path to the input JSON file for import/export")
	importOutputFilePath := flag.String("import-output", "dist/urls.json", "Path to the output JSON file for import/export")
	processImports := flag.Bool("import-urls", false, "Import urls from import/export JSON file")

	previewInputFilePath := flag.String("preview-input", "dist/urls.json", "Path to the input JSON file containing URLs for previews")
	previewOutputFilePath := flag.String("preview-output", "dist/previews.json", "Path to the output JSON file for link previews")
	generatePreviews := flag.Bool("generate-preview", false, "Generate link previews from URLs")

	flag.Parse()

	if *generatePreviews {
		previews.GenerateLinkPreviews(*previewInputFilePath, *previewOutputFilePath)
		log.Println("URL Processor program completed successfully")
		return
	}

	if *processImports {
		imports.ProcessImport(*importInputFilePath, *importOutputFilePath)
		log.Println("URL Processor program completed successfully")
		return
	}

	log.Println("No valid flags provided. Use -generate-import or -generate-previews to run the program.")
}
