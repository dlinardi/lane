.PHONY: build test lint install clean

BINARY_NAME=lane
BUILD_DIR=./bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

lint:
	golangci-lint run

install:
	go install .

clean:
	rm -rf $(BUILD_DIR)
