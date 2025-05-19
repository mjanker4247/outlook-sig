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

# Cross-compilation settings
CROSS_BUILD_DIR := $(BUILD_DIR)/windows
CROSS_BINARY := $(CROSS_BUILD_DIR)$(SEP)$(APP_NAME).exe

# Signing configuration
CERT_FILE ?= code-sign-certificate.pfx
CERT_URL ?= https://www.program-website.com
SIGNED_BINARY := $(CROSS_BUILD_DIR)$(SEP)$(APP_NAME)-signed.exe
NATIVE_SIGNED_BINARY := $(BUILD_DIR)$(SEP)$(APP_NAME)-signed.exe

.PHONY: all build clean run copy-templates cross-build-windows sign-windows sign-native-windows

all: clean build copy-templates

build:
	@echo "==> Building $(APP_NAME) for $(OS_NAME)..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@go build -o "$(BINARY)" $(SRC_DIR)

copy-templates:
	@echo "==> Copying templates directory..."
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@$(CP) "$(TEMPLATES_DIR)" "$(BUILD_DIR)$(SEP)templates"

run: build copy-templates
	@echo "==> Running application..."
	@$(BINARY)

clean:
	@echo "==> Cleaning build directory..."
	@if exist "$(BUILD_DIR)" $(RM) "$(BUILD_DIR)"

cross-build-windows:
	@echo "==> Cross-compiling $(APP_NAME) for Windows..."
	@if not exist "$(CROSS_BUILD_DIR)" mkdir "$(CROSS_BUILD_DIR)"
	@GOOS=windows GOARCH=amd64 go build -o "$(CROSS_BINARY)" $(SRC_DIR)
	@echo "==> Copying templates for Windows build..."
	@$(CP) "$(TEMPLATES_DIR)" "$(CROSS_BUILD_DIR)$(SEP)templates"
	@echo "==> Windows build complete: $(CROSS_BINARY)"

sign-windows: cross-build-windows
	@echo "==> Signing cross-compiled Windows executable..."
	@if [ ! -f "$(CERT_FILE)" ]; then \
		echo "Error: Certificate file $(CERT_FILE) not found"; \
		exit 1; \
	fi
	@osslsigncode sign \
		-pkcs12 "$(CERT_FILE)" \
		-askpass \
		-n "$(APP_NAME)" \
		-i "$(CERT_URL)" \
		-in "$(CROSS_BINARY)" \
		-out "$(SIGNED_BINARY)"
	@echo "==> Signed executable created: $(SIGNED_BINARY)"

sign-native-windows: build
	@echo "==> Signing native Windows executable..."
	@if not exist "$(CERT_FILE)" ( \
		echo Error: Certificate file $(CERT_FILE) not found && \
		exit /b 1 \
	)
	@osslsigncode sign \
		-pkcs12 "$(CERT_FILE)" \
		-askpass \
		-n "$(APP_NAME)" \
		-i "$(CERT_URL)" \
		-in "$(BINARY)" \
		-out "$(NATIVE_SIGNED_BINARY)"
	@echo "==> Signed executable created: $(NATIVE_SIGNED_BINARY)"
