# Delayed Commit Feature Documentation

## Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [Time Formats](#time-formats)
- [How It Works](#how-it-works)
- [Command Reference](#command-reference)
- [FAQ / Troubleshooting](#faq--troubleshooting)
- [Best Practices](#best-practices)
- [Configuration Examples](#configuration-examples)

---

## Overview

The Delayed Commit feature allows you to schedule commits outside of restricted time ranges (e.g., work hours). When you attempt to commit during restricted hours, you'll be presented with an interactive UI to choose when the commit should appear to have been made.

### Why Use This Feature?

- **Maintain Consistent Commit Timestamps**: Keep your git history clean by avoiding commits during work hours
- **Privacy & Work-Life Balance**: Prevent revealing your actual working patterns through commit timestamps
- **Professional Appearance**: Ensure commits appear at appropriate times (e.g., avoid 2 AM commits)
- **Timezone Management**: Control when your commits appear regardless of when you actually work

### Key Benefits

✅ **Flexible Time Selection**: Choose from suggested times or enter custom times  
✅ **Multiple Format Support**: Enter times in 24-hour, 12-hour, or shorthand formats  
✅ **Git-Native Integration**: Uses Git's built-in `--date` flag for authentic timestamps  
✅ **Configurable Restrictions**: Set your own restricted hours and suggestion intervals  
✅ **Easy Toggle**: Enable/disable the feature with simple commands

---

## Quick Start

### 1. Configure Delayed Commit Settings

```bash
gocommit --config-delayed
```

You'll be prompted to enter:
- **Restricted start hour** (0-23): e.g., `9` for 9 AM
- **Restricted end hour** (0-23): e.g., `17` for 5 PM  
- **Suggestion interval** (minutes): e.g., `20`

Example session:
```
Configure Delayed Commit Settings
===================================

Enter restricted start hour (0-23, e.g., 9 for 9 AM): 9
Enter restricted end hour (0-23, e.g., 17 for 5 PM): 17
Enter suggestion interval in minutes (e.g., 20, 30, 60): 20

Delayed commit settings configured successfully!
Restricted hours: 09:00 - 17:00
Suggestion interval: 20 minutes
```

### 2. Enable the Feature

```bash
gocommit --enable-delayed
```

Expected output:
```
Delayed commit feature enabled successfully!
```

### 3. Make a Commit During Restricted Hours

```bash
# Stage your changes
git add .

# Run gocommit during restricted hours (e.g., at 2 PM)
gocommit
```

You'll see the time selection interface:
```
┌──────────────────────────────────────────────────────────────┐
│ Commit during restricted hours detected                      │
│ Restricted hours: 09:00 - 17:00                             │
│                                                              │
│ Select commit time:                                          │
│  → 17:20                                                     │
│  • 17:40                                                     │
│  • 18:00                                                     │
│  • 18:20                                                     │
│  • 18:40                                                     │
│  • 19:00                                                     │
│  • Enter custom time (c)                                     │
│                                                              │
│ ↑↓: Move  Enter: Select  c: Custom  Esc: Cancel             │
└──────────────────────────────────────────────────────────────┘
```

---

## Configuration

### Configure Using Command Line

Use the [`--config-delayed`](#gocommit---config-delayed) flag to configure settings interactively:

```bash
gocommit --config-delayed
```

### Configuration Options

| Option | Description | Valid Range | Default |
|--------|-------------|-------------|---------|
| **Restricted Start Hour** | Hour when restrictions begin | 0-23 | 9 |
| **Restricted End Hour** | Hour when restrictions end | 0-23 | 17 |
| **Suggestion Interval** | Minutes between suggested times | 1-1440 | 20 |

### Configuration File

Settings are stored in [`~/.gocommit.json`](config/config.go:11):

```json
{
  "api_key": "AIzaXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
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

### Enable/Disable Feature

**Enable:**
```bash
gocommit --enable-delayed
```

**Disable:**
```bash
gocommit --disable-delayed
```

**Check Current Status:**  
View the config file:
```bash
cat ~/.gocommit.json
```

---

## Usage Examples

### Example 1: Basic Usage During Restricted Hours

**Scenario**: Committing at 2:30 PM (14:30) when restricted hours are 9 AM - 5 PM

```bash
$ gocommit
⏰ Commit during restricted hours detected!
Restricted hours: 09:00 - 17:00

# Time selection UI appears
# User selects 18:00 using arrow keys and Enter

Scheduling commit for: 18:00
✓ Successfully created commit with message: feat(api): add user authentication endpoint
Commit timestamp: 2024-10-21 18:00:00 +0700
```

### Example 2: Selecting from Suggested Times

**Scenario**: Navigate through suggested times and select one

```bash
# Time selection UI with 20-minute intervals
┌──────────────────────────────────────────────────────────────┐
│ Select commit time:                                          │
│  • 17:20                                                     │
│  → 17:40    ← Currently selected                            │
│  • 18:00                                                     │
│  • 18:20                                                     │
└──────────────────────────────────────────────────────────────┘

# Press Enter to confirm
Scheduling commit for: 17:40
Commit timestamp: 2024-10-21 17:40:00 +0700
```

### Example 3: Entering a Custom Time

**Scenario**: Use custom time instead of suggestions

```bash
# In time selection UI, press 'c' or select "Enter custom time"
┌──────────────────────────────────────────────────────────────┐
│ Enter custom commit time                                     │
│                                                              │
│ Accepted formats:                                            │
│  - HH:MM (24-hour): 18:45, 09:30                           │
│  - HHhMM: 18h45, 9h30                                       │
│  - HH:MM AM/PM: 6:45 PM, 9:30 AM                           │
│                                                              │
│ Time: 19:15_                                                │
│                                                              │
│ Enter: Confirm  Esc: Cancel                                  │
└──────────────────────────────────────────────────────────────┘

# Type custom time and press Enter
Scheduling commit for: 19:15
Commit timestamp: 2024-10-21 19:15:00 +0700
```

### Example 4: Using Outside Restricted Hours

**Scenario**: Committing at 7 PM (19:00) when restricted hours are 9 AM - 5 PM

```bash
$ gocommit
# No time selection UI shown - proceeds directly to commit

✓ Successfully created commit with message: fix(auth): resolve token expiration issue
```

The feature only activates when committing during restricted hours.

### Example 5: Combining with Other Flags

**Scenario**: Using delayed commit with detailed messages and icons

```bash
# Generate detailed commit messages with icons during restricted hours
gocommit -d --icon

# Time selection UI appears (if during restricted hours)
# After selecting time, commit is created with detailed message and icon
```

---

## Time Formats

The delayed commit feature supports multiple time input formats for flexibility:

### 24-Hour Format (Recommended)

Format: `HH:MM`

Examples:
- `18:00` → 6:00 PM
- `09:30` → 9:30 AM
- `23:45` → 11:45 PM
- `00:15` → 12:15 AM

### 12-Hour Format with AM/PM

Format: `HH:MM AM` or `HH:MM PM`

Examples:
- `6:00 PM` → 18:00
- `9:30 AM` → 09:30
- `11:45 PM` → 23:45
- `12:15 AM` → 00:15

### Shorthand Format

Format: `HHhMM`

Examples:
- `18h00` → 18:00
- `9h30` → 09:30
- `23h45` → 23:45
- `0h15` → 00:15

### Format Validation

The parser validates:
- Hours must be 0-23 (24-hour) or 1-12 (12-hour with AM/PM)
- Minutes must be 0-59
- Invalid formats show helpful error messages

**Invalid Examples:**
```bash
25:00      # Error: hour must be between 0 and 23
18:65      # Error: minute must be between 0 and 59
6PM        # Error: invalid time format. Use HH:MM AM/PM
```

---

## How It Works

### Git Timestamp Manipulation

The delayed commit feature uses Git's native [`--date`](main.go:1407) flag to set custom commit timestamps. This is a standard Git feature that allows you to specify when a commit should appear to have been made.

### Two Types of Timestamps

Git tracks two separate timestamps for each commit:

1. **Author Date**: When the changes were originally made
2. **Committer Date**: When the commit was recorded to the repository

The delayed commit feature sets **both** timestamps to your chosen time using:
- `git commit --date` flag for author date
- `GIT_COMMITTER_DATE` environment variable for committer date

### Command Execution

Normal commit:
```bash
git commit -m "feat: add new feature"
```

Delayed commit (behind the scenes):
```bash
# Sets both author and committer dates
GIT_COMMITTER_DATE="2024-10-21T18:00:00+07:00" \
git commit -m "feat: add new feature" \
  --date "2024-10-21T18:00:00+07:00"
```

### Viewing Commit Timestamps

Use `git log --format=fuller` to see both timestamps:

```bash
$ git log --format=fuller -1

commit abc123def456...
Author:     John Doe <john@example.com>
AuthorDate: Mon Oct 21 18:00:00 2024 +0700
Commit:     John Doe <john@example.com>
CommitDate: Mon Oct 21 18:00:00 2024 +0700

    feat: add new feature
```

### Impact on Git History

✅ **Authentic Timestamps**: Uses Git's official timestamp mechanism  
✅ **Full Compatibility**: Works with all Git tools and services (GitHub, GitLab, etc.)  
✅ **No Corruption**: Doesn't modify or corrupt Git history  
✅ **Audit Trail**: Timestamps are permanent and properly recorded

⚠️ **Note**: While the feature modifies timestamps, this is a legitimate Git feature. The commits are authentic and don't violate Git's integrity checks.

---

## Command Reference

### `gocommit --config-delayed`

Configure delayed commit settings interactively.

**Usage:**
```bash
gocommit --config-delayed
```

**Prompts:**
- Restricted start hour (0-23)
- Restricted end hour (0-23)
- Suggestion interval in minutes

**Example:**
```bash
$ gocommit --config-delayed
Configure Delayed Commit Settings
===================================

Enter restricted start hour (0-23, e.g., 9 for 9 AM): 9
Enter restricted end hour (0-23, e.g., 17 for 5 PM): 17
Enter suggestion interval in minutes (e.g., 20, 30, 60): 20

Delayed commit settings configured successfully!
```

---

### `gocommit --enable-delayed`

Enable the delayed commit feature.

**Usage:**
```bash
gocommit --enable-delayed
```

**Output:**
```
Delayed commit feature enabled successfully!
```

**Note**: You must configure settings first using [`--config-delayed`](#gocommit---config-delayed) before enabling.

---

### `gocommit --disable-delayed`

Disable the delayed commit feature.

**Usage:**
```bash
gocommit --disable-delayed
```

**Output:**
```
Delayed commit feature disabled successfully!
```

**Note**: Your configuration is preserved and can be re-enabled later.

---

### `gocommit` (with delayed commit active)

Regular commit command that respects delayed commit settings when enabled.

**Usage:**
```bash
# Stage changes
git add .

# Commit with gocommit
gocommit
```

**Behavior:**
- **Outside restricted hours**: Commits immediately with current timestamp
- **During restricted hours**: Shows time selection UI
- **Feature disabled**: Always commits immediately

**Compatible Flags:**
- `-d` : Detailed commit messages
- `--icon` : Use emoji icons
- All other gocommit flags work normally

---

## FAQ / Troubleshooting

### What happens if I commit outside restricted hours?

The delayed commit feature only activates during restricted hours. If you commit outside the configured time range, gocommit works normally and uses the current timestamp.

```bash
# Restricted hours: 9:00 - 17:00
# Current time: 19:00 (7 PM)
$ gocommit
# Proceeds directly to commit - no time selection UI
```

---

### Can I use this with other gocommit flags?

Yes! The delayed commit feature works seamlessly with all gocommit flags:

```bash
# Detailed messages with delayed commit
gocommit -d

# Icon mode with delayed commit
gocommit --icon

# Both flags together
gocommit -d --icon
```

The time selection UI appears after message selection, before the commit is executed.

---

### How do I check my current configuration?

View your configuration file:

```bash
cat ~/.gocommit.json
```

Or check specific settings:

```bash
# View entire config in pretty format
cat ~/.gocommit.json | json_pp  # On Linux/Mac with json_pp installed

# Or use jq for formatted output
cat ~/.gocommit.json | jq .delayed_commit
```

---

### What if I cancel the time selection?

Press `Esc` during time selection to cancel:

```bash
# Press Esc in time selection UI
Commit cancelled.
```

Your staged changes remain intact - nothing is committed.

---

### Does this affect my git history integrity?

No. The delayed commit feature uses Git's official `--date` flag, which is a standard feature. The commits are authentic and don't violate Git's integrity:

✅ Commit hashes are valid  
✅ Repository integrity is maintained  
✅ All Git tools recognize the commits  
✅ Timestamps are properly recorded in Git's data structures

**Important**: This is different from rewriting history (like `git rebase`). The commits are created once with the specified timestamp, not modified after creation.

---

### Can I set timestamps in the past or future?

**Past timestamps**: Yes, Git allows past timestamps. However, very old timestamps (> 1 year) may look suspicious.

**Future timestamps**: Yes, Git accepts future timestamps. This might be useful for scheduled releases or planning purposes.

**Validation**: The tool validates that the hour/minute are valid (0-23, 0-59) but doesn't restrict the date to today.

---

### What if the commit fails with the custom timestamp?

The tool includes automatic fallback:

```bash
# If commit with custom timestamp fails
Warning: Failed to commit with custom date, using current time: <error>
# Automatically retries with current timestamp
```

This ensures your commit succeeds even if there's an issue with the timestamp.

---

### Why don't I see the time selection UI?

Check these conditions:

1. **Is the feature enabled?**
   ```bash
   cat ~/.gocommit.json | grep "enabled"
   # Should show: "enabled": true
   ```

2. **Are you committing during restricted hours?**
   - Check your configured hours: `cat ~/.gocommit.json`
   - Verify current time matches restricted range

3. **Is the configuration valid?**
   - Start hour and end hour must be different
   - Both hours must be between 0-23

---

### How do I change my restricted hours?

Run the configuration command again:

```bash
gocommit --config-delayed
```

Enter new values when prompted. Previous settings will be overwritten.

---

## Best Practices

### Recommended Settings for Different Use Cases

#### Avoiding Work Hours (9-5 Schedule)
```bash
gocommit --config-delayed
# Start: 9
# End: 17
# Interval: 20
```
**Use case**: Keep commits outside typical work hours for privacy or professional appearance.

---

#### Night Owl Schedule (Avoiding Late Night Commits)
```bash
gocommit --config-delayed
# Start: 22
# End: 6
# Interval: 30
```
**Use case**: Avoid showing 2 AM commits by scheduling them for morning hours.

---

#### Weekend Work (Scheduling for Monday)
```bash
# Temporarily set very wide restricted hours
gocommit --config-delayed
# Start: 0
# End: 23
# Interval: 60
```
**Use case**: When working on weekends, schedule all commits for Monday morning.  
**Note**: Currently only supports same-day scheduling.

---

### When to Enable/Disable

**Enable when:**
- 🔒 Working on professional/public repositories
- 🌙 Working outside normal hours regularly
- 🏢 Maintaining work-life boundaries
- 🌍 Working across timezones

**Disable when:**
- 📝 Accurate timestamps are required for audit purposes
- 🔍 Debugging time-sensitive issues
- 👥 Working with a team that relies on actual timestamps
- 📊 Analyzing actual work patterns

---

### Tips for Choosing Restriction Hours

1. **Match Your Schedule**: Set restricted hours to when you don't want commits to appear
2. **Consider Timezones**: Account for your timezone and your team's timezone
3. **Be Consistent**: Use the same hours regularly to avoid pattern suspicion
4. **Allow Buffer Time**: Set end hour slightly after your target time to allow for flexibility

**Example**: If you want commits to appear after 6 PM:
```
Start: 9 (9 AM)
End: 18 (6 PM)
```
This gives you a buffer - commits will be suggested starting at 6:00 PM.

---

### Tips for Choosing Intervals

**20 minutes** (default): Good balance of options without overwhelming  
**30 minutes**: Cleaner timestamps (e.g., 18:00, 18:30, 19:00)  
**60 minutes**: Minimal options, on-the-hour timestamps only  
**15 minutes**: More granular control, more options to choose from

---

## Configuration Examples

### Example 1: Standard Office Hours

**Scenario**: Avoid commits during 9-5 work hours

```json
{
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval": 20
  }
}
```

**Suggested times** (if committing at 2 PM):
- 17:20, 17:40, 18:00, 18:20, 18:40, 19:00, 19:20, 19:40

---

### Example 2: Late Night Work

**Scenario**: Avoid showing commits between 10 PM and 6 AM

```json
{
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 22,
    "restricted_end_hour": 6,
    "suggestion_interval": 30
  }
}
```

**Suggested times** (if committing at 1 AM):
- 06:30, 07:00, 07:30, 08:00, 08:30, 09:00

**Note**: This handles overnight ranges properly.

---

### Example 3: Flexible Scheduling

**Scenario**: Fine-grained control with 15-minute intervals

```json
{
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 8,
    "restricted_end_hour": 18,
    "suggestion_interval": 15
  }
}
```

**Suggested times** (if committing at 3 PM):
- 18:15, 18:30, 18:45, 19:00, 19:15, 19:30, 19:45, 20:00

---

### Example 4: Minimal Options

**Scenario**: Clean hourly timestamps only

```json
{
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 9,
    "restricted_end_hour": 17,
    "suggestion_interval": 60
  }
}
```

**Suggested times** (if committing at 2 PM):
- 18:00, 19:00, 20:00, 21:00, 22:00, 23:00

---

### Example 5: Extended Hours

**Scenario**: Avoid commits during extended work day (7 AM - 7 PM)

```json
{
  "delayed_commit": {
    "enabled": true,
    "restricted_start_hour": 7,
    "restricted_end_hour": 19,
    "suggestion_interval": 25
  }
}
```

**Suggested times** (if committing at noon):
- 19:25, 19:50, 20:15, 20:40, 21:05, 21:30, 21:55, 22:20

---

## Related Documentation

- [Main README](README.md) - General gocommit documentation
- [Architecture Document](DELAYED_COMMIT_ARCHITECTURE.md) - Technical implementation details
- [Installation Guide](INSTALL.md) - Setup instructions
- [Configuration Guide](config/config.go) - Config file structure

---

## Support

For issues, questions, or feature requests related to the delayed commit feature:

1. Check this documentation for answers
2. Review the [FAQ section](#faq--troubleshooting)
3. Check your configuration: `cat ~/.gocommit.json`
4. Open an issue on GitHub with:
   - Your configuration (redact API key)
   - Expected behavior
   - Actual behavior
   - Steps to reproduce

---

*Last Updated: 2024-10-21*