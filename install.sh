#!/bin/bash

# GoCommit Installation Script
# This script downloads and installs the latest release of gocommit

set -e

# Default values
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
REPO="thanhphuchuynh/gocommit"
BINARY_NAME="gocommit"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os arch
    
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case "$os" in
        linux*)
            os="linux"
            ;;
        darwin*)
            os="darwin"
            ;;
        freebsd*)
            os="freebsd"
            ;;
        mingw*|msys*|cygwin*)
            os="windows"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    PLATFORM="${os}"
    ARCH="${arch}"
    
    if [ "$os" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    log_info "Detected platform: $PLATFORM-$ARCH"
}

# Get latest release version
get_latest_version() {
    log_info "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "curl or wget is required to download the binary"
        exit 1
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        log_error "Could not fetch the latest version"
        exit 1
    fi
    
    log_info "Latest version: $LATEST_VERSION"
}

# Download binary
download_binary() {
    local download_url temp_file
    
    download_url="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME%.*}-$PLATFORM-$ARCH"
    if [ "$PLATFORM" = "windows" ]; then
        download_url="${download_url}.exe"
    fi
    
    temp_file="/tmp/$BINARY_NAME"
    
    log_info "Downloading from: $download_url"
    
    if command -v curl >/dev/null 2>&1; then
        if curl -L -o "$temp_file" "$download_url" --silent --show-error; then
            log_success "Downloaded successfully"
        else
            log_error "Download failed"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if wget -O "$temp_file" "$download_url" --quiet; then
            log_success "Downloaded successfully"
        else
            log_error "Download failed"
            exit 1
        fi
    else
        log_error "curl or wget is required to download the binary"
        exit 1
    fi
    
    if [ ! -f "$temp_file" ]; then
        log_error "Download failed - file not found"
        exit 1
    fi
    
    echo "$temp_file"
}

# Install binary
install_binary() {
    local temp_file="$1"
    local install_path="$INSTALL_DIR/$BINARY_NAME"
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Creating install directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR" || {
            log_error "Failed to create install directory. Try running with sudo or choose a different directory."
            exit 1
        }
    fi
    
    # Check if we have write permissions
    if [ ! -w "$INSTALL_DIR" ]; then
        log_warning "No write permission to $INSTALL_DIR. Trying with sudo..."
        sudo cp "$temp_file" "$install_path"
        sudo chmod +x "$install_path"
    else
        cp "$temp_file" "$install_path"
        chmod +x "$install_path"
    fi
    
    # Cleanup
    rm -f "$temp_file"
    
    log_success "Installed $BINARY_NAME to $install_path"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        log_success "Installation verified! Version: $version"
        log_info "You can now use '$BINARY_NAME' from anywhere in your terminal"
    else
        log_warning "Binary installed but not found in PATH"
        log_info "You may need to add $INSTALL_DIR to your PATH or restart your terminal"
        log_info "Or run the binary directly: $INSTALL_DIR/$BINARY_NAME"
    fi
}

# Main installation process
main() {
    echo "GoCommit Installation Script"
    echo "============================"
    echo
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --install-dir DIR    Install directory (default: /usr/local/bin)"
                echo "  --help, -h           Show this help message"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    detect_platform
    get_latest_version
    
    local temp_file
    temp_file=$(download_binary)
    
    if [ -z "$temp_file" ] || [ ! -f "$temp_file" ]; then
        log_error "Failed to download binary"
        exit 1
    fi
    
    install_binary "$temp_file"
    verify_installation
    
    echo
    log_success "GoCommit installation completed!"
    echo
    echo "Next steps:"
    echo "  - Run 'gocommit --help' to see available commands"
    echo "  - Visit https://github.com/$REPO for documentation"
}

# Run main function
main "$@"
