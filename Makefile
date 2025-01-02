# Build targets
.PHONY: all build test lint coverage clean

# Default target
all: build test lint coverage

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./... -v

# Run linter
lint:
	golangci-lint run > lint.txt

# Generate test coverage
coverage:
	go test ./... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -func=coverage.txt

# Clean build artifacts
clean:
	rm -f coverage.txt lint.txt
	go clean
