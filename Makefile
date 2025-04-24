# Smart Makefile for Windows and Mac/Linux

APP_NAME=SignatureInstaller
BUILD_DIR=build
TEMPLATE_DIR=templates
SRC_DIR=./cmd/signature-installer
OS := $(shell uname -s)

# Determine OS
ifeq ($(OS),Windows_NT)
	OS_NAME := windows
	RM := del /Q /F
	MKDIR := mkdir
	CP := xcopy /E /I /Y
	SEP := \\
else
	OS_NAME := $(shell uname -s | tr '[:upper:]' '[:lower:]')
	RM := rm -rf
	MKDIR := mkdir -p
	CP := cp -r
	SEP := /
endif

# Directories
BIN_DIR := bin$(SEP)$(OS_NAME)
TEMPLATES_DIR := templates
DIST_DIR := dist$(SEP)$(OS_NAME)

# Set output binary based on OS
ifeq ($(OS),Windows_NT)
	BINARY=$(BUILD_DIR)/$(APP_NAME).exe
	COPY=cmd /C "if not exist $(BUILD_DIR) mkdir $(BUILD_DIR) & xcopy /E /I /Y $(TEMPLATE_DIR) $(BUILD_DIR)\$(TEMPLATE_DIR)"
else
	BINARY=$(BUILD_DIR)/$(APP_NAME)
	COPY=mkdir -p $(BUILD_DIR)/$(TEMPLATE_DIR) && cp -r $(TEMPLATE_DIR)/* $(BUILD_DIR)/$(TEMPLATE_DIR)
endif

.PHONY: all build clean run copy-templates

all: clean build copy-templates

build:
	@echo "==> Building $(APP_NAME) for $(OS_NAME)..."
	@$(MKDIR) "$(BIN_DIR)"
	@go build -o "$(BIN_DIR)$(SEP)$(APP_NAME).exe" $(SRC_DIR)

copy-templates:
	@echo "==> Copying templates directory..."
	@$(MKDIR) "$(DIST_DIR)"
	@$(CP) "$(TEMPLATES_DIR)" "$(DIST_DIR)$(SEP)templates"

run: build copy-templates
	@echo "==> Running application..."
	@$(BINARY)

clean:
	@echo "==> Cleaning build directories..."
	@if exist "$(BIN_DIR)" $(RM) "$(BIN_DIR)"
	@if exist "$(DIST_DIR)" $(RM) "$(DIST_DIR)"
