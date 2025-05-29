# Makefile for managing the project

GO_VERSION := $(shell grep '^go ' go.mod | awk '{print $$2}')
MOD_NAME := $(shell grep '^module ' go.mod | awk '{print $$2}')
BUILD_BIN := bin/$(MOD_NAME)

.PHONY: all test lint lint-fix fmt qlty-fmt qlty-check qlty-smells qlty-metrics qlty coverage build run-import build-run-import run-preview build-run-preview build-run run clean setup hooks

all: fmt test lint coverage qlty

setup:
	@echo "Setting up Go $(GO_VERSION)"
	@command -v go >/dev/null 2>&1 || { \
		echo "Go is not installed. Please install Go $(GO_VERSION)"; \
		exit 1; \
	}
	@go mod tidy

hooks:
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@ln -sf ../../scripts/commit-msg.sh .git/hooks/commit-msg
	@ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit
	@echo "Git hooks set up successfully."

test:
	@echo "Running tests..."
	@go clean -testcache
	@go test ./... -v

lint:
	@echo "Running golangci-lint..."
	@command -v bin/golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.5; \
	}
	@bin/golangci-lint run ./...

lint-fix:
	@echo "Running golangci-lint with fix..."
	@command -v bin/golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.5; \
	}
	@bin/golangci-lint run --fix ./...

fmt: 
	@echo "Formatting code with gofmt..."
	@command -v gofmt >/dev/null 2>&1 || { \
		echo "Installing gofmt..."; \
		go install golang.org/x/tools/cmd/gofmt@latest; \
	}
	@find . -name '*.go' -not -path './vendor/*' -exec gofmt -w {} \;
	@echo "Running goimports..."
	@command -v goimports >/dev/null 2>&1 || { \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	}
	@find . -name '*.go' -not -path './vendor/*' -exec goimports -w {} \;

qlty-fmt:
	@echo "Running quality checks and formatting..."
	@command -v qlty >/dev/null 2>&1 || { \
		echo "Installing qlty..."; \
		curl https://qlty.sh | bash; \
		export QLTY_TELEMETRY="off"; \
	}
	@qlty fmt --all

qlty-check:
	@echo "Running qlty lint..."
	@command -v qlty >/dev/null 2>&1 || { \
		echo "Installing qlty..."; \
		curl https://qlty.sh | bash; \
		export QLTY_TELEMETRY="off"; \
	}
	@qlty check --sample=12

qlty-smells:
	@echo "Running qlty smells..."
	@command -v qlty >/dev/null 2>&1 || { \
		echo "Installing qlty..."; \
		curl https://qlty.sh | bash; \
		export QLTY_TELEMETRY="off"; \
	}
	@qlty smells --all

qlty-metrics:
	@echo "Running qlty metrics..."
	@command -v qlty >/dev/null 2>&1 || { \
		echo "Installing qlty..."; \
		curl https://qlty.sh | bash; \
		export QLTY_TELEMETRY="off"; \
	}
	@qlty metrics --max-depth=5 --sort complexity --all

qlty:
	@echo "Running all qlty checks..."
	@command -v qlty >/dev/null 2>&1 || { \
		echo "Installing qlty..."; \
		curl https://qlty.sh | bash; \
		export QLTY_TELEMETRY="off"; \
	}
	@make qlty-fmt
	@make qlty-check
	@make qlty-smells
	@make qlty-metrics

coverage:
	@echo "Running tests with coverage..."
	@go clean -testcache
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' > coverage.txt
	@COVERAGE=$$(cat coverage.txt); \
	if [ $$(echo "$$COVERAGE < 70" | bc -l) -eq 1 ]; then \
		echo "Coverage is below 70%: $$COVERAGE%"; \
		exit 1; \
	else \
		echo "Coverage is sufficient: $$COVERAGE%"; \
	fi
	@rm -f coverage.out coverage.txt

build:
	@echo "Building the project..."
	@go build -o $(BUILD_BIN) .

build-run-import:
	@echo "Running the import process..."
	@$(BUILD_BIN) -import-urls -import-input=imports/export.json -import-output=dist/urls.json

build-run-preview:
	@echo "Running the preview generation..."
	@$(BUILD_BIN) -generate-preview -preview-input=dist/urls.json -preview-output=dist/previews.json

build-run:
	@echo "Running the application..."
	make build-run-import && make build-run-preview

run-import:
	@echo "Running the import process..."
	@go run . -import-urls -import-input=imports/export.json -import-output=dist/urls.json

run-preview:
	@echo "Running the preview generation..."
	@go run . -generate-preview -preview-input=dist/urls.json -preview-output=dist/previews.json

run:
	@echo "Running the application..."
	make run-import && make run-preview

clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.txt
	@rm -rf $(BUILD_BIN)