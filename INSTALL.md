# GoCommit Installation Guide

This guide provides multiple ways to install GoCommit on your system.

## Quick Installation

### Linux / macOS / FreeBSD

**One-line installation (recommended):**
```bash
curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh | bash
```

**Or download and run:**
```bash
wget https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh
chmod +x install.sh
./install.sh
```

### Windows

**PowerShell (recommended):**
```powershell
# Run in PowerShell as Administrator (optional, for system-wide install)
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

**Or download and run:**
1. Download [install.ps1](https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.ps1)
2. Right-click and "Run with PowerShell"

## Installation Options

### Custom Installation Directory

**Linux/macOS:**
```bash
# Install to custom directory
./install.sh --install-dir ~/.local/bin

# Install to system-wide location (requires sudo)
sudo ./install.sh --install-dir /usr/local/bin
```

**Windows:**
```powershell
# Install to custom directory
.\install.ps1 -InstallDir "C:\tools\bin"
```

### Manual Installation

#### 1. Download Pre-built Binaries

Visit the [releases page](https://github.com/thanhphuchuynh/gocommit/releases/latest) and download the appropriate binary for your system:

- **Linux x64**: `gocommit-linux-amd64`
- **Linux ARM64**: `gocommit-linux-arm64`
- **macOS x64**: `gocommit-darwin-amd64`
- **macOS ARM64** (Apple Silicon): `gocommit-darwin-arm64`
- **Windows x64**: `gocommit-windows-amd64.exe`
- **FreeBSD x64**: `gocommit-freebsd-amd64`

#### 2. Verify Download (Optional but Recommended)

Download the `checksums.txt` file and verify your binary:

**Linux/macOS:**
```bash
# Download checksums
wget https://github.com/thanhphuchuynh/gocommit/releases/latest/download/checksums.txt

# Verify binary (example for Linux x64)
sha256sum -c checksums.txt --ignore-missing
```

**Windows (PowerShell):**
```powershell
# Get checksum of downloaded file
Get-FileHash -Algorithm SHA256 gocommit-windows-amd64.exe

# Compare with checksums.txt manually
```

#### 3. Install Binary

**Linux/macOS/FreeBSD:**
```bash
# Make executable
chmod +x gocommit-*

# Move to PATH (choose one)
sudo mv gocommit-* /usr/local/bin/gocommit        # System-wide
mv gocommit-* ~/.local/bin/gocommit               # User-only
mv gocommit-* ~/bin/gocommit                      # User bin directory
```

**Windows:**
```powershell
# Create bin directory in user profile
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin"

# Move binary
Move-Item gocommit-windows-amd64.exe "$env:USERPROFILE\bin\gocommit.exe"

# Add to PATH (user)
$path = [Environment]::GetEnvironmentVariable("PATH", "User")
$newPath = "$path;$env:USERPROFILE\bin"
[Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
```

## Build from Source

### Prerequisites

- Go 1.21 or later
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/thanhphuchuynh/gocommit.git
cd gocommit

# Build
go build -o gocommit

# Install (optional)
go install
```

### Cross-compilation

Build for different platforms:

```bash
# Linux x64
GOOS=linux GOARCH=amd64 go build -o gocommit-linux-amd64

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o gocommit-darwin-arm64

# Windows x64
GOOS=windows GOARCH=amd64 go build -o gocommit-windows-amd64.exe
```

## Package Managers

### Homebrew (macOS/Linux)

*Coming soon...*

### Chocolatey (Windows)

*Coming soon...*

### APT (Ubuntu/Debian)

*Coming soon...*

### AUR (Arch Linux)

*Coming soon...*

## Verification

After installation, verify that GoCommit is working:

```bash
# Check version
gocommit --version

# Show help
gocommit --help

# Check installation path
which gocommit  # Linux/macOS
where gocommit  # Windows
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. **Restart your terminal/PowerShell**
2. **Check PATH**: Ensure the installation directory is in your PATH
3. **Manual PATH update**:

**Linux/macOS (add to ~/.bashrc or ~/.zshrc):**
```bash
export PATH="$HOME/.local/bin:$PATH"  # If installed to ~/.local/bin
export PATH="$HOME/bin:$PATH"         # If installed to ~/bin
```

**Windows (PowerShell profile):**
```powershell
# Add to $PROFILE
$env:PATH += ";$env:USERPROFILE\bin"
```

### Permission Denied

**Linux/macOS:**
```bash
# Make sure binary is executable
chmod +x /path/to/gocommit

# For system-wide installation, use sudo
sudo cp gocommit /usr/local/bin/
```

**Windows:**
- Run PowerShell as Administrator for system-wide installation
- Check antivirus software (may block unsigned executables)

### Download Issues

1. **Check internet connection**
2. **Firewall/proxy**: Ensure GitHub releases are accessible
3. **Manual download**: Use browser to download from releases page

### Platform Not Supported

If your platform isn't supported:
1. Try building from source
2. Open an issue on GitHub for platform support request

## Updating

### Using Installation Scripts

Re-run the installation script to get the latest version:

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh | bash

# Windows
.\install.ps1
```

### Manual Update

1. Download the latest binary from releases
2. Replace the existing binary
3. Verify the new version: `gocommit --version`

## Uninstallation

### Remove Binary

**Linux/macOS:**
```bash
# System-wide installation
sudo rm /usr/local/bin/gocommit

# User installation
rm ~/.local/bin/gocommit
# or
rm ~/bin/gocommit
```

**Windows:**
```powershell
# Remove binary
Remove-Item "$env:USERPROFILE\bin\gocommit.exe"

# Remove from PATH (optional)
# Edit PATH environment variable manually or use system settings
```

### Clean Up Configuration

GoCommit stores configuration in:
- Linux/macOS: `~/.config/gocommit/`
- Windows: `%APPDATA%\gocommit\`

Remove these directories if you want to completely uninstall GoCommit.

## Support

If you encounter issues:

1. Check this troubleshooting guide
2. Search [existing issues](https://github.com/thanhphuchuynh/gocommit/issues)
3. Create a [new issue](https://github.com/thanhphuchuynh/gocommit/issues/new) with:
   - Operating system and architecture
   - Installation method used
   - Error messages (if any)
   - Output of `gocommit --version` (if working)
