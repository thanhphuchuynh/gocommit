# GoCommit Project Flow Documentation

This document explains how the GoCommit project works, from development to deployment and usage.

## 📋 Project Overview

GoCommit is an AI-powered git commit message generator that uses Google's Gemini API to analyze staged changes and generate meaningful commit messages following conventional commit formats.

## 🏗️ Project Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          GoCommit Architecture                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │    main.go  │────│  config/    │────│     logger/         │  │
│  │   (CLI)     │    │ config.go   │    │   logger.go         │  │
│  │             │    │ (API keys)  │    │ (logging system)    │  │
│  └─────────────┘    └─────────────┘    └─────────────────────┘  │
│         │                                         │              │
│         │                                         ▼              │
│         ▼                                  ┌─────────────────┐   │
│  ┌─────────────┐                          │ ~/.gocommit/    │   │
│  │   Git Diff  │                          │ gocommit_       │   │
│  │  Analysis   │                          │ requests.log    │   │
│  └─────────────┘                          └─────────────────┘   │
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │              Google Gemini AI API                           ││
│  │         (Commit Message Generation)                         ││
│  └─────────────────────────────────────────────────────────────┘│
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────┐                                                 │
│  │ Git Commit  │                                                 │
│  │  Execution  │                                                 │
│  └─────────────┘                                                 │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 🔄 Development Workflow

### 1. Code Development Flow

```mermaid
graph TD
    A[Developer writes code] --> B[git add .]
    B --> C[gocommit]
    C --> D[Analyze git diff]
    D --> E[Send to Gemini AI]
    E --> F[Generate commit message]
    F --> G[Execute git commit]
    G --> H[Log interaction]
```

### 2. CI/CD Pipeline Flow

```mermaid
graph LR
    A[Push to main] --> B[Build Workflow]
    A --> C[CI Workflow]
    D[Create Tag] --> E[Release Workflow]
    
    B --> B1[Ubuntu Build]
    B --> B2[macOS Build]
    B --> B3[Run Tests]
    
    C --> C1[Linting]
    C --> C2[Coverage]
    C --> C3[Security Scan]
    
    E --> E1[Multi-platform Build]
    E --> E2[Generate Checksums]
    E --> E3[Create Release]
    E --> E4[Upload Binaries]
```

## 🚀 Installation Flow

### Quick Installation Process

```mermaid
graph TD
    A[User runs install command] --> B{Detect Platform}
    B --> C[Linux/macOS/FreeBSD]
    B --> D[Windows]
    
    C --> E[Download install.sh]
    D --> F[Download install.ps1]
    
    E --> G[Detect Architecture]
    F --> G
    
    G --> H[Fetch Latest Release]
    H --> I[Download Binary]
    I --> J[Verify Download]
    J --> K[Install to PATH]
    K --> L[Verify Installation]
    L --> M[Installation Complete]
```

### Platform-Specific Installation

| Platform | Command | Script |
|----------|---------|--------|
| **Linux/macOS** | `curl -sSL https://raw.githubusercontent.com/thanhphuchuynh/gocommit/main/install.sh \| bash` | `install.sh` |
| **Windows** | `Invoke-WebRequest -Uri "..." -OutFile "install.ps1"; .\install.ps1` | `install.ps1` |
| **Manual** | Download from releases page | N/A |

## 📊 Usage Flow

### Basic Usage Workflow

```mermaid
sequenceDiagram
    participant User
    participant GoCommit
    participant Git
    participant Gemini
    participant Logger
    
    User->>Git: git add .
    User->>GoCommit: gocommit
    GoCommit->>Git: git diff --cached
    Git-->>GoCommit: diff output
    GoCommit->>Gemini: Send diff + prompt
    Gemini-->>GoCommit: Generated commit message
    GoCommit->>Logger: Log request/response
    GoCommit->>User: Show commit message
    User->>GoCommit: Confirm
    GoCommit->>Git: git commit -m "message"
    Git-->>GoCommit: Commit success
    GoCommit->>User: Commit completed
```

### Configuration Flow

```mermaid
graph TD
    A[First Run] --> B{API Key Set?}
    B -->|No| C[gocommit --config]
    B -->|Yes| D[Ready to use]
    
    C --> E[Prompt for API key]
    E --> F[Validate key format]
    F --> G{Valid?}
    G -->|No| E
    G -->|Yes| H[Save to ~/.config/gocommit/]
    H --> D
```

## 🔧 Core Components

### 1. Main Application (`main.go`)
- **Purpose**: Command-line interface and orchestration
- **Features**:
  - Command-line argument parsing
  - Git integration
  - User interaction
  - Workflow coordination

### 2. Configuration (`config/config.go`)
- **Purpose**: API key management and validation
- **Features**:
  - Secure API key storage
  - Key format validation
  - Configuration file management

### 3. Logging (`logger/logger.go`)
- **Purpose**: Request/response logging for analysis
- **Features**:
  - JSON structured logging
  - Request/response tracking
  - Error logging
  - Performance metrics

### 4. Build System
- **Makefile**: Local development commands
- **GitHub Actions**: Automated CI/CD
- **Install Scripts**: Cross-platform installation

## 📈 Release Process

### Version Release Flow

```mermaid
graph LR
    A[Create Tag] --> B[Push Tag]
    B --> C[GitHub Actions Trigger]
    C --> D[Build Multi-platform]
    D --> E[Generate Checksums]
    E --> F[Create GitHub Release]
    F --> G[Upload Binaries]
    G --> H[Update Documentation]
```

### Release Artifacts

For each release, the following binaries are built:
- `gocommit-linux-amd64`
- `gocommit-linux-arm64`
- `gocommit-darwin-amd64`
- `gocommit-darwin-arm64`
- `gocommit-windows-amd64.exe`
- `gocommit-freebsd-amd64`
- `checksums.txt`

## 📝 Logging and Analysis

### Log Flow

```mermaid
graph TD
    A[User Action] --> B[GoCommit Execution]
    B --> C[Log Request]
    C --> D[Send to Gemini]
    D --> E[Log Response]
    E --> F[User Choice]
    F --> G[Log Final Result]
    G --> H[Store in ~/.gocommit/gocommit_requests.log]
```

### Log Analysis

```bash
# Analyze logs
go run tools/analyze_logs.go

# Output insights:
# - Success rates
# - Common commit types
# - Usage patterns
# - Error analysis
```

## 🔒 Security Flow

### API Key Security

```mermaid
graph TD
    A[User enters API key] --> B[Validate format]
    B --> C[Store in config file]
    C --> D[Set file permissions 600]
    D --> E[Use in API requests]
    E --> F[Never log API key]
```

### Security Measures
- API keys stored locally only
- Configuration files have restricted permissions
- No sensitive data in logs
- HTTPS-only communication
- Binary checksums for verification

## 📚 Documentation Structure

```
├── README.md              # Main project documentation
├── INSTALL.md             # Installation instructions
├── FLOW.md                # This workflow documentation
├── README_LOGGING_CONFIG.md # Logging configuration
├── README_GITHUB_ACTIONS.md # CI/CD documentation
└── tools/
    └── analyze_logs.go    # Log analysis tool
```

## 🎯 Key Features Flow

### 1. AI Integration
- Analyzes git diff context
- Generates conventional commit messages
- Learns from user preferences (via logs)

### 2. Cross-Platform Support
- Supports Linux, macOS, Windows, FreeBSD
- Architecture detection (amd64, arm64)
- Platform-specific installers

### 3. Developer Experience
- One-command installation
- Simple configuration
- Comprehensive logging
- Easy updates

### 4. Quality Assurance
- Automated testing
- Linting and formatting
- Security scanning
- Multi-platform builds

## 🔄 Update Flow

### Updating GoCommit

```mermaid
graph TD
    A[User wants to update] --> B[Re-run install script]
    B --> C[Fetch latest version]
    C --> D[Download new binary]
    D --> E[Replace existing binary]
    E --> F[Preserve configuration]
    F --> G[Update complete]
```

### Automated Updates (Future)
- Version checking
- Automatic update notifications
- Background updates
- Rollback capability

## 🤝 Contributing Flow

### Development Process

```mermaid
graph LR
    A[Fork Repository] --> B[Create Feature Branch]
    B --> C[Make Changes]
    C --> D[Run Tests]
    D --> E[Submit PR]
    E --> F[CI Checks]
    F --> G[Code Review]
    G --> H[Merge]
```

### Requirements for Contributors
- Go 1.21+ for development
- Follow conventional commits
- Add tests for new features
- Update documentation
- Ensure cross-platform compatibility

---

This flow documentation provides a comprehensive overview of how GoCommit works from development to deployment and usage. For specific implementation details, refer to the individual source files and documentation.
