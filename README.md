# GoCommit

An AI-powered git commit message generator supporting multiple AI providers (Google Gemini, OpenRouter, and Ollama). This tool automatically generates meaningful commit messages based on your staged changes, following conventional commit formats, with comprehensive logging for prompt analysis and improvement.

📚 **[View Full Documentation](https://thanhphuchuynh.github.io/gocommit/)**

## Features

- **Interactive Mode**: Full control over commit messages
  - Accept / Edit / Regenerate / Preview Diff / Quit
  - Split-screen diff view for reviewing changes
  - Regeneration loop for perfect messages
- Automatically analyzes git diff of staged changes
- Generates meaningful commit messages using AI (Gemini, OpenRouter, or Ollama)
- Supports multiple AI models:
  - **Cloud**: OpenRouter (Claude, GPT-4, Llama, etc.)
  - **Local**: Ollama (CodeLlama, Llama3, Mistral, etc.) - 100% free & private
- **Auto mode** (`--auto` or `-y`) for fast commits and CI/CD
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
- **One of the following**:
  - **Gemini**: API key starting with "AIza" (39 characters)
  - **OpenRouter**: API key starting with "sk-or-" or "sk-"
  - **Ollama**: [Local installation](https://ollama.ai) (no API key needed, 100% free)

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
   - **Gemini** (Google) - Cloud-based, requires API key
   - **OpenRouter** - Access to Claude, GPT-4, Llama, etc., requires API key
   - **Ollama** - Local AI, 100% free, no API key needed

2. **Enter API Key** (Gemini/OpenRouter only):
   - For Gemini: Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
   - For OpenRouter: Get your API key from [openrouter.ai](https://openrouter.ai/)
   - For Ollama: No API key needed

3. **Configure Model**:
   - **OpenRouter**: Select model (default: Claude 3.5 Sonnet)
   - **Ollama**: Select local model (default: CodeLlama 7B)

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

### Using Ollama (Local AI)

**Ollama** allows you to run AI models locally on your machine - 100% free, completely private, and works offline!

#### Prerequisites

1. **Install Ollama**: Visit [ollama.ai](https://ollama.ai) and install for your platform
2. **Pull a model**: Download a model of your choice

```bash
# Recommended models for commit messages
ollama pull codellama:7b      # 3.8GB - Good balance (default)
ollama pull llama3:8b          # 4.7GB - Best quality
ollama pull mistral:7b         # 4.1GB - Fast
ollama pull deepseek-coder:6.7b  # 3.8GB - Code-focused
ollama pull codellama:13b      # 7.4GB - Higher quality, slower
```

#### Setup

```bash
# 1. Make sure Ollama is running
ollama serve  # Usually runs automatically after installation

# 2. Configure gocommit for Ollama
gocommit --config

# Select option 3 (Ollama)
# Choose your model or press Enter for default (codellama:7b)
```

#### Configuration Example

```bash
$ gocommit --config

=== GoCommit Configuration ===

Select AI Provider:
  1. Gemini (Google)
  2. OpenRouter (Claude, GPT-4, Llama, etc.)
  3. Ollama (Local AI)

Enter choice (1, 2, or 3) [1]: 3

Configure Ollama:
Enter Ollama endpoint [http://localhost:11434]:

Configure Ollama model:
  Recommended models:
    - codellama:7b (default, 3.8GB)
    - llama3:8b (4.7GB)
    - mistral:7b (4.1GB)
    - deepseek-coder:6.7b (3.8GB)
    - codellama:13b (7.4GB, better quality)

Enter model name [codellama:7b]: llama3:8b

ℹ️  Make sure Ollama is running and the model is pulled:
   ollama pull llama3:8b

✓ Configuration saved successfully!
  Provider: ollama
  Endpoint: http://localhost:11434
  Model: llama3:8b
```

#### Model Comparison

| Model | Size | Quality | Speed | Use Case |
|-------|------|---------|-------|----------|
| `codellama:7b` | 3.8GB | ⭐⭐⭐ | Fast | General commits (default) |
| `llama3:8b` | 4.7GB | ⭐⭐⭐⭐ | Fast | Best balanced option |
| `mistral:7b` | 4.1GB | ⭐⭐⭐ | Very Fast | Quick commits |
| `deepseek-coder:6.7b` | 3.8GB | ⭐⭐⭐⭐ | Medium | Code-focused |
| `codellama:13b` | 7.4GB | ⭐⭐⭐⭐ | Slower | High quality |

#### Benefits of Ollama

- ✅ **$0 Cost** - Completely free forever
- ✅ **100% Private** - Data never leaves your machine
- ✅ **Offline** - Works without internet
- ✅ **Fast** - No network latency
- ✅ **No API Limits** - Unlimited usage
- ✅ **Open Source** - Full transparency

#### Troubleshooting

**Ollama not responding?**
```bash
# Check if Ollama is running
curl http://localhost:11434

# If not, start it
ollama serve
```

**Model not found?**
```bash
# Pull the model first
ollama pull codellama:7b

# List installed models
ollama list
```

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

### Interactive Mode (Default)

1. Stage your changes using git add:
```bash
git add .
```

2. Run gocommit:
```bash
gocommit
```

The tool will:
1. Analyze your staged changes with AI
2. Present 3 AI-generated commit message options
3. Let you select a message
4. **Show interactive prompt with options:**
   - **[A]ccept**: Commit with the selected message
   - **[E]dit**: Open your editor to modify the message
   - **[R]egenerate**: Generate new AI suggestions
   - **[D]iff**: Toggle split-screen view (diff + message)
   - **[Q]uit**: Cancel the commit operation

### Auto Mode (Fast Commits)

Skip the interactive prompt for quick commits or CI/CD workflows:

```bash
# Using --auto flag
gocommit --auto

# Using -y shorthand
gocommit -y
```

In auto mode, gocommit will:
1. Generate commit messages
2. Show selection UI
3. Commit immediately after selection (no interactive prompt)

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
