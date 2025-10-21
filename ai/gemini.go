package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// CommitResponse represents the JSON response from AI
type CommitResponse struct {
	Messages []string `json:"messages"`
}

// GeminiProvider implements the Provider interface for Google's Gemini AI
type GeminiProvider struct {
	apiKey string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
	}
}

// Name returns the provider name
func (g *GeminiProvider) Name() string {
	return "gemini"
}

// ValidateConfig validates the provider configuration
func (g *GeminiProvider) ValidateConfig() error {
	if g.apiKey == "" {
		return fmt.Errorf("API key is required")
	}
	// Google AI API keys typically start with "AIza" and are 39 characters long
	if len(g.apiKey) != 39 || !strings.HasPrefix(g.apiKey, "AIza") {
		return fmt.Errorf("invalid API key format: Google AI API keys should start with 'AIza' and be 39 characters long")
	}
	return nil
}

// GenerateMessages generates commit messages using Gemini AI
func (g *GeminiProvider) GenerateMessages(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	// Get the appropriate prompt template
	promptTemplate := GetPromptTemplate(req.Detailed, req.UseIcons)
	prompt := fmt.Sprintf(promptTemplate, req.Diff, req.LastCommit)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	// Get the text content from the response
	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if str, ok := part.(genai.Text); ok {
			text += string(str)
		}
	}

	// Parse the response based on mode
	messages, err := g.parseResponse(text, req.Detailed, req.UseIcons)
	if err != nil {
		return nil, err
	}

	return &GenerateResponse{
		Messages: messages,
		Provider: g.Name(),
	}, nil
}

// parseResponse parses the AI response into commit messages
func (g *GeminiProvider) parseResponse(text string, detailed, useIcons bool) ([]string, error) {
	cleanText := strings.TrimSpace(text)
	var finalMessages []string

	if detailed && !useIcons {
		// Parse detailed format (plain text with --- separators) for non-icon mode only
		finalMessages = g.parseDetailedResponse(cleanText)
	} else {
		// Parse JSON format for regular mode
		finalMessages = g.parseJSONResponse(cleanText)
	}

	if len(finalMessages) == 0 {
		return nil, fmt.Errorf("no valid commit messages found in response: %s", cleanText)
	}

	// Ensure we have at least 3 messages
	for len(finalMessages) < 3 {
		// If we don't have enough messages, duplicate the last one with slight variation
		lastMsg := finalMessages[len(finalMessages)-1]
		finalMessages = append(finalMessages, lastMsg)
	}

	// Take only the first 3 messages
	if len(finalMessages) > 3 {
		finalMessages = finalMessages[:3]
	}

	return finalMessages, nil
}

// parseDetailedResponse parses detailed format responses (plain text with --- separators)
func (g *GeminiProvider) parseDetailedResponse(cleanText string) []string {
	var finalMessages []string

	// Split by --- separators and extract commit messages
	parts := strings.Split(cleanText, "---")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			finalMessages = append(finalMessages, part)
		}
	}

	// If we don't have 3 messages from --- separator, try alternative parsing
	if len(finalMessages) < 3 {
		// Look for commit messages that start with type prefixes
		lines := strings.Split(cleanText, "\n")
		var currentMessage strings.Builder
		messages := []string{}

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Check if line starts with a commit type
			if strings.HasPrefix(line, "feat:") || strings.HasPrefix(line, "fix:") ||
				strings.HasPrefix(line, "docs:") || strings.HasPrefix(line, "style:") ||
				strings.HasPrefix(line, "refactor:") || strings.HasPrefix(line, "perf:") ||
				strings.HasPrefix(line, "test:") || strings.HasPrefix(line, "chore:") ||
				strings.HasPrefix(line, "feat(") || strings.HasPrefix(line, "fix(") ||
				strings.HasPrefix(line, "docs(") || strings.HasPrefix(line, "style(") ||
				strings.HasPrefix(line, "refactor(") || strings.HasPrefix(line, "perf(") ||
				strings.HasPrefix(line, "test(") || strings.HasPrefix(line, "chore(") {

				// If we have a previous message, save it
				if currentMessage.Len() > 0 {
					messages = append(messages, strings.TrimSpace(currentMessage.String()))
					currentMessage.Reset()
				}
				// Start new message
				currentMessage.WriteString(line)
			} else if currentMessage.Len() > 0 {
				// Add to current message body
				currentMessage.WriteString("\n" + line)
			}
		}

		// Add the last message if any
		if currentMessage.Len() > 0 {
			messages = append(messages, strings.TrimSpace(currentMessage.String()))
		}

		if len(messages) >= 3 {
			finalMessages = messages[:3]
		} else if len(messages) > 0 {
			finalMessages = messages
		}
	}

	return finalMessages
}

// parseJSONResponse parses JSON format responses
func (g *GeminiProvider) parseJSONResponse(cleanText string) []string {
	// Try to extract JSON from markdown code blocks
	if strings.Contains(cleanText, "```json") {
		// Extract JSON from json code block
		start := strings.Index(cleanText, "```json")
		if start != -1 {
			start += 7 // Skip "```json"
			// Skip any whitespace/newlines after ```json
			for start < len(cleanText) && (cleanText[start] == '\n' || cleanText[start] == '\r' || cleanText[start] == ' ' || cleanText[start] == '\t') {
				start++
			}
			end := strings.Index(cleanText[start:], "```")
			if end != -1 {
				cleanText = strings.TrimSpace(cleanText[start : start+end])
			}
		}
	} else if strings.Contains(cleanText, "```") {
		// Extract JSON from generic code block
		start := strings.Index(cleanText, "```")
		if start != -1 {
			start += 3 // Skip "```"
			// Skip any whitespace/newlines after ```
			for start < len(cleanText) && (cleanText[start] == '\n' || cleanText[start] == '\r' || cleanText[start] == ' ' || cleanText[start] == '\t') {
				start++
			}
			end := strings.Index(cleanText[start:], "```")
			if end != -1 {
				cleanText = strings.TrimSpace(cleanText[start : start+end])
			}
		}
	}

	// Find JSON object in the response - look for complete JSON structure
	jsonStart := strings.Index(cleanText, "{")
	if jsonStart == -1 {
		return nil
	}

	// Find the matching closing brace for the JSON object
	jsonText := ""
	braceCount := 0
	inQuotes := false
	escapeNext := false

	for i := jsonStart; i < len(cleanText); i++ {
		char := cleanText[i]

		if escapeNext {
			escapeNext = false
			continue
		}

		if char == '\\' && inQuotes {
			escapeNext = true
			continue
		}

		if char == '"' && !escapeNext {
			inQuotes = !inQuotes
		}

		if !inQuotes {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					jsonText = cleanText[jsonStart : i+1]
					break
				}
			}
		}
	}

	if jsonText == "" {
		return nil
	}

	// Parse JSON response
	var commitResponse CommitResponse
	if err := json.Unmarshal([]byte(jsonText), &commitResponse); err != nil {
		return nil
	}

	if len(commitResponse.Messages) == 0 {
		return nil
	}

	return commitResponse.Messages
}
