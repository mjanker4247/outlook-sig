.PHONY: build-mac build-win clean deps

# Build configuration
BINARY_NAME=signature-installer
VERSION=1.0.0
BUILD_DIR=build
MAC_DIR=$(BUILD_DIR)/mac
WIN_DIR=$(BUILD_DIR)/win

# Detect OS
ifeq ($(OS),Windows_NT)
    MKDIR=mkdir
    RMDIR=rmdir /s /q
    CP=copy
    RM=del /q
    SEP=\\
else
    MKDIR=mkdir -p
    RMDIR=rm -rf
    CP=cp -r
    RM=rm -f
    SEP=/
endif

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed"

# Build for macOS
build-mac: deps
	@echo "Building for macOS..."
	@$(MKDIR) $(MAC_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(MAC_DIR)$(SEP)$(BINARY_NAME) ./cmd/signature-installer
	@$(CP) templates$(SEP)* $(MAC_DIR)$(SEP)
	@echo "Build complete. Output in $(MAC_DIR)"

# Build for Windows
build-win: deps
	@echo "Building for Windows..."
	@$(MKDIR) $(WIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(WIN_DIR)$(SEP)$(BINARY_NAME).exe ./cmd/signature-installer
	@$(CP) templates$(SEP)* $(WIN_DIR)$(SEP)
	@echo "Build complete. Output in $(WIN_DIR)"

# Build for all platforms
build-all: build-mac build-win
	@echo "Build complete for all platforms"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@$(RMDIR) $(BUILD_DIR)
	@echo "Clean complete"

# Default target
all: build-all 