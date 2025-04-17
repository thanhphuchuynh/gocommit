#!/bin/bash

# Exit on error
set -e

echo "Installing GoCommit..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ "$(printf '%s\n' "1.21" "$GO_VERSION" | sort -V | head -n1)" != "1.21" ]]; then
    echo "Error: Go version 1.21 or higher is required. Current version: $GO_VERSION"
    exit 1
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download the repository
echo "Downloading GoCommit..."
if command -v git &> /dev/null; then
    git clone https://github.com/thanhphuchuynh/gocommit.git .
else
    echo "Git not found. Downloading using curl/wget..."
    if command -v curl &> /dev/null; then
        curl -L https://github.com/thanhphuchuynh/gocommit/archive/refs/heads/main.tar.gz | tar xz --strip-components=1
    elif command -v wget &> /dev/null; then
        wget -qO- https://github.com/thanhphuchuynh/gocommit/archive/refs/heads/main.tar.gz | tar xz --strip-components=1
    else
        echo "Error: Neither curl nor wget is installed. Please install one of them."
        exit 1
    fi
fi

# Build the application
echo "Building GoCommit..."
go build

# Check if build was successful
if [ ! -f "gocommit" ]; then
    echo "Error: Build failed. gocommit binary not found."
    exit 1
fi

# Determine installation directory
INSTALL_DIR="/usr/local/bin"
if [ "$(uname)" == "Darwin" ]; then
    # macOS
    INSTALL_DIR="/usr/local/bin"
elif [ "$(uname)" == "Linux" ]; then
    # Linux
    INSTALL_DIR="/usr/local/bin"
else
    echo "Error: Unsupported operating system."
    exit 1
fi

# Install the binary
echo "Installing to $INSTALL_DIR..."
sudo mv gocommit "$INSTALL_DIR/"

# Set executable permissions
sudo chmod +x "$INSTALL_DIR/gocommit"

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"

echo "Installation complete! You can now use 'gocommit' from anywhere."
echo "To configure your API key, run: gocommit --config" 
