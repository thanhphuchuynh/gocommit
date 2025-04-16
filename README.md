# GoCommit

An AI-powered git commit message generator using Google's Gemini API. This tool automatically generates meaningful commit messages based on your staged changes, following conventional commit formats.

## Features

- Automatically analyzes git diff of staged changes
- Generates meaningful commit messages using Gemini AI
- Follows conventional commit format (type(scope): description)
- Ensures commit message quality and consistency

## Prerequisites

- Go 1.21 or higher
- Git
- Gemini API key

## Installation

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

## Configuration

Set your Gemini API key as an environment variable:

```bash
export GEMINI_API_KEY='your-api-key-here'
```

## Usage

1. Stage your changes using git add:
```bash
git add .
```

2. Run gocommit:
```bash
./gocommit
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