.PHONY: help build test clean deps run-examples simulate-receiver simulate-sender

# Default target
help:
	@echo "CRE Networking Test - Makefile Commands"
	@echo "========================================"
	@echo "  make deps          - Download Go dependencies"
	@echo "  make build         - Build the project"
	@echo "  make test          - Run tests"
	@echo "  make run-examples  - Run example usage code"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make simulate-receiver - Simulate receiver workflow locally"
	@echo "  make simulate-sender   - Simulate sender workflow locally"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build the project
build:
	@echo "Building project..."
	go build -o bin/cre-networking-test ./examples

# Run tests
test:
	@echo "Running tests..."
	go test -v ./tests/...

# Run examples
run-examples:
	@echo "Running examples..."
	go run ./examples/example_usage.go

# Simulate receiver workflow
simulate-receiver:
	@echo "Simulating receiver workflow..."
	@echo "Note: This requires CRE CLI and proper configuration"
	@echo "Example command:"
	@echo "  cre workflow simulate --workflow receiver-workflow --input '{\"message\":\"test\"}'"

# Simulate sender workflow
simulate-sender:
	@echo "Simulating sender workflow..."
	@echo "Note: This requires CRE CLI and proper configuration"
	@echo "Example command:"
	@echo "  cre workflow simulate --workflow sender-workflow --input '{\"message\":\"test\",\"target_url\":\"...\"}'"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf *.log
	go clean

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/"; \
	fi

