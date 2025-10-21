# Delayed Commit Feature - Architecture Design Document

## Table of Contents
1. [Overview](#overview)
2. [Requirements Summary](#requirements-summary)
3. [Architecture Design](#architecture-design)
4. [Configuration Schema](#configuration-schema)
5. [Data Structures](#data-structures)
6. [Function Signatures](#function-signatures)
7. [Workflow Diagrams](#workflow-diagrams)
8. [Integration Points](#integration-points)
9. [User Interface Flow](#user-interface-flow)
10. [Git Command Modification](#git-command-modification)
11. [Edge Cases & Error Handling](#edge-cases--error-handling)
12. [Implementation Checklist](#implementation-checklist)

---

## Overview

The Delayed Commit feature allows users to schedule commits outside of restricted time ranges (e.g., work hours). When a commit is attempted during restricted hours, the user will be presented with a time selection UI to choose when the commit should appear to have been made.

### Key Goals
- Prevent commits during defined work hours (configurable)
- Provide intuitive time selection interface
- Maintain Git commit history integrity
- Enable/disable feature via configuration
- Use Git's native `--date` flag for timestamp control

---

## Requirements Summary

Based on user requirements and clarifications:

1. **Time Restrictions**: Single daily range (e.g., 9h-17h)
2. **Time Intervals**: Every 20 minutes for suggested times
3. **Date Handling**: Today only (use immediate time if no available time today)
4. **Default Behavior**: Show time selection UI immediately when committing during restricted hours
5. **Git Integration**: Use `git commit --date` for custom timestamps
6. **Configuration**: Toggle feature on/off in [`~/.gocommit.json`](config/config.go:11)

---

## Architecture Design

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Main Flow                            │
│  (main.go lines 925-1019)                                   │
└────────────┬────────────────────────────────────────────────┘
             │
             ├─→ 1. Get staged changes (git diff --cached)
             ├─→ 2. Generate commit messages via AI
             ├─→ 3. User selects/edits commit message
             │
             ├─→ 4. *** NEW: Check Delayed Commit Feature ***
             │        │
             │        ├─→ Is feature enabled?
             │        │    └─→ No: Proceed to commit (step 6)
             │        │
             │        ├─→ Is current time in restricted range?
             │        │    └─→ No: Proceed to commit (step 6)
             │        │
             │        └─→ Yes to both: Show Time Selection UI
             │             │
             │             ├─→ Generate suggested times (20min intervals)
             │             ├─→ Display time options in termbox UI
             │             ├─→ Allow manual time entry
             │             └─→ Get selected timestamp
             │
             ├─→ 5. *** NEW: Prepare commit command with timestamp ***
             │        │
             │        └─→ Build git commit command with --date flag
             │
             └─→ 6. Execute commit & show success message
```

### Component Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                     config/config.go                         │
│  - Extended Config struct with DelayedCommit fields          │
│  - New getter/setter functions for delayed commit config     │
└──────────────────┬───────────────────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────────────────┐
│                      main.go                                  │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ New Functions:                                         │  │
│  │  - isTimeInRestrictedRange(currentTime, config)        │  │
│  │  - generateSuggestedTimes(config)                      │  │
│  │  - getCommitTimestamp(config)                          │  │
│  │  - showTimeSelectionUI(suggestedTimes)                 │  │
│  │  - validateTimeInput(timeStr)                          │  │
│  │  - formatTimeForGit(timestamp)                         │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                               │
│  Modified Functions:                                          │
│  - main() - Add delayed commit logic before git commit       │
└───────────────────────────────────────────────────────────────┘
```

---

## Configuration Schema

### Extended Config Structure

Add to [`config/config.go:14`](config/config.go:14):

```go
type Config struct {
    APIKey         string              `json:"api_key"`
    LoggingEnabled bool                `json:"logging_enabled"`
    IconMode       bool                `json:"icon_mode"`
    DelayedCommit  DelayedCommitConfig `json:"delayed_commit"` // NEW
}

type DelayedCommitConfig struct {
    Enabled              bool   `json:"enabled"`
    RestrictedStartHour  int    `json:"restricted_start_hour"`  // 24-hour format (0-23)
    RestrictedEndHour    int    `json:"restricted_end_hour"`    // 24-hour format (0-23)
    SuggestionIntervalMin int   `json:"suggestion_interval_min"` // Minutes between suggestions
}
```

### Configuration File Example

`~/.gocommit.json`:
```json
{
  "api_key": "AIzaXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "logging_enabled": true,
  "icon_mode": false,
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval_min": 20
  }
}
```

### Default Configuration Values

```go
func NewDefaultDelayedCommitConfig() DelayedCommitConfig {
    return DelayedCommitConfig{
        Enabled:              false,  // Disabled by default
        RestrictedStartHour:  9,      // 9 AM
        RestrictedEndHour:    17,     // 5 PM
        SuggestionIntervalMin: 20,    // Every 20 minutes
    }
}
```

---

## Data Structures

### Time Option Structure

```go
// TimeOption represents a suggested commit time option
type TimeOption struct {
    DisplayText string    // Human-readable: "18:00 (6:00 PM)"
    Timestamp   time.Time // Actual timestamp to use
    IsManual    bool      // True if this is the "Enter custom time" option
}
```

### Time Selection Result

```go
// TimeSelectionResult contains the user's time choice
type TimeSelectionResult struct {
    Selected  bool      // Whether user selected a time (false if cancelled)
    Timestamp time.Time // The chosen timestamp
    UseNow    bool      // True if user wants to use current time anyway
}
```

---

## Function Signatures

### Configuration Functions (config/config.go)

```go
// Get delayed commit configuration
func GetDelayedCommitConfig() (DelayedCommitConfig, error)

// Set delayed commit configuration
func SetDelayedCommitConfig(cfg DelayedCommitConfig) error

// Enable/disable delayed commit feature
func SetDelayedCommitEnabled(enabled bool) error

// Check if delayed commit is enabled
func IsDelayedCommitEnabled() (bool, error)

// Set restricted hours
func SetRestrictedHours(startHour, endHour int) error
```

### Time Validation Functions (main.go)

```go
// Check if current time falls within restricted range
// Returns: (isRestricted bool, restrictionEnd time.Time, err error)
func isTimeInRestrictedRange(currentTime time.Time, config DelayedCommitConfig) (bool, time.Time, error)

// Validate hour is in valid range (0-23)
func validateHour(hour int) error

// Parse user-entered time string (formats: "18:00", "18h", "6pm", "18:20")
// Returns: (hour int, minute int, err error)
func parseTimeString(timeStr string) (int, int, error)
```

### Time Generation Functions (main.go)

```go
// Generate list of suggested commit times outside restricted hours
// Returns times for today only, starting after restriction ends
func generateSuggestedTimes(config DelayedCommitConfig, currentTime time.Time) []TimeOption

// Create a TimeOption with formatted display text
func createTimeOption(timestamp time.Time) TimeOption

// Get the next available time after restrictions (same day only)
// Returns: (nextTime time.Time, isToday bool)
func getNextAvailableTime(currentTime time.Time, config DelayedCommitConfig) (time.Time, bool)
```

### UI Functions (main.go)

```go
// Show time selection interface with suggested times
// Similar to getUserChoice() but for time selection
func showTimeSelectionUI(options []TimeOption, restrictionEndTime time.Time) (TimeSelectionResult, error)

// Draw time selection screen in termbox
func drawTimeOptions(options []TimeOption, selected int, customTimeInput string, showingCustomInput bool)

// Handle custom time input mode
func handleCustomTimeInput() (time.Time, error)

// Format timestamp for display (e.g., "18:00 (6:00 PM) - Today")
func formatTimeDisplay(t time.Time) string
```

### Git Command Functions (main.go)

```go
// Format timestamp for git --date flag
// Uses ISO 8601 format: "2024-10-21T18:00:00+07:00"
func formatTimeForGit(timestamp time.Time) string

// Execute git commit with custom timestamp
// Uses both --date (author date) and environment variable for committer date
func executeDelayedCommit(message string, timestamp time.Time) error
```

### Main Integration Function (main.go)

```go
// Main function to handle delayed commit workflow
// Called after user selects commit message, before actual commit
// Returns: (timestamp time.Time, shouldUseCustomTime bool, err error)
func handleDelayedCommit(config DelayedCommitConfig) (time.Time, bool, error)
```

---

## Workflow Diagrams

### Overall Flow

```
┌─────────────────────────────────────────────────────────────┐
│ START: User runs gocommit                                   │
└────────────┬────────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────────┐
│ Get git diff & generate commit messages                    │
└────────────┬───────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────────┐
│ User selects/edits commit message                          │
└────────────┬───────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────────┐
│ Check: Is delayed commit enabled?                          │
└────────┬───────────────────────────────────────┬───────────┘
         │ NO                                     │ YES
         │                                        │
         ▼                                        ▼
    ┌────────────┐              ┌────────────────────────────┐
    │ Use current│              │ Check: Is current time in  │
    │ time       │              │ restricted range?          │
    └─────┬──────┘              └────────┬──────────┬────────┘
          │                              │ NO       │ YES
          │                              │          │
          │                              ▼          ▼
          │                         ┌────────┐  ┌──────────────┐
          │                         │ Use    │  │ Show Time    │
          │                         │current │  │ Selection UI │
          │                         │time    │  └──────┬───────┘
          │                         └───┬────┘         │
          │                             │              │
          │                             │              ▼
          │                             │  ┌──────────────────────┐
          │                             │  │ User selects time or │
          │                             │  │ enters custom time   │
          │                             │  └──────────┬───────────┘
          │                             │             │
          └─────────────────────────────┴─────────────┘
                                        │
                                        ▼
                        ┌───────────────────────────────┐
                        │ Execute git commit            │
                        │ with selected timestamp       │
                        └───────────────┬───────────────┘
                                        │
                                        ▼
                        ┌───────────────────────────────┐
                        │ Show success message          │
                        │ END                           │
                        └───────────────────────────────┘
```

### Time Selection UI Flow

```
┌──────────────────────────────────────────────────────────────┐
│ Time Selection UI                                            │
└────────────┬─────────────────────────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────────────────────────┐
│ Display:                                                      │
│  - "Commit during work hours (9:00-17:00) detected"         │
│  - "Select commit time:"                                     │
│  - Suggested times (20min intervals after restriction ends)  │
│    • 17:20 (5:20 PM)                                        │
│    • 17:40 (5:40 PM)                                        │
│    • 18:00 (6:00 PM)     ← Selected                         │
│    • 18:20 (6:20 PM)                                        │
│    • 18:40 (6:40 PM)                                        │
│    • 19:00 (7:00 PM)                                        │
│    • Enter custom time...                                    │
│  - Instructions: "↑↓: Move  Enter: Select  Esc: Cancel"     │
└────────────┬─────────────────────────────────────────────────┘
             │
             ├─→ User presses ↑/↓: Update selection
             │
             ├─→ User presses Enter on time: Return selected time
             │
             ├─→ User presses Enter on "Enter custom time...":
             │   ┌──────────────────────────────────────────────┐
             │   │ Show custom time input prompt:               │
             │   │ "Enter time (HH:MM or HH:MMpm): _"          │
             │   └────────┬─────────────────────────────────────┘
             │            │
             │            ├─→ User enters valid time: Return parsed time
             │            ├─→ User enters invalid time: Show error, retry
             │            └─→ User presses Esc: Return to selection
             │
             └─→ User presses Esc: Cancel, use current time anyway
```

### Time Validation Logic

```
┌──────────────────────────────────────────────────────────────┐
│ isTimeInRestrictedRange(currentTime, config)                 │
└────────────┬─────────────────────────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────────────────────────┐
│ Extract hour from currentTime                                 │
└────────────┬─────────────────────────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────────────────────────┐
│ Is hour >= RestrictedStartHour AND                           │
│    hour < RestrictedEndHour ?                                │
└────────┬───────────────────────────────────────┬─────────────┘
         │ YES                                    │ NO
         │                                        │
         ▼                                        ▼
┌──────────────────────┐              ┌──────────────────────┐
│ Calculate restriction│              │ Not in restricted    │
│ end time for today   │              │ range, proceed with  │
│ Return: true,        │              │ current time         │
│   endTime, nil       │              │ Return: false,       │
└──────────────────────┘              │   zero, nil          │
                                      └──────────────────────┘
```

---

## Integration Points

### 1. Configuration Integration (config/config.go)

**Location**: After [`SetIconMode()`](config/config.go:122)

**New Functions to Add**:
```go
// Lines 131-160 (approximate)
func GetDelayedCommitConfig() (DelayedCommitConfig, error) {
    config, err := LoadConfig()
    if err != nil {
        return DelayedCommitConfig{}, err
    }
    
    // Return default config if not set
    if config.DelayedCommit.RestrictedStartHour == 0 && 
       config.DelayedCommit.RestrictedEndHour == 0 {
        return NewDefaultDelayedCommitConfig(), nil
    }
    
    return config.DelayedCommit, nil
}

func SetDelayedCommitConfig(cfg DelayedCommitConfig) error {
    config, err := LoadConfig()
    if err != nil {
        return err
    }
    
    config.DelayedCommit = cfg
    return SaveConfig(config)
}

// Additional helper functions...
```

### 2. Main Flow Integration (main.go)

**Location**: Between line 997 (after [`getUserChoice()`](main.go:994)) and line 1006 (before git commit)

**Integration Code**:
```go
// Lines 998-1004 (new code block)

// Handle delayed commit if enabled
timestamp := time.Now()
useCustomTimestamp := false

delayedConfig, err := config.GetDelayedCommitConfig()
if err != nil {
    log.Printf("Warning: Failed to get delayed commit config: %v", err)
} else if delayedConfig.Enabled {
    customTime, shouldDelay, err := handleDelayedCommit(delayedConfig)
    if err != nil {
        log.Fatal("Delayed commit error:", err)
    }
    if shouldDelay {
        timestamp = customTime
        useCustomTimestamp = true
    }
}

// Get the full AI response for logging
aiResponse := strings.Join(messages, "\n")

// Log successful request
logger.LogSuccess(gitDiff, lastCommitMsg, prompt, aiResponse, messages, commitMsg)

// Create git commit with the chosen message (MODIFIED)
var cmd *exec.Cmd
if useCustomTimestamp {
    cmd = executeDelayedCommit(commitMsg, timestamp)
} else {
    cmd = exec.Command("git", "commit", "-m", commitMsg)
}
```

### 3. Command Line Flag Integration (main.go)

**Location**: After line 931 (flag definitions)

**New Flags**:
```go
// Line 932 (add new flags)
configDelayedFlag := flag.Bool("config-delayed", false, "Configure delayed commit settings")
enableDelayedFlag := flag.Bool("enable-delayed", false, "Enable delayed commit feature")
disableDelayedFlag := flag.Bool("disable-delayed", false, "Disable delayed commit feature")
```

**Flag Handler** (after line 948):
```go
// Lines 949-975 (new flag handlers)
if *configDelayedFlag {
    // Interactive configuration
    fmt.Println("Configure Delayed Commit Feature")
    fmt.Print("Enter restricted start hour (0-23, e.g., 9 for 9 AM): ")
    var startHour int
    fmt.Scanln(&startHour)
    
    fmt.Print("Enter restricted end hour (0-23, e.g., 17 for 5 PM): ")
    var endHour int
    fmt.Scanln(&endHour)
    
    fmt.Print("Enter suggestion interval in minutes (e.g., 20): ")
    var interval int
    fmt.Scanln(&interval)
    
    cfg := config.DelayedCommitConfig{
        Enabled:              true,
        RestrictedStartHour:  startHour,
        RestrictedEndHour:    endHour,
        SuggestionIntervalMin: interval,
    }
    
    if err := config.SetDelayedCommitConfig(cfg); err != nil {
        log.Fatalf("Failed to save delayed commit config: %v", err)
    }
    fmt.Println("Delayed commit configured successfully!")
    return
}

if *enableDelayedFlag {
    if err := config.SetDelayedCommitEnabled(true); err != nil {
        log.Fatalf("Failed to enable delayed commit: %v", err)
    }
    fmt.Println("Delayed commit enabled!")
    return
}

if *disableDelayedFlag {
    if err := config.SetDelayedCommitEnabled(false); err != nil {
        log.Fatalf("Failed to disable delayed commit: %v", err)
    }
    fmt.Println("Delayed commit disabled!")
    return
}
```

---

## User Interface Flow

### Time Selection Screen Layout

```
┌──────────────────────────────────────────────────────────────┐
│ Commit during work hours (09:00-17:00) detected             │
│ Current time: 14:30                                          │
│                                                              │
│ Select commit time:                                          │
│                                                              │
│  • 17:20 (5:20 PM) - Today                                  │
│  • 17:40 (5:40 PM) - Today                                  │
│  → 18:00 (6:00 PM) - Today                                  │
│  • 18:20 (6:20 PM) - Today                                  │
│  • 18:40 (6:40 PM) - Today                                  │
│  • 19:00 (7:00 PM) - Today                                  │
│  • 19:20 (7:20 PM) - Today                                  │
│  • Enter custom time...                                      │
│                                                              │
│ ↑↓: Move  Enter: Select  Esc: Use current time anyway       │
└──────────────────────────────────────────────────────────────┘
```

### Custom Time Input Screen

```
┌──────────────────────────────────────────────────────────────┐
│ Enter custom commit time                                     │
│                                                              │
│ Time: 18:45_                                                │
│                                                              │
│ Accepted formats:                                            │
│  - HH:MM (24-hour): 18:45, 09:30                           │
│  - HHhMM: 18h45, 9h30                                       │
│  - HH:MM AM/PM: 6:45 PM, 9:30 AM                           │
│                                                              │
│ Enter: Confirm  Esc: Cancel                                  │
└──────────────────────────────────────────────────────────────┘
```

### UI Implementation Using Existing Framework

The UI will use the same termbox-go framework as [`getUserChoice()`](main.go:851) and [`editMessage()`](main.go:615), providing consistent user experience.

**Key Functions**:
- `showTimeSelectionUI()` - Similar to `getUserChoice()` but for time options
- `drawTimeOptions()` - Similar to `drawMessages()` but for time display
- `handleCustomTimeInput()` - Similar to `editMessage()` but simpler (time only)

---

## Git Command Modification

### Using git commit --date

Git provides the `--date` flag to set custom commit dates. There are two timestamps in Git:
1. **Author Date**: When the changes were made (controlled by `--date`)
2. **Committer Date**: When the commit was recorded (controlled by `GIT_COMMITTER_DATE` env var)

### Implementation Strategy

**Option 1: Set Both Dates (Recommended)**
```go
func executeDelayedCommit(message string, timestamp time.Time) *exec.Cmd {
    // Format timestamp for git (ISO 8601)
    dateStr := timestamp.Format(time.RFC3339)
    
    // Create command with --date flag
    cmd := exec.Command("git", "commit", "-m", message, "--date", dateStr)
    
    // Also set committer date via environment variable
    cmd.Env = append(os.Environ(), 
        fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))
    
    return cmd
}
```

**Option 2: Set Author Date Only (Simpler)**
```go
func executeDelayedCommit(message string, timestamp time.Time) *exec.Cmd {
    dateStr := timestamp.Format(time.RFC3339)
    return exec.Command("git", "commit", "-m", message, "--date", dateStr)
}
```

### Date Format Options

Git accepts multiple date formats. Recommended: **ISO 8601 (RFC3339)**

```go
func formatTimeForGit(t time.Time) string {
    // ISO 8601 format: "2024-10-21T18:00:00+07:00"
    return t.Format(time.RFC3339)
}
```

**Alternative formats** (if needed):
- Unix timestamp: `t.Unix()` → `"@1697904000"`
- Git internal: `t.Format("Mon Jan 2 15:04:05 2006 -0700")`
- Relative: `"2 hours ago"` (not recommended for this use case)

### Modified Commit Execution

**Current code** (main.go:1006):
```go
cmd := exec.Command("git", "commit", "-m", commitMsg)
```

**New code** (with delayed commit support):
```go
var cmd *exec.Cmd
if useCustomTimestamp {
    dateStr := timestamp.Format(time.RFC3339)
    cmd = exec.Command("git", "commit", "-m", commitMsg, "--date", dateStr)
    cmd.Env = append(os.Environ(), 
        fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))
} else {
    cmd = exec.Command("git", "commit", "-m", commitMsg)
}
```

### Verification

Users can verify the commit timestamp using:
```bash
git log --format=fuller
```

This shows both author and committer dates:
```
commit abc123...
Author:     User Name <user@email.com>
AuthorDate: Mon Oct 21 18:00:00 2024 +0700
Commit:     User Name <user@email.com>
CommitDate: Mon Oct 21 18:00:00 2024 +0700

    feat: add new feature
```

---

## Edge Cases & Error Handling

### 1. Configuration Edge Cases

| Case | Handling |
|------|----------|
| **Invalid hour range** (e.g., start=20, end=8) | Validate on save: return error "Start hour must be less than end hour" |
| **Same start/end hour** | Allow (effectively disables restriction for that hour) |
| **Hours outside 0-23** | Validate on save: return error "Hours must be between 0 and 23" |
| **Negative interval** | Use default (20 minutes) |
| **Zero interval** | Use default (20 minutes) |
| **Very large interval** (>60min) | Allow but warn user |
| **Missing config** | Use default values silently |

### 2. Time Selection Edge Cases

| Case | Handling |
|------|----------|
| **Commit at 16:55 (5min before restriction ends)** | Still show UI; include 17:00, 17:20, etc. |
| **Commit at 23:50 (near midnight)** | Only show remaining times today; if none, show "No times available today, using current time" |
| **All suggested times in past** | Show "No future times available today" and use current time |
| **User selects time in the past** | Show warning: "Selected time is in the past. Continue anyway? (Y/n)" |
| **User selects time within restricted range** | Show error: "Time must be outside restricted hours (9:00-17:00)" and return to selection |
| **Custom time input is empty** | Return to selection menu |
| **Custom time input is invalid format** | Show error message: "Invalid format. Use HH:MM (e.g., 18:30)" and allow retry |

### 3. Git Command Edge Cases

| Case | Handling |
|------|----------|
| **git commit fails with custom date** | Fall back to current time and retry |
| **Timestamp in future** | Allow (Git supports this) |
| **Timestamp very far in past** (>1 year) | Allow but log warning |
| **Time zone issues** | Always use local timezone from `time.Now()` |

### 4. UI Edge Cases

| Case | Handling |
|------|----------|
| **Terminal too small** | Use minimum width (80 chars), truncate if needed |
| **User presses Esc during time selection** | Cancel selection, use current time (show confirmation) |
| **No suggested times generated** | Show "No times available" and manual entry option only |
| **Custom input screen exits unexpectedly** | Return to main selection menu |

### 5. Concurrent/Race Conditions

| Case | Handling |
|------|----------|
| **Config file modified during execution** | Load config once at start, don't reload |
| **System time changes during execution** | Use captured time from start of function |

### Error Message Examples

```go
// Time validation errors
"Invalid time format. Please use HH:MM (e.g., 18:30)"
"Time must be outside restricted hours (09:00-17:00)"
"Selected time is in the past. Continue anyway?"

// Configuration errors
"Start hour must be less than end hour"
"Hours must be between 0 and 23"
"Failed to load delayed commit configuration"

// Git command errors
"Failed to create commit with custom timestamp"
"Git command failed. Retrying with current time..."

// UI errors
"Time selection cancelled by user"
"No times available today outside restricted hours"
```

### Error Recovery Strategy

```go
func handleDelayedCommit(config DelayedCommitConfig) (time.Time, bool, error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic in delayed commit: %v", r)
            // Return current time as fallback
        }
    }()
    
    // Main logic with error handling...
    // On any critical error, fall back to current time
}
```

---

## Implementation Checklist

### Phase 1: Configuration (config/config.go)

- [ ] Add `DelayedCommitConfig` struct to `Config`
- [ ] Implement `GetDelayedCommitConfig()`
- [ ] Implement `SetDelayedCommitConfig()`
- [ ] Implement `SetDelayedCommitEnabled()`
- [ ] Implement `IsDelayedCommitEnabled()`
- [ ] Implement `SetRestrictedHours()`
- [ ] Add validation functions for hours
- [ ] Write unit tests for config functions

### Phase 2: Time Validation Logic (main.go)

- [ ] Implement `isTimeInRestrictedRange()`
- [ ] Implement `validateHour()`
- [ ] Implement `parseTimeString()` with multiple format support
- [ ] Implement `getNextAvailableTime()`
- [ ] Write unit tests for time validation
- [ ] Handle edge cases (midnight, invalid ranges)

### Phase 3: Time Generation (main.go)

- [ ] Implement `generateSuggestedTimes()`
- [ ] Implement `createTimeOption()`
- [ ] Implement `formatTimeDisplay()`
- [ ] Ensure times are only for today
- [ ] Handle no-times-available scenario
- [ ] Write unit tests for time generation

### Phase 4: User Interface (main.go)

- [ ] Implement `showTimeSelectionUI()`
- [ ] Implement `drawTimeOptions()` using termbox
- [ ] Implement `handleCustomTimeInput()`
- [ ] Add keyboard navigation (↑↓, Enter, Esc)
- [ ] Add visual feedback (highlighting, colors)
- [ ] Handle terminal resize events
- [ ] Test UI with various terminal sizes

### Phase 5: Git Integration (main.go)

- [ ] Implement `formatTimeForGit()`
- [ ] Implement `executeDelayedCommit()`
- [ ] Modify main() to check delayed commit config
- [ ] Modify main() to call `handleDelayedCommit()`
- [ ] Set both author and committer dates
- [ ] Test with actual git repository
- [ ] Verify timestamps in git log

### Phase 6: Command Line Flags (main.go)

- [ ] Add `--config-delayed` flag
- [ ] Add `--enable-delayed` flag
- [ ] Add `--disable-delayed` flag
- [ ] Implement interactive configuration
- [ ] Add help text for new flags
- [ ] Update documentation

### Phase 7: Error Handling & Edge Cases

- [ ] Add validation for all user inputs
- [ ] Implement error recovery mechanisms
- [ ] Add user-friendly error messages
- [ ] Handle terminal/UI errors gracefully
- [ ] Add logging for debugging
- [ ] Test all edge cases listed above

### Phase 8: Testing & Documentation

- [ ] Write unit tests for all new functions
- [ ] Write integration tests
- [ ] Test with real git repositories
- [ ] Update README.md with feature documentation
- [ ] Add examples to documentation
- [ ] Create user guide for delayed commit feature
- [ ] Update INSTALL.md if needed

### Phase 9: Polish & Optimization

- [ ] Code review and refactoring
- [ ] Performance optimization
- [ ] Memory leak checks
- [ ] Add metrics/telemetry (if applicable)
- [ ] Final testing on multiple platforms
- [ ] Prepare release notes

---

## Example Usage Scenarios

### Scenario 1: Committing During Work Hours

```bash
$ gocommit
# (User has staged changes)
# AI generates commit messages
# User selects: "feat(api): add user authentication endpoint"

┌──────────────────────────────────────────────────────────────┐
│ Commit during work hours (09:00-17:00) detected             │
│ Current time: 14:30                                          │
│                                                              │
│ Select commit time:                                          │
│  → 17:20 (5:20 PM) - Today                                  │
│  • 17:40 (5:40 PM) - Today                                  │
│  • 18:00 (6:00 PM) - Today                                  │
│  ...                                                         │
└──────────────────────────────────────────────────────────────┘

# User selects 18:00
Successfully created commit with message: feat(api): add user authentication endpoint
Commit timestamp: 2024-10-21 18:00:00 +0700
```

### Scenario 2: Custom Time Entry

```bash
$ gocommit
# ...after message selection...

┌──────────────────────────────────────────────────────────────┐
│ Select commit time:                                          │
│  • 17:20 (5:20 PM) - Today                                  │
│  • 17:40 (5:40 PM) - Today                                  │
│  → Enter custom time...                                      │
└──────────────────────────────────────────────────────────────┘

# User selects "Enter custom time..."

┌──────────────────────────────────────────────────────────────┐
│ Enter custom commit time                                     │
│ Time: 19:15_                                                │
└──────────────────────────────────────────────────────────────┘

Successfully created commit with message: feat(api): add user authentication endpoint
Commit timestamp: 2024-10-21 19:15:00 +0700
```

### Scenario 3: Committing Outside Restricted Hours

```bash
$ gocommit
# Current time: 19:30 (outside 9-17 restriction)
# No time selection UI shown
# Proceeds directly to commit with current timestamp

Successfully created commit with message: feat(api): add user authentication endpoint
```

### Scenario 4: Configuring the Feature

```bash
$ gocommit --config-delayed
Configure Delayed Commit Feature
Enter restricted start hour (0-23, e.g., 9 for 9 AM): 9
Enter restricted end hour (0-23, e.g., 17 for 5 PM): 17
Enter suggestion interval in minutes (e.g., 20): 20
Delayed commit configured successfully!

$ gocommit --enable-delayed
Delayed commit enabled!

$ gocommit --disable-delayed
Delayed commit disabled!
```

---

## Architecture Diagrams

### System Context Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        User                                  │
└────────────┬────────────────────────────────────────────────┘
             │
             │ Runs gocommit
             ▼
┌──────────────────────────────────────────────────────────────┐
│                     gocommit CLI                             │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  Main Application (main.go)                           │  │
│  │  - Git operations                                      │  │
│  │  - AI message generation                               │  │
│  │  - User interaction (termbox UI)                       │  │
│  │  - Delayed commit logic ← NEW                         │  │
│  └─────────────┬────────────────────┬─────────────────────┘  │
│                │                    │                         │
│  ┌─────────────▼──────────┐  ┌─────▼──────────────────────┐  │
│  │ Config Module          │  │ Logger Module              │  │
│  │ (config/config.go)     │  │ (logger/logger.go)         │  │
│  │ - Load/Save config     │  │ - Request logging          │  │
│  │ - Delayed commit cfg ←│  │                            │  │
│  └────────────────────────┘  └────────────────────────────┘  │
└────────┬────────────────────────────────────────┬────────────┘
         │                                        │
         │ Reads/Writes                          │ Executes
         ▼                                        ▼
┌────────────────────┐                  ┌─────────────────────┐
│  ~/.gocommit.json  │                  │   Git Repository    │
│  Configuration     │                  │   (.git/)           │
└────────────────────┘                  └─────────────────────┘
```

### Component Interaction Diagram

```
main() ──┐
         │
         ├─→ config.GetDelayedCommitConfig() ──→ Load from ~/.gocommit.json
         │                                      
         ├─→ getGitDiff() ──────────────────→ git diff --cached
         │                                      
         ├─→ generateCommitMessages() ───────→ Google AI API
         │                                      
         ├─→ getUserChoice() ────────────────→ termbox UI (message selection)
         │                                      
         ├─→ handleDelayedCommit() ──┐
         │                           │
         │                           ├─→ isTimeInRestrictedRange()
         │                           │
         │                           ├─→ generateSuggestedTimes()
         │                           │
         │                           └─→ showTimeSelectionUI() ──→ termbox UI (time selection)
         │                                   │
         │                                   └─→ handleCustomTimeInput()
         │                                      
         ├─→ logger.LogSuccess() ────────────→ ~/.gocommit/gocommit_requests.log
         │                                      
         └─→ executeDelayedCommit() ─────────→ git commit --date
             OR
             exec.Command("git", "commit")
```

### Data Flow Diagram

```
User Input
    │
    ▼
┌─────────────────┐
│ Staged Changes  │
└────────┬────────┘
         │
         ▼
┌──────────────────────────────────────┐
│ AI generates 3 commit messages       │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│ User selects/edits message           │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│ Load delayed commit config           │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│ Is feature enabled?                  │
└──┬──────────────────────────────┬────┘
   │ NO                           │ YES
   │                              │
   │                              ▼
   │                   ┌─────────────────────┐
   │                   │ Is time restricted? │
   │                   └──┬──────────────┬───┘
   │                      │ NO           │ YES
   │                      │              │
   │                      │              ▼
   │                      │   ┌──────────────────────┐
   │                      │   │ Generate time options│
   │                      │   └──────────┬───────────┘
   │                      │              │
   │                      │              ▼
   │                      │   ┌──────────────────────┐
   │                      │   │ Show time UI         │
   │                      │   └──────────┬───────────┘
   │                      │              │
   │                      │              ▼
   │                      │   ┌──────────────────────┐
   │                      │   │ Get user selection   │
   │                      │   └──────────┬───────────┘
   │                      │              │
   └──────────────────────┴──────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │ Commit with timestamp │
              └───────────┬───────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │ Update git repository │
              └───────────────────────┘
```

---

## Summary

This architecture provides a comprehensive solution for the delayed commit feature with:

✅ **Clear Configuration**: Extended [`Config struct`](config/config.go:14) with dedicated delayed commit settings
✅ **Robust Time Logic**: Comprehensive time validation and suggestion generation  
✅ **Intuitive UI**: Consistent termbox interface matching existing patterns
✅ **Git Integration**: Proper use of `git commit --date` for timestamp control
✅ **Error Handling**: Extensive edge case coverage with graceful fallbacks
✅ **Extensibility**: Clean separation of concerns for future enhancements

The design is ready for implementation with clear function signatures, integration points, and testing guidelines. All new code maintains consistency with the existing gocommit codebase patterns and conventions.