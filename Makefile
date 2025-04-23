.PHONY: build-mac build-win clean deps

# Build configuration
BINARY_NAME=signature-installer
VERSION=1.0.0
BUILD_DIR=build
MAC_DIR=$(BUILD_DIR)/mac
WIN_DIR=$(BUILD_DIR)/win

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed"

# Build for macOS
build-mac: deps
	@echo "Building for macOS..."
	@mkdir -p $(MAC_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(MAC_DIR)/$(BINARY_NAME) ./cmd/signature-installer
	@cp -r templates $(MAC_DIR)/
	@echo "Build complete. Output in $(MAC_DIR)"

# Build for Windows
build-win: deps
	@echo "Building for Windows..."
	@mkdir -p $(WIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(WIN_DIR)/$(BINARY_NAME).exe ./cmd/signature-installer
	@cp -r templates $(WIN_DIR)/
	@echo "Build complete. Output in $(WIN_DIR)"

# Build for all platforms
build-all: build-mac build-win
	@echo "Build complete for all platforms"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Default target
all: build-all 