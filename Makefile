.PHONY: build test lint install clean

BINARY_NAME=lane
BUILD_DIR=./bin
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/dlinardi/lane/cmd.version=$(VERSION) -X github.com/dlinardi/lane/cmd.commit=$(COMMIT) -X github.com/dlinardi/lane/cmd.date=$(DATE)"

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

lint:
	golangci-lint run

install:
	go install $(LDFLAGS) .

clean:
	rm -rf $(BUILD_DIR)
