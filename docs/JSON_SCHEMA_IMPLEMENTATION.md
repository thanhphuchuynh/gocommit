# JSON Schema Implementation for AI Providers

## Overview

This document describes the implementation of structured JSON output using LangChain-style approaches for both Gemini and OpenRouter AI providers in the GoCommit project.

## Implementation Details

### 1. Gemini Provider (`ai/gemini.go`)

**Approach**: Temperature-based consistency
- Set model temperature to `0.3` for more deterministic and consistent JSON responses
- The Gemini SDK version (0.5.0) doesn't support `ResponseMIMEType` directly
- Rely on well-structured prompts with clear JSON schema examples

```go
temperature := float32(0.3) // Lower temperature for more consistent JSON
model.SetTemperature(temperature)
```

### 2. OpenRouter Provider (`ai/openrouter.go`)

**Approach**: Structured output with JSON Schema validation + Fallback parsing

⚠️ **IMPORTANT - Model Compatibility**: The `response_format` parameter with strict JSON schema is **ONLY supported by specific models**:
- ✅ OpenAI GPT-4o and newer (e.g., `openai/gpt-4o`, `openai/chatgpt-4o-latest`)
- ✅ Fireworks-provided models
- ❌ Most other models (Gemini, Claude, Llama, etc.) do NOT support `response_format`

**Implementation Strategy**:
- For supported models: Uses `json_schema` type with strict validation
- For unsupported models: Relies on enhanced prompting + fallback text parser
- Automatically detects model capabilities and adjusts approach
- Schema enforces exactly 3 commit messages as an array

```go
type OpenRouterResponseFmt struct {
    Type       string                 `json:"type"`
    JSONSchema *OpenRouterJSONSchema  `json:"json_schema,omitempty"`
}

type OpenRouterJSONSchema struct {
    Name   string         `json:"name"`
    Schema map[string]any `json:"schema"`
    Strict bool           `json:"strict"`
}
```

**Schema Definition**:
```json
{
  "type": "object",
  "properties": {
    "messages": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 3,
      "maxItems": 3
    }
  },
  "required": ["messages"],
  "additionalProperties": false
}
```

## Benefits

### 1. **Consistent JSON Output**
- Both providers now return more reliable JSON responses
- Reduced parsing errors and edge cases
- Strict schema validation for supported OpenRouter models ensures exact format
- Fallback text parser handles non-compliant responses gracefully

### 2. **Better Error Handling**
- Structured output reduces the need for complex parsing logic
- JSON schema validation catches malformed responses early (when supported)
- Lower temperature for Gemini reduces hallucinations
- Multiple parsing strategies ensure robustness

### 3. **LangChain-Style Approach**
- Follows industry best practices for structured LLM outputs
- Similar to LangChain's `StructuredOutputParser`
- Type-safe output format with clear expectations
- Graceful degradation for models without structured output support

### 4. **Mode-Specific Handling**
- JSON schema is applied for regular and icon modes (when supported)
- Detailed non-icon mode uses plain text parsing for better formatting
- Flexible approach supports different output formats
- Model-aware configuration adapts to capabilities

## Response Formats by Mode

| Mode | Use Icons | Output Format | Schema Applied |
|------|-----------|---------------|----------------|
| Regular | No | JSON | Yes |
| Regular | Yes | JSON | Yes |
| Detailed | No | Plain text (---) | No |
| Detailed | Yes | JSON | Yes |

## Testing

Comprehensive tests verify:
- JSON parsing with and without markdown code blocks
- Emoji handling in commit messages
- Detailed vs regular mode responses
- Both Gemini and OpenRouter providers
- Edge cases and malformed inputs

Run tests:
```bash
go test -v ./ai/...
```

## Example Response

**Input Prompt**: Generate 3 commit messages for a git diff

**Expected JSON Output**:
```json
{
  "messages": [
    "feat(api): add user authentication endpoint",
    "fix(db): resolve connection timeout issue",
    "refactor(utils): simplify error handling logic"
  ]
}
```

**With Icons**:
```json
{
  "messages": [
    "✨ feat(api): add user authentication endpoint",
    "🐛 fix(db): resolve connection timeout issue",
    "🛠 refactor(utils): simplify error handling logic"
  ]
}
```

**With Detailed Bodies**:
```json
{
  "messages": [
    "✨ feat(api): add user authentication endpoint\n\nImplements JWT-based authentication with token refresh capability.\nIncludes rate limiting and security headers for protection.",
    "🐛 fix(db): resolve connection timeout issue\n\nIncreases connection pool size and adds retry logic.\nImproves error handling for transient database failures.",
    "🛠 refactor(utils): simplify error handling logic\n\nConsolidates error types into a single package.\nAdds structured logging for better debugging."
  ]
}
```

## Model Recommendations

### For Best Results with Structured Outputs

**Via OpenRouter**:
- ✅ `openai/gpt-4o` - Full JSON schema support, highly reliable
- ✅ `openai/chatgpt-4o-latest` - Latest GPT-4o with improvements
- ⚠️ `anthropic/claude-3.5-sonnet` - Good prompt following, no strict schema
- ⚠️ `google/gemini-2.0-flash-exp:free` - Free but inconsistent JSON formatting
- ⚠️ `meta-llama/llama-3.1-70b-instruct` - Good performance, relies on prompting

**Direct Provider**:
- ✅ Gemini via Google AI API - Works with temperature tuning + strong prompts

### Troubleshooting

**If you see "no valid commit messages found" errors**:
1. Check if your model supports `response_format` (see list above)
2. Consider switching to `openai/gpt-4o` for guaranteed JSON compliance
3. The fallback parser should handle most cases, but some models may still struggle
4. Check logs to see the actual AI response format

## Future Improvements

1. **Upgrade Gemini SDK**: When a newer version with `ResponseMIMEType` support is available, implement proper JSON mode
2. **Schema Evolution**: Add more complex schemas for additional metadata (author, timestamp, tags)
3. **Validation Layer**: Add runtime JSON schema validation before parsing
4. **Retry Logic**: Implement automatic retries with schema validation on parsing failures
5. **Model Detection**: Auto-detect model capabilities via OpenRouter API metadata

## References

- [OpenRouter Structured Outputs](https://openrouter.ai/docs/structured-outputs)
- [Gemini JSON Mode Documentation](https://ai.google.dev/gemini-api/docs/json-mode)
- [LangChain Structured Output Parser](https://python.langchain.com/docs/modules/model_io/output_parsers/structured)
- [OpenAI JSON Schema](https://platform.openai.com/docs/guides/structured-outputs)
