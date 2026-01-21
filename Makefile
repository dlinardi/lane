.PHONY: build test lint install clean help

BINARY_NAME=lane
BUILD_DIR=./bin
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/dlinardi/lane/cmd.version=$(VERSION) -X github.com/dlinardi/lane/cmd.commit=$(COMMIT) -X github.com/dlinardi/lane/cmd.date=$(DATE)"

build: ## Build to ./bin/lane
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test: ## Run tests with coverage
	go test -v -cover ./...

lint: ## Run golangci-lint
	golangci-lint run

install: ## Install to $GOPATH/bin
	go install $(LDFLAGS) .

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
