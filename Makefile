# Makefile for BatchInvoice PDF

.PHONY: all build clean run test deps help

# Variables
APP_NAME=batchinvoice-pdf
VERSION=1.0.0
BUILD_DIR=build
GO=go
GOFLAGS=-v

# Default target
all: deps build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Build for current platform
build:
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Build for all platforms (only current OS - Fyne uses CGO/OpenGL, no cross-compile)
build-all:
	@chmod +x build.sh
	@./build.sh

# Run the application
run:
	@echo "🚀 Running $(APP_NAME)..."
	$(GO) run main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GO) test ./... -v

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Format code
fmt:
	@echo "📝 Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  make deps       - Install dependencies"
	@echo "  make build      - Build for current platform"
	@echo "  make build-all  - Build for current platform (Fyne cannot cross-compile)"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"
	@echo "  make help       - Show this help message"
