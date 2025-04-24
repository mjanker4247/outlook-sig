# Makefile for Windows and Mac/Linux

# Application name
APP_NAME := SignatureInstaller

# Directories
BUILD_DIR := build
TEMPLATES_DIR := templates
SRC_DIR := ./cmd/signature-installer

# Determine OS
ifeq ($(OS),Windows_NT)
	OS_NAME := windows
	RM := del /Q /F
	MKDIR := if not exist
	CP := xcopy /E /I /Y
	SEP := \\
	BINARY := $(BUILD_DIR)$(SEP)$(APP_NAME).exe
else
	OS_NAME := $(shell uname -s | tr '[:upper:]' '[:lower:]')
	RM := rm -rf
	MKDIR := mkdir -p
	CP := cp -r
	SEP := /
	BINARY := $(BUILD_DIR)$(SEP)$(APP_NAME)
endif

.PHONY: all build clean run copy-templates

all: clean build copy-templates

build:
	@echo "==> Building $(APP_NAME) for $(OS_NAME)..."
	@$(MKDIR) "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@go build -o "$(BINARY)" $(SRC_DIR)

copy-templates:
	@echo "==> Copying templates directory..."
	@$(MKDIR) "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@$(CP) "$(TEMPLATES_DIR)" "$(BUILD_DIR)$(SEP)templates"

run: build copy-templates
	@echo "==> Running application..."
	@$(BINARY)

clean:
	@echo "==> Cleaning build directory..."
	@if exist "$(BUILD_DIR)" $(RM) "$(BUILD_DIR)"
