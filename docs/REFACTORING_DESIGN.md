# GoCommit Refactoring Design Document

## Executive Summary

This document outlines a comprehensive refactoring plan to split main.go (1649 lines) into a modular, maintainable architecture following clean code principles and functional programming patterns where appropriate.

**Current State**: Monolithic 1649-line main.go with mixed concerns  
**Target State**: Modular architecture with clear separation of concerns  
**Approach**: Phased refactoring maintaining backward compatibility  
**Estimated Phases**: 6 phases over multiple iterations

---

## Table of Contents

1. [Current Architecture Analysis](#current-architecture-analysis)
2. [Proposed Architecture](#proposed-architecture)
3. [Module Design](#module-design)
4. [AI Provider Abstraction](#ai-provider-abstraction)
5. [Function Migration Plan](#function-migration-plan)
6. [Dependency Graph](#dependency-graph)
7. [Refactoring Phases](#refactoring-phases)
8. [Testing Strategy](#testing-strategy)
9. [Backward Compatibility](#backward-compatibility)
10. [Implementation Guidelines](#implementation-guidelines)

---

## Current Architecture Analysis

### File Statistics

- **Total Lines**: 1649
- **Functions**: ~25 functions
- **Responsibilities**: 6+ distinct concerns mixed together
- **Complexity**: High coupling, low cohesion

### Current Structure Breakdown

```
main.go (1649 lines)
├── Imports & Constants (1-20)
├── Data Structures (22-26)
├── Prompt Templates (27-283) - 4 templates × ~60 lines each
├── Git Operations (285-531)
│   ├── getGitDiff() - 8 lines
│   ├── generateCommitMessages() - 237 lines (COMPLEX)
│   └── getLastCommitMessage() - 7 lines
├── UI Functions (533-920)
│   ├── drawMessages() - 83 lines
│   ├── editMessage() - 237 lines (COMPLEX)
│   └── getUserChoice() - 68 lines
├── Delayed Commit Logic (922-1428)
│   ├── isTimeInRestrictedRange() - 15 lines
│   ├── parseTimeString() - 102 lines
│   ├── validateHour() - 5 lines
│   ├── generateSuggestedTimes() - 39 lines
│   ├── getNextAvailableTime() - 12 lines
│   ├── showTimeSelectionUI() - 121 lines
│   ├── handleCustomTimeInput() - 145 lines
│   └── executeDelayedCommit() - 44 lines
├── Utilities (1430-1433)
│   └── isValidAPIKey() - 4 lines
└── Main Function (1435-1649)
    └── main() - 215 lines (VERY COMPLEX)
```

### Key Issues Identified

1. **Single Responsibility Violation**: main.go handles AI, Git, UI, time logic, config
2. **Long Functions**: `generateCommitMessages()` (237 lines), `editMessage()` (237 lines)
3. **High Complexity**: main() orchestrates everything (215 lines)
4. **Hard-coded AI Provider**: Gemini API directly coupled
5. **Mixed Abstraction Levels**: Low-level termbox calls mixed with high-level logic
6. **Difficult to Test**: Monolithic structure makes unit testing challenging
7. **Code Duplication**: Prompt templates have similar structure

---

## Proposed Architecture

### Directory Structure

```
/home/tphuc/coding/gocommit/
├── main.go                      # Entry point (80-100 lines)
├── go.mod
├── go.sum
│
├── config/                      # Configuration (EXISTING - enhanced)
│   ├── config.go               # Config management
│   └── config_test.go          # Config tests
│
├── logger/                      # Logging (EXISTING)
│   ├── logger.go               # Request logging
│   └── logger_test.go          # Logger tests
│
├── ai/                          # AI Provider Abstraction (NEW)
│   ├── provider.go             # Provider interface
│   ├── gemini.go               # Gemini implementation
│   ├── openrouter.go           # OpenRouter implementation (future)
│   ├── prompts.go              # Prompt templates & formatting
│   ├── response.go             # Response parsing
│   └── ai_test.go              # AI tests
│
├── git/                         # Git Operations (NEW)
│   ├── diff.go                 # Get git diff
│   ├── commit.go               # Commit operations
│   ├── history.go              # Git history queries
│   └── git_test.go             # Git tests
│
├── ui/                          # Terminal UI (NEW)
│   ├── ui.go                   # UI interface
│   ├── message_selector.go    # Message selection UI
│   ├── message_editor.go       # Message editing UI
│   ├── time_selector.go        # Time selection UI
│   ├── renderer.go             # Termbox rendering utilities
│   └── ui_test.go              # UI tests
│
├── commit/                      # Commit Orchestration (NEW)
│   ├── workflow.go             # Commit workflow orchestration
│   ├── delayed.go              # Delayed commit logic
│   ├── validation.go           # Time validation
│   └── commit_test.go          # Commit tests
│
└── internal/                    # Internal utilities (NEW)
    ├── time/
    │   ├── parser.go           # Time string parsing
    │   ├── validation.go       # Time validation utilities
    │   └── time_test.go        # Time tests
    └── utils/
        ├── validation.go       # General validation
        └── utils_test.go       # Utils tests
```

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         main.go                                 │
│                    (CLI Entry Point)                            │
│  - Flag parsing                                                 │
│  - High-level orchestration                                     │
└────────────┬────────────────────────────────────────────────────┘
             │
             ├──────────────────────────────────────────────┐
             │                                              │
             ▼                                              ▼
┌─────────────────────────┐                   ┌──────────────────────┐
│    commit/workflow.go    │                   │   config/config.go   │
│  Commit Orchestration    │◄──────────────────│  Configuration Mgmt  │
│  - Main workflow         │                   │  - Load/Save config  │
│  - Delayed commit logic  │                   └──────────────────────┘
└────────┬────────────────┘
         │
         ├────────────┬────────────┬────────────┬──────────────┐
         │            │            │            │              │
         ▼            ▼            ▼            ▼              ▼
┌──────────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌──────────────┐
│  git/        │ │  ai/    │ │  ui/    │ │ logger/  │ │ internal/    │
│  Git Ops     │ │ AI Integ│ │ Term UI │ │ Logging  │ │ Utilities    │
│              │ │         │ │         │ │          │ │              │
│ - diff.go    │ │ provider│ │ selector│ │ logger.go│ │ time/        │
│ - commit.go  │ │ gemini  │ │ editor  │ │          │ │ validation/  │
│ - history.go │ │ prompts │ │ time sel│ │          │ │              │
└──────────────┘ └─────────┘ └─────────┘ └──────────┘ └──────────────┘
```

---

## Module Design

### 1. **ai/** - AI Provider Abstraction

**Purpose**: Decouple AI provider implementation from business logic

#### `ai/provider.go` - Provider Interface

```go
package ai

import "context"

// Provider defines the interface for AI commit message generation
type Provider interface {
    // GenerateMessages generates commit messages from a git diff
    GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
    
    // Name returns the provider name
    Name() string
    
    // ValidateConfig validates provider-specific configuration
    ValidateConfig() error
}

// GenerateRequest contains the input for message generation
type GenerateRequest struct {
    GitDiff       string
    LastCommit    string
    Detailed      bool
    IconMode      bool
}

// GenerateResponse contains generated commit messages
type GenerateResponse struct {
    Messages      []string
    RawResponse   string
    PromptUsed    string
}

// Config holds provider configuration
type Config struct {
    Provider  string // "gemini", "openrouter"
    APIKey    string
    Model     string // e.g., "gemini-2.0-flash"
    Timeout   int    // seconds
}
```

#### `ai/gemini.go` - Gemini Implementation

```go
package ai

import (
    "context"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

// GeminiProvider implements Provider for Google Gemini
type GeminiProvider struct {
    apiKey string
    model  string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
    if model == "" {
        model = "gemini-2.0-flash"
    }
    return &GeminiProvider{
        apiKey: apiKey,
        model:  model,
    }
}

// GenerateMessages implements Provider.GenerateMessages
func (g *GeminiProvider) GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
    // Implementation from current generateCommitMessages()
    // Lines 294-522 refactored
}

// Name implements Provider.Name
func (g *GeminiProvider) Name() string {
    return "gemini"
}

// ValidateConfig implements Provider.ValidateConfig
func (g *GeminiProvider) ValidateConfig() error {
    return validateGeminiAPIKey(g.apiKey)
}

// validateGeminiAPIKey validates Gemini API key format
func validateGeminiAPIKey(apiKey string) error {
    // From isValidAPIKey() - lines 1430-1433
}
```

#### `ai/prompts.go` - Prompt Templates

```go
package ai

// PromptType defines the type of prompt to use
type PromptType int

const (
    PromptStandard PromptType = iota
    PromptDetailed
    PromptIcon
    PromptIconDetailed
)

// GetPrompt returns the appropriate prompt template
func GetPrompt(promptType PromptType) string {
    // Returns templates from lines 27-283
}

// FormatPrompt fills in the prompt template
func FormatPrompt(template, gitDiff, lastCommit string) string {
    return fmt.Sprintf(template, gitDiff, lastCommit)
}
```

#### `ai/response.go` - Response Parsing

```go
package ai

// ResponseParser handles parsing of AI responses
type ResponseParser struct {
    detailed bool
    iconMode bool
}

// ParseResponse parses AI response into commit messages
func (p *ResponseParser) ParseResponse(text string) ([]string, error) {
    // Extract parsing logic from generateCommitMessages()
    // Lines 344-522
}

// parseJSON parses JSON format responses
func parseJSON(text string) ([]string, error) {
    // Lines 414-512
}

// parseDetailed parses detailed format responses
func parseDetailed(text string) ([]string, error) {
    // Lines 351-408
}
```

**Responsibilities**:
- ✅ Abstract AI provider implementations
- ✅ Manage prompt templates and formatting
- ✅ Parse AI responses into structured data
- ✅ Validate provider configurations
- ✅ Handle provider-specific errors

**Public Interface**:
- `Provider` interface for implementations
- `NewGeminiProvider()` factory
- `NewOpenRouterProvider()` factory (future)

---

### 2. **git/** - Git Operations

**Purpose**: Encapsulate all Git interactions

#### `git/diff.go`

```go
package git

import (
    "os/exec"
    "fmt"
)

// GetStagedDiff returns the staged changes
func GetStagedDiff() (string, error) {
    // From getGitDiff() - lines 285-292
    cmd := exec.Command("git", "diff", "--cached")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("error getting git diff: %v", err)
    }
    return string(output), nil
}

// HasStagedChanges checks if there are staged changes
func HasStagedChanges() (bool, error) {
    diff, err := GetStagedDiff()
    if err != nil {
        return false, err
    }
    return diff != "", nil
}
```

#### `git/history.go`

```go
package git

// GetLastCommitMessage returns the last commit message
func GetLastCommitMessage() (string, error) {
    // From getLastCommitMessage() - lines 524-531
    cmd := exec.Command("git", "log", "-1", "--pretty=%B")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("error getting last commit message: %v", err)
    }
    return strings.TrimSpace(string(output)), nil
}
```

#### `git/commit.go`

```go
package git

import (
    "os/exec"
    "time"
)

// CommitOptions contains options for creating a commit
type CommitOptions struct {
    Message   string
    Timestamp *time.Time // nil for current time
}

// CreateCommit creates a git commit
func CreateCommit(opts CommitOptions) error {
    if opts.Timestamp != nil {
        return createCommitWithDate(opts.Message, *opts.Timestamp)
    }
    return createCommitNow(opts.Message)
}

// createCommitNow creates a commit with current timestamp
func createCommitNow(message string) error {
    cmd := exec.Command("git", "commit", "-m", message)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// createCommitWithDate creates a commit with custom timestamp
func createCommitWithDate(message string, timestamp time.Time) error {
    // From executeDelayedCommit() - lines 1389-1428
    dateStr := timestamp.Format(time.RFC3339)
    cmd := exec.Command("git", "commit", "-m", message, "--date", dateStr)
    cmd.Env = append(os.Environ(), 
        fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
```

**Responsibilities**:
- ✅ Execute git commands
- ✅ Get staged changes
- ✅ Query git history
- ✅ Create commits (normal and with custom timestamp)
- ✅ Handle git command errors

**Public Interface**:
- `GetStagedDiff() (string, error)`
- `GetLastCommitMessage() (string, error)`
- `CreateCommit(CommitOptions) error`

---

### 3. **ui/** - Terminal UI

**Purpose**: Handle all terminal interactions using termbox

#### `ui/ui.go` - UI Interface

```go
package ui

import "github.com/nsf/termbox-go"

// UI provides terminal user interface operations
type UI interface {
    // Initialize initializes the UI
    Initialize() error
    
    // Close closes the UI
    Close() error
    
    // SelectMessage shows message selection UI
    SelectMessage(messages []string) (string, error)
    
    // EditMessage shows message editing UI
    EditMessage(initial string) (string, error)
    
    // SelectTime shows time selection UI
    SelectTime(options []string) (string, bool, error)
}

// TermboxUI implements UI using termbox-go
type TermboxUI struct {
    initialized bool
}

// NewTermboxUI creates a new termbox UI
func NewTermboxUI() *TermboxUI {
    return &TermboxUI{}
}
```

#### `ui/message_selector.go`

```go
package ui

// MessageSelector handles message selection UI
type MessageSelector struct {
    messages []string
    selected int
}

// Show displays the message selection UI
func (m *MessageSelector) Show() (string, error) {
    // From getUserChoice() - lines 853-920
    // From drawMessages() - lines 533-615
}
```

#### `ui/message_editor.go`

```go
package ui

// MessageEditor handles message editing UI
type MessageEditor struct {
    initialMessage string
}

// Edit shows the message editing UI
func (e *MessageEditor) Edit() (string, error) {
    // From editMessage() - lines 617-851
}
```

#### `ui/time_selector.go`

```go
package ui

// TimeSelector handles time selection UI
type TimeSelector struct {
    options []string
}

// Show displays time selection and returns selected time
// Returns: (selectedTime string, isCustom bool, error)
func (t *TimeSelector) Show() (string, bool, error) {
    // From showTimeSelectionUI() - lines 1120-1242
    // From handleCustomTimeInput() - lines 1244-1387
}
```

#### `ui/renderer.go` - Rendering Utilities

```go
package ui

import "github.com/nsf/termbox-go"

// RenderText renders text at specified position
func RenderText(x, y int, text string, fg, bg termbox.Attribute) {
    for i, ch := range text {
        termbox.SetCell(x+i, y, ch, fg, bg)
    }
}

// RenderBullet renders a bullet point
func RenderBullet(x, y int, selected bool, fg, bg termbox.Attribute) {
    bullet := "•"
    if selected {
        bullet = "→"
    }
    termbox.SetCell(x, y, []rune(bullet)[0], fg, bg)
}

// GetTerminalSize returns terminal dimensions
func GetTerminalSize() (width, height int) {
    w, h := termbox.Size()
    if w < 10 {
        w = 80
    }
    if h < 10 {
        h = 24
    }
    return w, h
}
```

**Responsibilities**:
- ✅ Display message selection interface
- ✅ Handle message editing
- ✅ Show time selection for delayed commits
- ✅ Manage keyboard input
- ✅ Render UI elements with termbox

**Public Interface**:
- `UI` interface for testing
- `NewTermboxUI() *TermboxUI`
- Message, editor, and time selectors

---

### 4. **commit/** - Commit Orchestration

**Purpose**: Orchestrate the commit workflow

#### `commit/workflow.go`

```go
package commit

import (
    "context"
    "github.com/thanhphuchuynh/ai"
    "github.com/thanhphuchuynh/git"
    "github.com/thanhphuchuynh/ui"
    "github.com/thanhphuchuynh/config"
)

// Workflow orchestrates the commit process
type Workflow struct {
    aiProvider ai.Provider
    ui         ui.UI
    config     *config.Config
}

// NewWorkflow creates a new commit workflow
func NewWorkflow(provider ai.Provider, ui ui.UI, cfg *config.Config) *Workflow {
    return &Workflow{
        aiProvider: provider,
        ui:         ui,
        config:     cfg,
    }
}

// Execute runs the complete commit workflow
func (w *Workflow) Execute(ctx context.Context) error {
    // 1. Get git diff
    diff, err := git.GetStagedDiff()
    if err != nil {
        return err
    }
    
    if diff == "" {
        return errors.New("no staged changes found")
    }
    
    // 2. Get last commit for context
    lastCommit, _ := git.GetLastCommitMessage()
    
    // 3. Generate commit messages
    messages, err := w.generateMessages(ctx, diff, lastCommit)
    if err != nil {
        return err
    }
    
    // 4. User selects/edits message
    selectedMsg, err := w.ui.SelectMessage(messages)
    if err != nil {
        return err
    }
    
    // 5. Handle delayed commit if enabled
    timestamp, err := w.handleDelayedCommit()
    if err != nil {
        return err
    }
    
    // 6. Create commit
    opts := git.CommitOptions{
        Message:   selectedMsg,
        Timestamp: timestamp,
    }
    
    return git.CreateCommit(opts)
}

// generateMessages generates commit messages using AI
func (w *Workflow) generateMessages(ctx context.Context, diff, lastCommit string) ([]string, error) {
    req := ai.GenerateRequest{
        GitDiff:    diff,
        LastCommit: lastCommit,
        Detailed:   w.config.DetailedMode,
        IconMode:   w.config.IconMode,
    }
    
    resp, err := w.aiProvider.GenerateMessages(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return resp.Messages, nil
}
```

#### `commit/delayed.go`

```go
package commit

import (
    "time"
    "github.com/thanhphuchuynh/config"
    "github.com/thanhphuchuynh/ui"
)

// DelayedCommitHandler handles delayed commit logic
type DelayedCommitHandler struct {
    config *config.DelayedCommitConfig
    ui     ui.UI
}

// HandleDelayedCommit checks if delayed commit is needed and handles it
// Returns timestamp to use (nil for current time)
func (h *DelayedCommitHandler) HandleDelayedCommit() (*time.Time, error) {
    if !h.config.Enabled {
        return nil, nil
    }
    
    now := time.Now()
    if !isTimeInRestrictedRange(now, h.config) {
        return nil, nil
    }
    
    // Show time selection UI
    options := generateSuggestedTimes(h.config.RestrictedEndHour, h.config.SuggestionInterval)
    selectedTime, isCustom, err := h.ui.SelectTime(options)
    if err != nil {
        return nil, err
    }
    
    // Parse selected time
    timestamp, err := parseTimeSelection(selectedTime, now)
    if err != nil {
        return nil, err
    }
    
    return &timestamp, nil
}

// isTimeInRestrictedRange checks if current time is restricted
func isTimeInRestrictedRange(now time.Time, cfg *config.DelayedCommitConfig) bool {
    // From lines 926-944
}

// generateSuggestedTimes generates time suggestions
func generateSuggestedTimes(endHour, interval int) []string {
    // From lines 1060-1101
}
```

#### `commit/validation.go`

```go
package commit

// ValidateCommitMessage validates a commit message
func ValidateCommitMessage(msg string) error {
    if msg == "" {
        return errors.New("commit message cannot be empty")
    }
    
    if len(msg) > 1000 {
        return errors.New("commit message too long")
    }
    
    return nil
}
```

**Responsibilities**:
- ✅ Orchestrate commit workflow
- ✅ Coordinate between ai, git, ui modules
- ✅ Handle delayed commit logic
- ✅ Validate commit messages
- ✅ Error handling and recovery

**Public Interface**:
- `NewWorkflow() *Workflow`
- `Execute(context.Context) error`

---

### 5. **internal/** - Internal Utilities

**Purpose**: Shared utilities not exposed publicly

#### `internal/time/parser.go`

```go
package time

import (
    "strconv"
    "strings"
    "time"
)

// ParseTimeString parses various time formats
func ParseTimeString(timeStr string) (hour, minute int, err error) {
    // From parseTimeString() - lines 946-1050
    // Supports: "HH:MM", "HH:MM AM/PM", "HHhMM"
}
```

#### `internal/time/validation.go`

```go
package time

// ValidateHour validates hour is in range 0-23
func ValidateHour(hour int) error {
    // From validateHour() - lines 1052-1058
}

// ValidateMinute validates minute is in range 0-59
func ValidateMinute(minute int) error {
    if minute < 0 || minute > 59 {
        return fmt.Errorf("minute must be between 0 and 59")
    }
    return nil
}
```

**Responsibilities**:
- ✅ Time parsing utilities
- ✅ Time validation
- ✅ General utility functions

---

## AI Provider Abstraction

### Interface Design

The AI provider abstraction enables supporting multiple AI services (Gemini, OpenRouter, OpenAI, etc.) through a common interface.

```go
// Provider interface (primary abstraction)
type Provider interface {
    GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
    Name() string
    ValidateConfig() error
}
```

### Provider Factory

```go
// NewProvider creates a provider based on configuration
func NewProvider(cfg Config) (Provider, error) {
    switch strings.ToLower(cfg.Provider) {
    case "gemini", "":
        return NewGeminiProvider(cfg.APIKey, cfg.Model), nil
    case "openrouter":
        return NewOpenRouterProvider(cfg.APIKey, cfg.Model), nil
    default:
        return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
    }
}
```

### OpenRouter Implementation (Future)

```go
package ai

// OpenRouterProvider implements Provider for OpenRouter
type OpenRouterProvider struct {
    apiKey string
    model  string
    baseURL string
}

func NewOpenRouterProvider(apiKey, model string) *OpenRouterProvider {
    return &OpenRouterProvider{
        apiKey:  apiKey,
        model:   model,
        baseURL: "https://openrouter.ai/api/v1",
    }
}

func (o *OpenRouterProvider) GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
    // Implementation using OpenRouter API
}
```

### Configuration Extension

Add to [`config/config.go`](config/config.go:21):

```go
type Config struct {
    APIKey         string              `json:"api_key"`
    LoggingEnabled bool                `json:"logging_enabled"`
    IconMode       bool                `json:"icon_mode"`
    DelayedCommit  DelayedCommitConfig `json:"delayed_commit"`
    
    // New fields for AI provider
    AIProvider     string              `json:"ai_provider"`      // "gemini", "openrouter"
    AIModel        string              `json:"ai_model"`         // Model name
    DetailedMode   bool                `json:"detailed_mode"`    // Generate detailed messages
}
```

---

## Function Migration Plan

### Phase-by-Phase Migration

| Current Location | New Location | Lines | Complexity | Phase |
|-----------------|--------------|-------|------------|-------|
| **Constants & Templates** |
| Prompt templates | `ai/prompts.go` | 27-283 | Low | 4 |
| **Data Structures** |
| `CommitResponse` | `ai/response.go` | 22-26 | Low | 4 |
| **Git Operations** |
| `getGitDiff()` | `git/diff.go` | 285-292 | Low | 3 |
| `getLastCommitMessage()` | `git/history.go` | 524-531 | Low | 3 |
| `executeDelayedCommit()` | `git/commit.go` | 1389-1428 | Medium | 3 |
| **AI Integration** |
| `generateCommitMessages()` | `ai/gemini.go` | 294-522 | High | 4 |
| JSON parsing logic | `ai/response.go` | 344-522 | High | 4 |
| `isValidAPIKey()` | `ai/gemini.go` | 1430-1433 | Low | 4 |
| **UI Functions** |
| `drawMessages()` | `ui/message_selector.go` | 533-615 | Medium | 2 |
| `editMessage()` | `ui/message_editor.go` | 617-851 | High | 2 |
| `getUserChoice()` | `ui/message_selector.go` | 853-920 | Medium | 2 |
| `showTimeSelectionUI()` | `ui/time_selector.go` | 1120-1242 | High | 2 |
| `handleCustomTimeInput()` | `ui/time_selector.go` | 1244-1387 | High | 2 |
| **Time/Delayed Commit** |
| `isTimeInRestrictedRange()` | `commit/delayed.go` | 926-944 | Low | 1 |
| `parseTimeString()` | `internal/time/parser.go` | 946-1050 | Medium | 1 |
| `validateHour()` | `internal/time/validation.go` | 1052-1058 | Low | 1 |
| `generateSuggestedTimes()` | `commit/delayed.go` | 1060-1101 | Low | 1 |
| `getNextAvailableTime()` | `commit/delayed.go` | 1103-1117 | Low | 1 |
| **Main Orchestration** |
| `main()` flag handling | `main.go` | 1435-1544 | Medium | 5 |
| `main()` workflow | `commit/workflow.go` | 1546-1649 | High | 5 |

### Migration Priority

**High Priority** (Complex, high-value refactoring):
1. `generateCommitMessages()` → `ai/gemini.go`
2. `editMessage()` → `ui/message_editor.go`
3. `main()` workflow → `commit/workflow.go`

**Medium Priority** (Moderate complexity):
4. UI functions → `ui/` package
5. Time parsing → `internal/time/`

**Low Priority** (Simple, low-risk):
6. Git operations → `git/` package
7. Utility functions → `internal/utils/`

---

## Dependency Graph

### Module Dependencies

```
main.go
    ├─→ commit/workflow
    │       ├─→ ai/provider
    │       │       ├─→ ai/prompts
    │       │       └─→ ai/response
    │       ├─→ git/
    │       │       ├─→ git/diff
    │       │       ├─→ git/commit
    │       │       └─→ git/history
    │       ├─→ ui/
    │       │       ├─→ ui/message_selector
    │       │       ├─→ ui/message_editor
    │       │       ├─→ ui/time_selector
    │       │       └─→ ui/renderer
    │       ├─→ commit/delayed
    │       │       └─→ internal/time/
    │       ├─→ config/
    │       └─→ logger/
    │
    └─→ config/
```

### Dependency Rules

1. **No Circular Dependencies**: Enforce strict DAG
2. **internal/ Cannot Be Imported**: Only by packages in same repo
3. **UI Abstraction**: UI interfaces, not concrete termbox calls
4. **Provider Abstraction**: AI providers behind interface
5. **Config Independence**: Config package has no dependencies except stdlib

### Import Constraints

```go
// ✅ Allowed
import "github.com/thanhphuchuynh/commit"
import "github.com/thanhphuchuynh/ai"

// ❌ Not Allowed - circular
// In ai/: import "github.com/thanhphuchuynh/commit"

// ✅ Allowed - internal within module
import "github.com/thanhphuchuynh/internal/time"

// ❌ Not Allowed - internal from outside
// External package: import "github.com/thanhphuchuynh/internal/time"
```

---

## Refactoring Phases

### Phase 1: Extract Time Utilities (Low Risk)

**Goal**: Extract time-related functions to reduce main.go complexity

**Tasks**:
1. Create `internal/time/` package
2. Move time parsing functions
3. Move time validation functions
4. Write comprehensive tests
5. Update main.go imports

**Files Modified**:
- NEW: `internal/time/parser.go`
- NEW: `internal/time/validation.go`
- NEW: `internal/time/time_test.go`
- MODIFY: `main.go` (remove ~150 lines)

**Testing**:
```bash
# Run tests
go test ./internal/time/...

# Verify main.go still works
go build -o gocommit
./gocommit --help
```

**Success Criteria**:
- ✅ All time tests pass
- ✅ main.go compiles without errors
- ✅ Existing functionality unchanged
- ✅ ~150 lines removed from main.go

---

### Phase 2: Extract UI Components (Medium Risk)

**Goal**: Separate UI logic into dedicated package

**Tasks**:
1. Create `ui/` package structure
2. Move message selection UI
3. Move message editor UI
4. Move time selection UI
5. Create UI interface for testing
6. Extract rendering utilities
7. Write UI tests (may require mocking termbox)

**Files Modified**:
- NEW: `ui/ui.go`
- NEW: `ui/message_selector.go`
- NEW: `ui/message_editor.go`
- NEW: `ui/time_selector.go`
- NEW: `ui/renderer.go`
- NEW: `ui/ui_test.go`
- MODIFY: `main.go` (remove ~450 lines)

**Testing**:
```bash
go test ./ui/...
go build -o gocommit
# Manual testing of UI
./gocommit
```

**Success Criteria**:
- ✅ UI tests pass
- ✅ Message selection works
- ✅ Message editing works
- ✅ Time selection works
- ✅ ~450 lines removed from main.go

---

### Phase 3: Extract Git Operations (Low Risk)

**Goal**: Encapsulate Git interactions

**Tasks**:
1. Create `git/` package
2. Move diff operations
3. Move commit operations
4. Move history operations
5. Add error handling
6. Write Git operation tests

**Files Modified**:
- NEW: `git/diff.go`
- NEW: `git/commit.go`
- NEW: `git/history.go`
- NEW: `git/git_test.go`
- MODIFY: `main.go` (remove ~50 lines)
- MODIFY: `commit/` (will use git package)

**Testing**:
```bash
go test ./git/...
# Integration test
./gocommit
```

**Success Criteria**:
- ✅ Git tests pass
- ✅ Can get staged diff
- ✅ Can create commits
- ✅ Can query history
- ✅ ~50 lines removed from main.go

---

### Phase 4: Extract AI Provider Logic (High Risk)

**Goal**: Abstract AI provider and enable multi-provider support

**Tasks**:
1. Create `ai/` package
2. Define Provider interface
3. Extract Gemini implementation
4. Extract prompt templates
5. Extract response parsing
6. Add provider factory
7. Write comprehensive AI tests

**Files Modified**:
- NEW: `ai/provider.go`
- NEW: `ai/gemini.go`
- NEW: `ai/prompts.go`
- NEW: `ai/response.go`
- NEW: `ai/ai_test.go`
- MODIFY: `main.go` (remove ~500 lines)
- MODIFY: `config/config.go` (add AI provider config)

**Testing**:
```bash
go test ./ai/...
# Integration test with real API
export GEMINI_API_KEY="your-key"
./gocommit
```

**Success Criteria**:
- ✅ AI tests pass
- ✅ Gemini provider works
- ✅ Response parsing works
- ✅ Provider abstraction complete
- ✅ ~500 lines removed from main.go

---

### Phase 5: Extract Commit Orchestration (High Risk)

**Goal**: Move workflow logic out of main.go

**Tasks**:
1. Create `commit/` package
2. Extract workflow orchestration
3. Extract delayed commit logic
4. Add validation functions
5. Write workflow tests
6. Integration testing

**Files Modified**:
- NEW: `commit/workflow.go`
- NEW: `commit/delayed.go`
- NEW: `commit/validation.go`
- NEW: `commit/commit_test.go`
- MODIFY: `main.go` (remove ~300 lines)

**Testing**:
```bash
go test ./commit/...
# Full integration test
./gocommit
./gocommit --enable-delayed
./gocommit --config-delayed
```

**Success Criteria**:
- ✅ Workflow tests pass
- ✅ Delayed commit works
- ✅ Normal commit works
- ✅ Error handling works
- ✅ ~300 lines removed from main.go

---

### Phase 6: Slim Down main.go (Low Risk)

**Goal**: Reduce main.go to minimal entry point

**Tasks**:
1. Keep only CLI flag parsing
2. Keep only high-level orchestration
3. Move flag handlers to appropriate packages
4. Clean up imports
5. Add main.go documentation
6. Final integration testing

**Files Modified**:
- MODIFY: `main.go` (reduce to 80-100 lines)
- MODIFY: `config/config.go` (add flag handler helpers)

**Final main.go Structure**:
```go
package main

import (
    "context"
    "flag"
    "log"
    
    "github.com/thanhphuchuynh/ai"
    "github.com/thanhphuchuynh/commit"
    "github.com/thanhphuchuynh/config"
    "github.com/thanhphuchuynh/ui"
)

func main() {
    // Parse flags
    cfg := parseFlags()
    
    // Handle config commands
    if handleConfigCommands(cfg) {
        return
    }
    
    // Setup components
    provider := ai.NewProvider(cfg.AI)
    uiManager := ui.NewTermboxUI()
    workflow := commit.NewWorkflow(provider, uiManager, cfg)
    
    // Execute commit workflow
    ctx := context.Background()
    if err := workflow.Execute(ctx); err != nil {
        log.Fatal(err)
    }
}

func parseFlags() *config.Config {
    // Flag parsing
}

func handleConfigCommands(cfg *config.Config) bool {
    // Handle --config, --enable-logging, etc.
}
```

**Testing**:
```bash
# Full regression testing
go test ./...
./gocommit --help
./gocommit --config
./gocommit
./gocommit -d
./gocommit --icon
./gocommit --enable-delayed
```

**Success Criteria**:
- ✅ main.go < 100 lines
- ✅ All tests pass
- ✅ All features work
- ✅ Clean architecture
- ✅ Ready for future enhancements

---

## Testing Strategy

### Unit Testing

**Test Coverage Goals**:
- **ai/**: 80%+ coverage
- **git/**: 70%+ coverage (mocking git commands)
- **ui/**: 60%+ coverage (manual testing supplement)
- **commit/**: 80%+ coverage
- **internal/**: 90%+ coverage

**Test Structure**:
```
pkg/
├── pkg.go
└── pkg_test.go
```

### Integration Testing

**Test Scenarios**:
1. **Normal commit flow**: Stage → Generate → Select → Commit
2. **Detailed mode**: Test detailed message generation
3. **Icon mode**: Test icon message generation
4. **Delayed commit**: Test time restriction and selection
5. **Custom time**: Test custom time input
6. **Edit message**: Test message editing
7. **Error cases**: Test various error scenarios

### Regression Testing Checklist

Before each phase completion:

```bash
# Build
go build -o gocommit

# Basic functionality
./gocommit --help
./gocommit --version

# Configuration
./gocommit --config
./gocommit --enable-logging
./gocommit --disable-logging

# Message generation (requires staged changes)
git add test.txt
./gocommit              # Normal mode
./gocommit -d           # Detailed mode
./gocommit --icon       # Icon mode

# Delayed commit
./gocommit --config-delayed
./gocommit --enable-delayed
./gocommit              # During restricted hours

# Edge cases
./gocommit              # No staged changes (should error)
./gocommit --config     # Invalid API key
```

### Automated Testing

Add to CI/CD pipeline:

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: go vet ./...
      - run: staticcheck ./...
```

---

## Backward Compatibility

### CLI Compatibility

**Guaranteed Compatibility**:
- ✅ All existing flags work unchanged
- ✅ Configuration file format compatible
- ✅ Same user experience
- ✅ No breaking changes to commands

**Flags to Preserve**:
```bash
--config              # Configure API key
--d                   # Detailed mode
--icon                # Icon mode
--enable-logging      # Enable logging
--disable-logging     # Disable logging
--config-delayed      # Configure delayed commit
--enable-delayed      # Enable delayed commit
--disable-delayed     # Disable delayed commit
```

### Configuration Compatibility

**Config File Migration**:

Old format (current):
```json
{
  "api_key": "AIza...",
  "logging_enabled": true,
  "icon_mode": false,
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval": 20
  }
}
```

New format (backward compatible):
```json
{
  "api_key": "AIza...",
  "logging_enabled": true,
  "icon_mode": false,
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval": 20
  },
  "ai_provider": "gemini",
  "ai_model": "gemini-2.0-flash",
  "detailed_mode": false
}
```

**Migration Strategy**:
```go
// In config.LoadConfig()
func LoadConfig() (*Config, error) {
    // Load existing config
    config, err := loadConfigFile()
    if err != nil {
        return defaultConfig(), nil
    }
    
    // Apply defaults for new fields
    if config.AIProvider == "" {
        config.AIProvider = "gemini"
    }
    if config.AIModel == "" {
        config.AIModel = "gemini-2.0-flash"
    }
    
    return config, nil
}
```

### API Compatibility

**Internal APIs**:
- New packages are additions, not replacements
- Existing config package enhanced, not replaced
- Logger package unchanged

**No Breaking Changes**:
- Users won't need to reconfigure
- Existing workflows continue working
- Gradual migration path for custom integrations

---

## Implementation Guidelines

### Code Style

**Follow Go Best Practices**:
```go
// ✅ Good - Clear function names
func GenerateCommitMessages(ctx context.Context, diff string) ([]string, error)

// ❌ Bad - Unclear abbreviations
func GenMsg(ctx context.Context, d string) ([]string, error)

// ✅ Good - Error wrapping
return fmt.Errorf("failed to parse response: %w", err)

// ❌ Bad - Lost context
return err

// ✅ Good - Interface naming
type Provider interface { ... }

// ❌ Bad - Redundant naming
type ProviderInterface interface { ... }
```

### Functional Programming Patterns

**Pure Functions** (where possible):
```go
// ✅ Pure function - no side effects, deterministic
func ParseTimeString(timeStr string) (hour, minute int, err error) {
    // Only depends on input, no external state
    // Same input always produces same output
}

// ✅ Pure function - immutable operations
func FormatPrompt(template, diff, commit string) string {
    return fmt.Sprintf(template, diff, commit)
}
```

**Higher-Order Functions**:
```go
// ✅ Function that returns a function
func ValidatorChain(validators ...func(string) error) func(string) error {
    return func(value string) error {
        for _, validator := range validators {
            if err := validator(value); err != nil {
                return err
            }
        }
        return nil
    }
}

// Usage
validator := ValidatorChain(
    validateNotEmpty,
    validateLength,
    validateFormat,
)
```

**Immutable Data Structures**:
```go
// ✅ Return new struct instead of modifying
func (c Config) WithDetailedMode(enabled bool) Config {
    newConfig := c
    newConfig.DetailedMode = enabled
    return newConfig
}

// ❌ Avoid mutable operations
func (c *Config) SetDetailedMode(enabled bool) {
    c.DetailedMode = enabled
}
```

**Function Composition**:
```go
// ✅ Compose small functions
func ProcessCommitMessage(msg string) string {
    return trim(lowercase(removeSpecialChars(msg)))
}
```

### Error Handling

**Consistent Error Patterns**:
```go
// ✅ Wrap errors with context
if err != nil {
    return fmt.Errorf("generating commit messages: %w", err)
}

// ✅ Custom error types for specific cases
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ✅ Error handling with recovery
defer func() {
    if r := recover(); r != nil {
        log.Printf("Recovered from panic: %v", r)
        // Fallback logic
    }
}()
```

### Package Organization

**Package Naming**:
```go
// ✅ Good - Short, clear
package ai
package git
package ui

// ❌ Bad - Verbose
package aiintegration
package gitoperations
```

**File Naming**:
```go
// ✅ Good - Descriptive
provider.go          // Interfaces
gemini.go           // Implementation
prompts.go          // Templates
response.go         // Parsing

// ❌ Bad - Vague
types.go
utils.go
helpers.go
```

### Documentation

**Package Documentation**:
```go
// Package ai provides AI-powered commit message generation.
//
// It supports multiple AI providers through a common interface,
// allowing easy integration of new providers.
//
// Example usage:
//
//	provider := ai.NewGeminiProvider(apiKey, model)
//	req := ai.GenerateRequest{
//		GitDiff: diff,
//		Detailed: true,
//	}
//	resp, err := provider.GenerateMessages(ctx, req)
package ai
```

**Function Documentation**:
```go
// GenerateMessages generates commit messages from a git diff.
//
// It takes a GenerateRequest containing the git diff and options,
// and returns a GenerateResponse with the generated messages.
//
// The function may return an error if:
//   - The API request fails
//   - The response cannot be parsed
//   - The API key is invalid
//
// Example:
//
//	req := GenerateRequest{GitDiff: diff}
//	resp, err := provider.GenerateMessages(ctx, req)
//	if err != nil {
//		return fmt.Errorf("generation failed: %w", err)
//	}
func (p *Provider) GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
```

### Testing Best Practices

**Table-Driven Tests**:
```go
func TestParseTimeString(t *testing.T) {
    tests := []struct {
        name       string
        input      string
        wantHour   int
        wantMinute int
        wantErr    bool
    }{
        {"24-hour format", "18:30", 18, 30, false},
        {"12-hour AM", "9:30 AM", 9, 30, false},
        {"12-hour PM", "6:30 PM", 18, 30, false},
        {"h format", "18h30", 18, 30, false},
        {"invalid format", "invalid", 0, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hour, minute, err := ParseTimeString(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseTimeString() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if hour != tt.wantHour || minute != tt.wantMinute {
                t.Errorf("ParseTimeString() = (%v, %v), want (%v, %v)", 
                    hour, minute, tt.wantHour, tt.wantMinute)
            }
        })
    }
}
```

**Mocking External Dependencies**:
```go
// Mock AI provider for testing
type MockProvider struct {
    GenerateFunc func(context.Context, GenerateRequest) (*GenerateResponse, error)
}

func (m *MockProvider) GenerateMessages(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
    if m.GenerateFunc != nil {
        return m.GenerateFunc(ctx, req)
    }
    return &GenerateResponse{Messages: []string{"mock message"}}, nil
}
```

---

## Success Metrics

### Quantitative Metrics

**Code Metrics**:
- ✅ main.go reduced from 1649 → <100 lines (94% reduction)
- ✅ Average function length: <50 lines
- ✅ Maximum function complexity: <15 (cyclomatic)
- ✅ Test coverage: >75% overall
- ✅ Number of packages: 7 (well-organized)

**Quality Metrics**:
- ✅ No circular dependencies
- ✅ All tests passing
- ✅ No linter warnings
- ✅ Documentation coverage: 100% of public APIs

### Qualitative Metrics

**Maintainability**:
- ✅ Easy to add new AI providers
- ✅ Easy to modify UI without affecting logic
- ✅ Clear separation of concerns
- ✅ Self-documenting code structure

**Developer Experience**:
- ✅ Clear package boundaries
- ✅ Intuitive naming conventions
- ✅ Comprehensive documentation
- ✅ Easy to test

**User Experience**:
- ✅ No breaking changes
- ✅ Same or better performance
- ✅ All features working
- ✅ Improved error messages

---

## Risk Assessment

### High-Risk Areas

1. **AI Response Parsing** (Phase 4)
   - **Risk**: Breaking existing response parsing
   - **Mitigation**: Extensive testing, keep fallbacks
   - **Rollback**: Can revert to original parsing

2. **UI Refactoring** (Phase 2)
   - **Risk**: Breaking termbox integration
   - **Mitigation**: Careful extraction, manual testing
   - **Rollback**: UI is self-contained, easy to revert

3. **Main Workflow** (Phase 5)
   - **Risk**: Breaking commit orchestration
   - **Mitigation**: Comprehensive integration tests
   - **Rollback**: Keep original main.go as backup

### Medium-Risk Areas

4. **Git Operations** (Phase 3)
   - **Risk**: Git command failures
   - **Mitigation**: Test on multiple platforms
   - **Rollback**: Simple functions, easy to fix

5. **Time Logic** (Phase 1)
   - **Risk**: Time parsing edge cases
   - **Mitigation**: Extensive unit tests
   - **Rollback**: Low impact, easy to fix

### Low-Risk Areas

6. **Configuration** (All phases)
   - **Risk**: Config compatibility
   - **Mitigation**: Backward-compatible design
   - **Rollback**: Config is versioned

---

## Timeline Estimate

### Conservative Estimate

| Phase | Tasks | Complexity | Estimated Time | Risk Level |
|-------|-------|------------|----------------|------------|
| Phase 1: Time Utils | 5 tasks | Low | 2-3 hours | Low |
| Phase 2: UI | 7 tasks | Medium | 4-6 hours | Medium |
| Phase 3: Git Ops | 6 tasks | Low | 2-4 hours | Low |
| Phase 4: AI Provider | 7 tasks | High | 6-8 hours | High |
| Phase 5: Workflow | 6 tasks | High | 5-7 hours | High |
| Phase 6: Main Slim | 6 tasks | Low | 2-3 hours | Low |
| **Testing & Polish** | Various | Medium | 4-6 hours | Medium |
| **Total** | 37 tasks | - | **25-37 hours** | - |

### Aggressive Estimate

With focused work: **15-20 hours** total

### Phased Rollout

- **Week 1**: Phases 1-2 (Time + UI)
- **Week 2**: Phases 3-4 (Git + AI)
- **Week 3**: Phases 5-6 (Workflow + Main)
- **Week 4**: Testing, documentation, polish

---

## Next Steps

### Immediate Actions

1. **Review & Approval**: Get stakeholder approval for this design
2. **Branch Strategy**: Create feature branch `refactor/modular-architecture`
3. **Backup**: Tag current version as `v1.x-pre-refactor`
4. **Setup**: Create initial package directories

### Phase 1 Kickoff

Once approved, begin Phase 1:

```bash
# Create branch
git checkout -b refactor/phase-1-time-utils

# Create package structure
mkdir -p internal/time
touch internal/time/parser.go
touch internal/time/validation.go
touch internal/time/time_test.go

# Begin implementation
# ... (follow Phase 1 tasks)
```

### Communication Plan

- **Progress Updates**: After each phase completion
- **Code Reviews**: Before merging each phase
- **Demo**: After Phase 6 for full integration
- **Documentation**: Update README with new architecture

---

## Appendix

### A. Detailed File Sizes

| Current File | Lines | Future Files | Estimated Lines |
|-------------|-------|--------------|-----------------|
| main.go | 1649 | main.go | 80-100 |
| | | ai/provider.go | 50 |
| | | ai/gemini.go | 200 |
| | | ai/prompts.go | 100 |
| | | ai/response.go | 150 |
| | | git/diff.go | 30 |
| | | git/commit.go | 80 |
| | | git/history.go | 30 |
| | | ui/message_selector.go | 120 |
| | | ui/message_editor.go | 250 |
| | | ui/time_selector.go | 200 |
| | | ui/renderer.go | 50 |
| | | commit/workflow.go | 150 |
| | | commit/delayed.go | 120 |
| | | internal/time/parser.go | 100 |
| | | internal/time/validation.go | 30 |
| **Total** | **1649** | **Total** | **~1740** |

Note: Total lines increase slightly due to package declarations, imports, and documentation, but complexity per file decreases dramatically.

### B. Go Module Structure

```
module github.com/thanhphuchuynh

go 1.21

require (
    github.com/google/generative-ai-go v0.5.0
    github.com/nsf/termbox-go v1.1.1
    google.golang.org/api v0.155.0
)
```

No additional dependencies needed for refactoring.

### C. OpenRouter Provider Specification

For future implementation:

```go
// OpenRouter API endpoint
const openRouterAPI = "https://openrouter.ai/api/v1/chat/completions"

// Request format
type OpenRouterRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Response format
type OpenRouterResponse struct {
    Choices []Choice `json:"choices"`
}

type Choice struct {
    Message Message `json:"message"`
}
```

---

## Summary

This refactoring design provides:

✅ **Clear Architecture**: 7 well-defined packages with single responsibilities  
✅ **AI Abstraction**: Support for multiple AI providers (Gemini, OpenRouter, future)  
✅ **Clean Code**: Small, focused functions following SOLID principles  
✅ **Functional Patterns**: Pure functions, immutability, composition where appropriate  
✅ **Testability**: Easy to unit test with clear interfaces  
✅ **Maintainability**: Easy to understand, modify, and extend  
✅ **Backward Compatibility**: No breaking changes to user experience  
✅ **Phased Approach**: 6 phases with clear milestones and rollback points  

**Main Complexity Reduction**: 1649 lines → ~100 lines in main.go (94% reduction)

The design is comprehensive, well-documented, and ready for implementation approval.