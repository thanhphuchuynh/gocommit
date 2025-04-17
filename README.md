# GoCommit

An AI-powered git commit message generator using Google's Gemini API. This tool automatically generates meaningful commit messages based on your staged changes, following conventional commit formats.

## Features

- Automatically analyzes git diff of staged changes
- Generates meaningful commit messages using Gemini AI
- Follows conventional commit format (type(scope): description)
- Ensures commit message quality and consistency
- Secure API key configuration and validation
- Available as a system-wide command
- Easy installation with automated script
- One-command installation from the internet

## Prerequisites

- Go 1.21 or higher
- Git (optional, for manual installation)
- Gemini API key (starts with "AIza" and is 39 characters long)

## Installation

### Option 1: One-command Installation (Recommended)
Run this command in your terminal:
```bash
curl -sSL https://raw.githubusercontent.com/tphuc/gocommit/main/install.sh | bash
```
or using wget:
```bash
wget -qO- https://raw.githubusercontent.com/tphuc/gocommit/main/install.sh | bash
```

The script will:
1. Download the latest version
2. Check for required dependencies
3. Build the application
4. Install it system-wide
5. Set proper permissions

### Option 2: Manual Installation
1. Clone the repository:
```bash
git clone https://github.com/tphuc/gocommit.git
cd gocommit
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build
```

4. Move the binary to a directory in your PATH:
```bash
# For macOS/Linux
sudo mv gocommit /usr/local/bin/

# For Windows (using PowerShell as admin)
Move-Item gocommit.exe C:\Windows\System32\
```

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

## License

MIT 