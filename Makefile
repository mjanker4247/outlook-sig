.PHONY: build-mac build-win clean deps

# Detect OS
ifeq ($(OS),Windows_NT)
    MKDIR=mkdir /p
    RMDIR=rmdir /S /Q
    CP=copy
    RM=del /Q /F
    SEP=\\
    EXE_EXT=.exe
else
    MKDIR=mkdir -p
    RMDIR=rm -rf
    CP=cp -r
    RM=rm -f
    SEP=/
    EXE_EXT=
endif

# Build configuration
BINARY_NAME=signature-installer
VERSION=1.0.0
BUILD_DIR=build
MAC_DIR=$(BUILD_DIR)$(SEP)mac
WIN_DIR=$(BUILD_DIR)$(SEP)win
ASSETS_DIR=assets

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed"

# Build for macOS
build-mac: deps
	@echo "Building for macOS..."
	@$(MKDIR) $(MAC_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(MAC_DIR)$(SEP)$(BINARY_NAME)$(EXE_EXT) .$(SEP)cmd$(SEP)signature-installer
	@$(MKDIR) $(MAC_DIR)$(SEP)$(ASSETS_DIR)
	@$(CP) $(ASSETS_DIR)$(SEP)* $(MAC_DIR)$(SEP)$(ASSETS_DIR)$(SEP)
	@echo "Build complete. Output in $(MAC_DIR)"

# Build for Windows
build-win: deps
	@echo "Building for Windows..."
	@$(MKDIR) $(WIN_DIR)
	go build -o $(WIN_DIR)$(SEP)$(BINARY_NAME).exe .$(SEP)cmd$(SEP)signature-installer
	@$(MKDIR) $(WIN_DIR)$(SEP)$(ASSETS_DIR)
	if exist $(ASSETS_DIR) xcopy /E /I /Y $(ASSETS_DIR)$(SEP)* $(WIN_DIR)$(SEP)$(ASSETS_DIR)$(SEP)
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