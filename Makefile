# Smart Makefile for Windows and Mac/Linux

APP_NAME=SignatureInstaller
BUILD_DIR=build
TEMPLATE_DIR=templates
SRC_DIR=./cmd/signature-installer
OS := $(shell uname -s)

# Set output binary based on OS
ifeq ($(OS),Windows_NT)
	BINARY=$(BUILD_DIR)/$(APP_NAME).exe
	COPY=cmd /C "if not exist $(BUILD_DIR) mkdir $(BUILD_DIR) & xcopy /E /I /Y $(TEMPLATE_DIR) $(BUILD_DIR)\$(TEMPLATE_DIR)"
	RM=cmd /C "if exist $(BUILD_DIR) rmdir /S /Q $(BUILD_DIR)"
else
	BINARY=$(BUILD_DIR)/$(APP_NAME)
	COPY=mkdir -p $(BUILD_DIR)/$(TEMPLATE_DIR) && cp -r $(TEMPLATE_DIR)/* $(BUILD_DIR)/$(TEMPLATE_DIR)
	RM=rm -rf $(BUILD_DIR)
endif

.PHONY: all build clean run copy-templates

all: build copy-templates

build:
	@echo "==> Building $(APP_NAME) for $(OS)..."
	@go build -o $(BINARY) $(SRC_DIR)

copy-templates:
	@echo "==> Copying templates directory..."
	@$(COPY)

run: build copy-templates
	@echo "==> Running application..."
	@$(BINARY)

clean:
	@echo "==> Cleaning up..."
	@$(RM)
