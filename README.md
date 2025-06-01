# GoCommit

An AI-powered git commit message generator using Google's Gemini API. This tool automatically generates meaningful commit messages based on your staged changes, following conventional commit formats, with comprehensive logging for prompt analysis and improvement.

## Features

- Automatically analyzes git diff of staged changes
- Generates meaningful commit messages using Gemini AI
- Follows conventional commit format (type(scope): description)
- Ensures commit message quality and consistency
- Secure API key configuration and validation
- Available as a system-wide command
- Easy installation with automated script
- One-command installation from the internet
- Comprehensive request/response logging for prompt improvement
- Built-in log analysis tools

## Prerequisites

- Go 1.21 or higher
- Git (optional, for manual installation)
- Gemini API key (starts with "AIza" and is 39 characters long)

## Installation

### Quick Installation (Recommended)

#### Linux / macOS / FreeBSD
```bash
curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh | bash
```

#### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.ps1" -OutFile "install.ps1"; .\install.ps1
```

### Manual Installation

#### Download Pre-built Binaries
Visit the [releases page](https://github.com/thanhphuchuynh/gocommit/releases/latest) and download the appropriate binary:

- **Linux x64**: `gocommit-linux-amd64`
- **Linux ARM64**: `gocommit-linux-arm64`
- **macOS x64**: `gocommit-darwin-amd64`
- **macOS ARM64** (Apple Silicon): `gocommit-darwin-arm64`
- **Windows x64**: `gocommit-windows-amd64.exe`
- **FreeBSD x64**: `gocommit-freebsd-amd64`

#### Build from Source
```bash
git clone https://github.com/thanhphuchuynh/gocommit.git
cd gocommit
go build -o gocommit
```

For detailed installation instructions, see [INSTALL.md](INSTALL.md).

## Updating

To update to the latest version, simply re-run the installation command:

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh | bash
```

**Windows:**
```powershell
.\install.ps1
```

Or download the latest binary from the [releases page](https://github.com/thanhphuchuynh/gocommit/releases/latest).

## Configuration

Configure your Gemini API key using the built-in configuration tool:

```bash
gocommit --config
```

The tool will:
1. Prompt you to enter your API key
2. Validate the key format (must start with "AIza" and be 39 characters long)
3. Securely save the key to your home directory

To update your API key later, simply run the configuration command again.

## Usage

1. Stage your changes using git add:
```bash
git add .
```

2. Run gocommit:
```bash
gocommit
```

The tool will:
1. Analyze your staged changes
2. Generate an appropriate commit message using Gemini AI
3. Create a commit with the generated message

## Commit Message Format

The generated commit messages follow the conventional commits format:

- feat: A new feature
- fix: A bug fix
- docs: Documentation changes
- style: Code style changes (formatting, missing semi-colons, etc)
- refactor: Code refactoring
- perf: Performance improvements
- test: Adding or updating tests
- chore: Maintenance tasks

## Logging & Analysis

gocommit automatically logs all requests and responses to help improve prompts over time. Each interaction captures:

- Git diff and commit context
- AI prompts sent and responses received
- User selections and final commit messages
- Success/failure status and error details
- Timestamps for analysis

Logs are stored locally at `~/.gocommit/gocommit_requests.log` in JSON format.

### Analyzing Logs

Use the built-in analysis tool to understand usage patterns:

```bash
go run tools/analyze_logs.go
```

This provides insights into:
- Success rates and error patterns
- Most common commit types
- User choice preferences
- Daily activity patterns

For detailed logging documentation, see [README_LOGGING.md](README_LOGGING.md).

## License

MIT
