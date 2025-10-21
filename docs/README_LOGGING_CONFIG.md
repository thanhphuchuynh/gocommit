# Logging Configuration

GoCommit supports configurable request logging to help with debugging and analysis. By default, logging is **disabled**.

## Configuration

### Enable Logging
```bash
gocommit --enable-logging
```

### Disable Logging  
```bash
gocommit --disable-logging
```

### Check Current Status
The logging configuration is stored in `~/.gocommit.json`:
```bash
cat ~/.gocommit.json
```

Example output:
```json
{
  "api_key": "AIzaSyCUNHL0gmVWKUDGQpySBCw0GhnJr3rrklw",
  "logging_enabled": true
}
```

## Log File Location

When logging is enabled, requests are logged to:
```
~/.gocommit/gocommit_requests.log
```

## What Gets Logged

When logging is enabled, the following information is captured for each request:

### Successful Requests
- Timestamp
- Git diff
- Last commit message (for context)
- Prompt sent to AI
- AI response
- Generated commit messages
- Selected commit message

### Failed Requests  
- Timestamp
- Git diff
- Last commit message (for context)
- Prompt sent to AI
- Error message

## Privacy and Security

- Logging is **disabled by default** to protect privacy
- When disabled, no requests are logged to disk
- Log files are stored locally on your machine only
- Consider the sensitivity of your git diffs before enabling logging
- You can delete log files manually if needed: `rm ~/.gocommit/gocommit_requests.log`

## Usage Examples

1. **Enable logging for debugging:**
   ```bash
   gocommit --enable-logging
   # Use gocommit normally
   # Check logs if needed
   ```

2. **Disable logging for privacy:**
   ```bash
   gocommit --disable-logging
   # Logging is now turned off
   ```

3. **Temporary logging session:**
   ```bash
   gocommit --enable-logging
   # Make some commits with logging
   gocommit --disable-logging
   # Back to no logging
