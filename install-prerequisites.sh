#!/bin/bash
set -euo pipefail

echo -e "\033[32mInstalling prerequisites for outlook-signature project...\033[0m"

# ── Homebrew ────────────────────────────────────────────────────────────────
if ! command -v brew &>/dev/null; then
    echo -e "\033[33mInstalling Homebrew...\033[0m"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    if [[ $(uname -m) == 'arm64' ]]; then
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
fi

# ── Core tools ──────────────────────────────────────────────────────────────
if ! command -v git &>/dev/null; then
    echo -e "\033[33mInstalling Git...\033[0m"
    brew install git
fi

if ! command -v go &>/dev/null || [[ ! "$(go version)" =~ "go1.24" ]]; then
    echo -e "\033[33mInstalling Go 1.24.2...\033[0m"
    brew install go@1.24.2
    echo 'export PATH="/usr/local/opt/go@1.24.2/bin:$PATH"' >> ~/.zshrc
    export PATH="/usr/local/opt/go@1.24.2/bin:$PATH"
fi

if ! command -v task &>/dev/null; then
    echo -e "\033[33mInstalling Task...\033[0m"
    brew install go-task
fi

# ── Linting & formatting ────────────────────────────────────────────────────
if ! command -v golangci-lint &>/dev/null; then
    echo -e "\033[33mInstalling golangci-lint...\033[0m"
    brew install golangci-lint
fi

echo -e "\033[33mInstalling goimports...\033[0m"
go install golang.org/x/tools/cmd/goimports@latest

# ── Cross-compilation (Docker + fyne-cross) ─────────────────────────────────
if ! command -v docker &>/dev/null; then
    echo -e "\033[33mInstalling Docker Desktop...\033[0m"
    brew install --cask docker
    echo -e "\033[36m  NOTE: Start Docker Desktop and complete setup before using cross-compilation tasks.\033[0m"
fi

echo -e "\033[33mInstalling fyne-cross...\033[0m"
go install github.com/fyne-io/fyne-cross@latest

# ── Signing tools ───────────────────────────────────────────────────────────
if ! command -v osslsigncode &>/dev/null; then
    echo -e "\033[33mInstalling osslsigncode (cross-platform Windows signing)...\033[0m"
    brew install osslsigncode
fi

# codesign and xcrun are provided by Xcode Command Line Tools
if ! command -v codesign &>/dev/null; then
    echo -e "\033[33mInstalling Xcode Command Line Tools (provides codesign, xcrun)...\033[0m"
    xcode-select --install || true
fi

echo ""
echo -e "\033[32mAll prerequisites installed!\033[0m"
echo -e "\033[32mRun 'task build' to build the project.\033[0m"
