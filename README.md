# GoCommit

An AI-powered git commit message generator supporting multiple AI providers (Google Gemini and OpenRouter). This tool automatically generates meaningful commit messages based on your staged changes, following conventional commit formats, with comprehensive logging for prompt analysis and improvement.

## Features

- Automatically analyzes git diff of staged changes
- Generates meaningful commit messages using AI (Gemini or OpenRouter)
- Supports multiple AI models through OpenRouter (Claude, GPT-4, Llama, etc.)
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
- API key from one of the supported providers:
  - **Gemini**: API key starting with "AIza" (39 characters)
  - **OpenRouter**: API key starting with "sk-or-" or "sk-"

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

### Interactive Setup

Configure your AI provider and API key using the built-in interactive configuration tool:

```bash
gocommit --config
```

The tool will guide you through:

1. **Select AI Provider**:
   - Choose between Gemini (Google) or OpenRouter (Claude, GPT-4, Llama, etc.)

2. **Enter API Key**:
   - For Gemini: Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
   - For OpenRouter: Get your API key from [openrouter.ai](https://openrouter.ai/)

3. **Configure Model** (OpenRouter only):
   - Optionally select a specific model or use the default (Claude 3.5 Sonnet)

**Example Configuration Flow:**

```bash
$ gocommit --config

=== GoCommit Configuration ===

Select AI Provider:
  1. Gemini (Google)
  2. OpenRouter (Claude, GPT-4, Llama, etc.)

Enter choice (1 or 2) [1]: 2

Enter OpenRouter API key: sk-or-v1-xxx...

Configure OpenRouter model (optional):
  Popular models:
    - anthropic/claude-3.5-sonnet (default)
    - anthropic/claude-3-opus
    - openai/gpt-4-turbo
    - openai/gpt-4
    - meta-llama/llama-3.1-70b-instruct

Enter model name (or press Enter for default): openai/gpt-4-turbo

✓ Configuration saved successfully!
  Provider: openrouter
  Model: openai/gpt-4-turbo

You can now use 'gocommit' to generate commit messages.
```

### Available OpenRouter Models

When using OpenRouter, you have access to various AI models:

- **Anthropic Claude**:
  - `anthropic/claude-3.5-sonnet` (default, recommended)
  - `anthropic/claude-3-opus`
  - `anthropic/claude-3-sonnet`

- **OpenAI**:
  - `openai/gpt-4-turbo`
  - `openai/gpt-4`
  - `openai/gpt-3.5-turbo`

- **Google**:
  - `google/gemini-pro-1.5`
  - `google/gemini-pro`

- **Meta**:
  - `meta-llama/llama-3.1-70b-instruct`
  - `meta-llama/llama-3.1-405b-instruct`

- **And many more** - see [OpenRouter models](https://openrouter.ai/models) for the complete list

### Changing Configuration

To switch providers or update your API key, simply run the configuration command again:

```bash
gocommit --config
```

You can also change the model for OpenRouter without reconfiguring everything:

```bash
gocommit --set-model "anthropic/claude-3-opus"
```

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

For detailed logging documentation, see [README_LOGGING_CONFIG.md](README_LOGGING_CONFIG.md).

## Project Documentation

- **[FLOW.md](FLOW.md)** - Complete project workflow and architecture documentation
- **[INSTALL.md](INSTALL.md)** - Detailed installation instructions
- **[README_GITHUB_ACTIONS.md](README_GITHUB_ACTIONS.md)** - CI/CD documentation
- **[README_LOGGING_CONFIG.md](README_LOGGING_CONFIG.md)** - Logging configuration

## License

MIT
