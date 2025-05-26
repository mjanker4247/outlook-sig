#!/bin/bash

echo -e "\033[32mInstalling prerequisites for outlook-signature project...\033[0m"

# Check if Homebrew is installed, if not install it
if ! command -v brew &> /dev/null; then
    echo -e "\033[33mInstalling Homebrew package manager...\033[0m"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Add Homebrew to PATH for Apple Silicon Macs
    if [[ $(uname -m) == 'arm64' ]]; then
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
fi

# Install Git if not present
if ! command -v git &> /dev/null; then
    echo -e "\033[33mInstalling Git...\033[0m"
    brew install git
fi

# Install Go if not present or wrong version
if ! command -v go &> /dev/null || [[ ! "$(go version)" =~ "go1.24.2" ]]; then
    echo -e "\033[33mInstalling Go 1.24.2...\033[0m"
    brew install go@1.24.2
    echo 'export PATH="/usr/local/opt/go@1.24.2/bin:$PATH"' >> ~/.zshrc
    source ~/.zshrc
fi

# Install Task if not present
if ! command -v task &> /dev/null; then
    echo -e "\033[33mInstalling Task...\033[0m"
    brew install go-task/tap/go-task
fi

echo -e "\n\033[32mAll prerequisites have been installed!\033[0m"
echo -e "\033[32mYou can now run 'task build' to build the project.\033[0m" 