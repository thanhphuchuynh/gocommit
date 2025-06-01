# GoCommit Installation Script for Windows
# This script downloads and installs the latest release of gocommit

param(
    [string]$InstallDir = "$env:USERPROFILE\bin",
    [switch]$Help
)

# Show help
if ($Help) {
    Write-Host "GoCommit Installation Script for Windows"
    Write-Host "Usage: .\install.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -InstallDir DIR    Install directory (default: $env:USERPROFILE\bin)"
    Write-Host "  -Help              Show this help message"
    exit 0
}

# Configuration
$Repo = "thanhphuchuynh/gocommit"
$BinaryName = "gocommit.exe"

# Colors for output
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

# Get latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest release information..."
    
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $version = $response.tag_name
        
        if (-not $version) {
            Write-Error "Could not fetch the latest version"
            exit 1
        }
        
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch release information: $($_.Exception.Message)"
        exit 1
    }
}

# Download binary
function Download-Binary {
    param([string]$Version, [string]$Architecture)
    
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/gocommit-windows-$Architecture.exe"
    $tempFile = "$env:TEMP\$BinaryName"
    
    Write-Info "Downloading from: $downloadUrl"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        
        if (-not (Test-Path $tempFile)) {
            Write-Error "Download failed"
            exit 1
        }
        
        Write-Success "Downloaded successfully"
        return $tempFile
    }
    catch {
        Write-Error "Download failed: $($_.Exception.Message)"
        exit 1
    }
}

# Install binary
function Install-Binary {
    param([string]$TempFile, [string]$InstallPath)
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        Write-Info "Creating install directory: $InstallDir"
        try {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        catch {
            Write-Error "Failed to create install directory: $($_.Exception.Message)"
            exit 1
        }
    }
    
    $finalPath = Join-Path $InstallDir $BinaryName
    
    try {
        Copy-Item $TempFile $finalPath -Force
        Write-Success "Installed $BinaryName to $finalPath"
        
        # Cleanup
        Remove-Item $TempFile -Force
        
        return $finalPath
    }
    catch {
        Write-Error "Installation failed: $($_.Exception.Message)"
        exit 1
    }
}

# Add to PATH
function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to PATH..."
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Added to PATH. Please restart your terminal or PowerShell session."
        return $true
    }
    else {
        Write-Info "Directory already in PATH"
        return $false
    }
}

# Verify installation
function Test-Installation {
    param([string]$BinaryPath)
    
    try {
        $version = & $BinaryPath --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Installation verified! Version: $version"
        }
        else {
            Write-Success "Binary installed successfully"
        }
        
        Write-Info "You can now use 'gocommit' from PowerShell"
    }
    catch {
        Write-Warning "Could not verify installation, but binary was copied successfully"
        Write-Info "Try running: $BinaryPath"
    }
}

# Main installation process
function Main {
    Write-Host "GoCommit Installation Script for Windows" -ForegroundColor Cyan
    Write-Host "=========================================" -ForegroundColor Cyan
    Write-Host ""
    
    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 3) {
        Write-Error "PowerShell 3.0 or higher is required"
        exit 1
    }
    
    $architecture = Get-Architecture
    Write-Info "Detected architecture: $architecture"
    
    $version = Get-LatestVersion
    $tempFile = Download-Binary -Version $version -Architecture $architecture
    $installedPath = Install-Binary -TempFile $tempFile -InstallPath $InstallDir
    
    # Add to PATH
    $pathAdded = Add-ToPath -Directory $InstallDir
    
    # Test installation
    if ($pathAdded) {
        Write-Warning "PATH updated. Please restart your terminal and run 'gocommit --help'"
    }
    else {
        Test-Installation -BinaryPath $installedPath
    }
    
    Write-Host ""
    Write-Success "GoCommit installation completed!"
    Write-Host ""
    Write-Host "Next steps:"
    Write-Host "  - Run 'gocommit --help' to see available commands"
    Write-Host "  - Visit https://github.com/$Repo for documentation"
    
    if ($pathAdded) {
        Write-Host "  - Restart your PowerShell/Command Prompt to use 'gocommit' command"
    }
}

# Run main function
try {
    Main
}
catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}
