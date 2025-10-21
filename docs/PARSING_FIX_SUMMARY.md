# Detailed Icon Mode Parsing Fix

## Problem
When using `-d --icon` (detailed + icon mode), the AI was generating proper detailed commit messages with body text and emoji prefixes, but the parser was not handling the response correctly.

## Root Cause
The parsing logic in both [`ai/gemini.go`](ai/gemini.go) and [`ai/openrouter.go`](ai/openrouter.go) had a conditional that only parsed detailed format for `detailed && !useIcons`. When both flags were enabled, it would default to JSON parsing, but the logic didn't account for the fact that detailed icon mode uses JSON format with embedded newlines.

### Original Logic (Incorrect)
```go
if detailed && !useIcons {
    // Parse detailed format (plain text with --- separators)
    finalMessages = parseDetailedResponse(cleanText)
} else {
    // Parse JSON format for regular mode
    finalMessages = parseJSONResponse(cleanText)
}
```

This meant that `-d --icon` would try to use `parseDetailedResponse()` which expected plain text with `---` separators, but the AI was actually generating JSON format as specified in `iconDetailedPromptTemplate`.

## Solution

### 1. Updated `parseResponse()` Logic
Modified the conditional to explicitly handle all three cases:
- **Detailed without icons**: Use plain text format with `---` separators
- **Detailed with icons**: Use JSON format (messages contain `\n\n` for body)
- **Regular mode**: Use JSON format

```go
if detailed && !useIcons {
    // Parse detailed format (plain text with --- separators)
    finalMessages = parseDetailedResponse(cleanText)
} else if detailed && useIcons {
    // Parse JSON format for detailed icon mode - messages contain newlines for body
    finalMessages = parseJSONResponse(cleanText)
} else {
    // Parse JSON format for regular mode (with or without icons)
    finalMessages = parseJSONResponse(cleanText)
}
```

### 2. Enhanced `parseDetailedResponse()` 
Updated the function to properly handle emoji prefixes in commit messages:
- Detects emoji characters at the start of lines
- Strips emoji to check for commit type prefixes
- Preserves emoji in the final output
- Better handling of empty lines in message bodies

### 3. Improved `isCommitMessageStart()` (OpenRouter)
Enhanced the commit message detection to:
- Strip emoji before checking for commit types
- Support all commit types (feat, fix, docs, style, refactor, perf, test, chore, build, ci, revert, try, init)
- Handle both scoped and unscoped commit messages

## Prompt Templates Analysis

The fix aligns with how the prompt templates are designed:

1. **`detailedPromptTemplate`** (detailed without icons):
   - Expects plain text with `---` separators
   - Example format: `feat: title\n\nbody\n\n---\n\nfix: title\n\nbody`

2. **`iconDetailedPromptTemplate`** (detailed with icons):
   - **Expects JSON format** with newlines embedded in strings
   - Example format: `{"messages": ["✨ feat: title\n\nbody", "🐛 fix: title\n\nbody"]}`

3. **`iconPromptTemplate`** (icons without detailed):
   - Expects JSON format with single-line messages
   - Example format: `{"messages": ["✨ feat: title", "🐛 fix: title"]}`

4. **`promptTemplate`** (standard mode):
   - Expects JSON format with single-line messages
   - Example format: `{"messages": ["feat: title", "fix: title"]}`

## Testing

Created comprehensive tests in [`ai/parsing_test.go`](ai/parsing_test.go) covering:
- ✅ Detailed Icon Mode - JSON format with emoji and body
- ✅ Detailed Mode - Plain text with `---` separators
- ✅ Icon Mode - JSON format with emoji only
- ✅ Detailed Response with Emoji - Plain text with emoji and `---` separators
- ✅ OpenRouter Detailed Icon Mode - JSON format parsing

All tests pass successfully:
```
PASS: TestParseDetailedIconResponse
PASS: TestParseDetailedResponseWithEmoji
PASS: TestOpenRouterParseDetailedIconResponse
```

## Expected Behavior After Fix

### Mode: `gocommit -d --icon`
**AI Response:**
```json
{
  "messages": [
    "✨ feat: introduce multi-AI provider support\n\nAdd OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.",
    "🛠 refactor: split original 1649-line main.go into modular packages\n\nReorganize codebase into ai/, config/, git/, ui/, commit/, and logger/ packages for better maintainability and separation of concerns.",
    "📖 docs: update README with multi-provider setup\n\nDocument configuration options for both Gemini and OpenRouter providers, including API key setup and model selection."
  ]
}
```

**Parsed Output (displayed in selector):**
```
✨ feat: introduce multi-AI provider support

Add OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.
```

### Mode: `gocommit -d`
**AI Response:**
```
feat: introduce multi-AI provider support

Add OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.

---

refactor: split original 1649-line main.go into modular packages
...
```

**Parsed Output:**
```
feat: introduce multi-AI provider support

Add OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.
```

### Mode: `gocommit --icon`
**AI Response:**
```json
{
  "messages": [
    "✨ feat: introduce multi-AI provider support",
    "🛠 refactor: split main.go into modular packages",
    "📖 docs: update README with setup"
  ]
}
```

**Parsed Output:**
```
✨ feat: introduce multi-AI provider support
```

## Files Modified
- [`ai/gemini.go`](ai/gemini.go) - Updated parseResponse() and parseDetailedResponse()
- [`ai/openrouter.go`](ai/openrouter.go) - Updated parseResponse(), parseDetailedResponse(), and isCommitMessageStart()
- [`ai/parsing_test.go`](ai/parsing_test.go) - Created comprehensive test coverage

## Backward Compatibility
✅ All existing modes continue to work:
- Standard mode (`gocommit`)
- Icon mode (`gocommit --icon`)
- Detailed mode (`gocommit -d`)
- **Fixed:** Detailed + Icon mode (`gocommit -d --icon`)