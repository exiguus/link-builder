package main

import (
	"flag"
	"log"
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

	processImport(*importInputFilePath, *importOutputFilePath, *validateHead)
	log.Println("URL Processor program completed successfully")
}
