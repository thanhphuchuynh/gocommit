package ai

import (
	"strings"
	"testing"
)

func TestParseDetailedIconResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // expected number of messages
	}{
		{
			name: "Detailed Icon Mode - JSON format",
			input: `{
  "messages": [
    "✨ feat: introduce multi-AI provider support\n\nAdd OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.",
    "🛠 refactor: split original 1649-line main.go into modular packages\n\nReorganize codebase into ai/, config/, git/, ui/, commit/, and logger/ packages for better maintainability and separation of concerns.",
    "📖 docs: update README with multi-provider setup\n\nDocument configuration options for both Gemini and OpenRouter providers, including API key setup and model selection."
  ]
}`,
			expected: 3,
		},
		{
			name: "Detailed Mode - Plain text with --- separators",
			input: `feat: introduce multi-AI provider support

Add OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.

---

refactor: split original 1649-line main.go into modular packages

Reorganize codebase into ai/, config/, git/, ui/, commit/, and logger/ packages for better maintainability and separation of concerns.

---

docs: update README with multi-provider setup

Document configuration options for both Gemini and OpenRouter providers, including API key setup and model selection.`,
			expected: 3,
		},
		{
			name: "Icon Mode - JSON format",
			input: `{
  "messages": [
    "✨ feat: add new feature",
    "🐛 fix: resolve bug",
    "📖 docs: update documentation"
  ]
}`,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with GeminiProvider
			g := &GeminiProvider{}

			var messages []string
			var err error

			if strings.Contains(tt.name, "Detailed Icon") {
				messages, err = g.parseResponse(tt.input, true, true)
			} else if strings.Contains(tt.name, "Detailed Mode") {
				messages, err = g.parseResponse(tt.input, true, false)
			} else {
				messages, err = g.parseResponse(tt.input, false, true)
			}

			if err != nil {
				t.Fatalf("parseResponse failed: %v", err)
			}

			if len(messages) != tt.expected {
				t.Errorf("expected %d messages, got %d", tt.expected, len(messages))
			}

			// Verify messages are not empty
			for i, msg := range messages {
				if msg == "" {
					t.Errorf("message %d is empty", i)
				}
				t.Logf("Message %d: %s", i+1, msg)
			}
		})
	}
}

func TestParseDetailedResponseWithEmoji(t *testing.T) {
	g := &GeminiProvider{}

	// Test input with --- separators (proper detailed format for non-icon mode)
	input := `✨ feat: introduce multi‑AI provider support

– Add OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.

---

🛠 refactor: split original 1649‑line main.go into modular packages

– Reorganize codebase into ai/, config/, git/, ui/, commit/, and logger/ packages for better maintainability and separation of concerns.

---

📖 docs: update README with multi‑provider setup

– Document configuration options for both Gemini and OpenRouter providers, including API key setup and model selection.`

	messages := g.parseDetailedResponse(input)

	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
		for i, msg := range messages {
			t.Logf("Message %d:\n%s\n", i+1, msg)
		}
		return
	}

	// Verify each message has emoji and body
	for i, msg := range messages {
		t.Logf("Message %d:\n%s\n", i+1, msg)

		// Check that emoji is preserved
		runes := []rune(msg)
		if len(runes) == 0 || runes[0] <= 127 {
			t.Errorf("message %d does not start with emoji", i+1)
		}

		// Check that message has multiple lines (title + body)
		if !strings.Contains(msg, "\n") {
			t.Logf("Warning: message %d might not have body text", i+1)
		}
	}
}

func TestOpenRouterParseDetailedIconResponse(t *testing.T) {
	p := &OpenRouterProvider{}

	input := `{
  "messages": [
    "✨ feat: introduce multi-AI provider support\n\nAdd OpenRouter provider alongside Gemini, expose provider and model keys in config, and validate keys for both services.",
    "🛠 refactor: split original 1649-line main.go into modular packages\n\nReorganize codebase into ai/, config/, git/, ui/, commit/, and logger/ packages for better maintainability and separation of concerns.",
    "📖 docs: update README with multi-provider setup\n\nDocument configuration options for both Gemini and OpenRouter providers, including API key setup and model selection."
  ]
}`

	messages, err := p.parseResponse(input, true, true)
	if err != nil {
		t.Fatalf("parseResponse failed: %v", err)
	}

	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}

	// Verify messages contain emoji and body
	for i, msg := range messages {
		t.Logf("Message %d:\n%s\n", i+1, msg)

		runes := []rune(msg)
		if len(runes) == 0 || runes[0] <= 127 {
			t.Errorf("message %d does not start with emoji", i+1)
		}

		if !strings.Contains(msg, "\n\n") {
			t.Errorf("message %d does not have proper body separator", i+1)
		}
	}
}

func TestJSONSchemaResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		detailed bool
		useIcons bool
	}{
		{
			name: "Strict JSON Schema - Regular Mode",
			input: `{
  "messages": [
    "feat(api): add user authentication endpoint",
    "fix(db): resolve connection timeout issue",
    "refactor(utils): simplify error handling logic"
  ]
}`,
			detailed: false,
			useIcons: false,
		},
		{
			name: "Strict JSON Schema - Icon Mode",
			input: `{
  "messages": [
    "✨ feat(api): add user authentication endpoint",
    "🐛 fix(db): resolve connection timeout issue",
    "🛠 refactor(utils): simplify error handling logic"
  ]
}`,
			detailed: false,
			useIcons: true,
		},
		{
			name: "Strict JSON Schema - Detailed Icon Mode",
			input: `{
  "messages": [
    "✨ feat(api): add user authentication endpoint\n\nImplements JWT-based authentication with token refresh capability.\nIncludes rate limiting and security headers for protection.",
    "🐛 fix(db): resolve connection timeout issue\n\nIncreases connection pool size and adds retry logic.\nImproves error handling for transient database failures.",
    "🛠 refactor(utils): simplify error handling logic\n\nConsolidates error types into a single package.\nAdds structured logging for better debugging."
  ]
}`,
			detailed: true,
			useIcons: true,
		},
		{
			name:     "JSON with markdown code block",
			input:    "```json\n{\n  \"messages\": [\n    \"feat: add feature\",\n    \"fix: resolve bug\",\n    \"docs: update readme\"\n  ]\n}\n```",
			detailed: false,
			useIcons: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test both providers
			providers := []struct {
				name   string
				parser func(string, bool, bool) ([]string, error)
			}{
				{"Gemini", (&GeminiProvider{}).parseResponse},
				{"OpenRouter", (&OpenRouterProvider{}).parseResponse},
			}

			for _, provider := range providers {
				t.Run(provider.name, func(t *testing.T) {
					messages, err := provider.parser(tt.input, tt.detailed, tt.useIcons)
					if err != nil {
						t.Fatalf("%s parseResponse failed: %v", provider.name, err)
					}

					if len(messages) != 3 {
						t.Errorf("%s: expected 3 messages, got %d", provider.name, len(messages))
					}

					// Verify all messages are non-empty
					for i, msg := range messages {
						if msg == "" {
							t.Errorf("%s: message %d is empty", provider.name, i)
						}
						t.Logf("%s Message %d: %s", provider.name, i+1, msg)
					}
				})
			}
		})
	}
}
