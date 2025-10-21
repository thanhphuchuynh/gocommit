package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenRouterProvider implements the Provider interface for OpenRouter API
type OpenRouterProvider struct {
	apiKey string
	model  string // e.g., "anthropic/claude-3.5-sonnet", "openai/gpt-4"
	client *http.Client
}

// OpenRouterRequest represents the request format for OpenRouter API
type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

// OpenRouterMessage represents a message in the request
type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response format from OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(apiKey, model string) *OpenRouterProvider {
	if model == "" {
		model = "anthropic/claude-3.5-sonnet" // default model
	}
	return &OpenRouterProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

// GenerateMessages generates commit messages using OpenRouter API
func (p *OpenRouterProvider) GenerateMessages(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// Get the appropriate prompt template
	promptTemplate := GetPromptTemplate(req.Detailed, req.UseIcons)
	prompt := fmt.Sprintf(promptTemplate, req.Diff, req.LastCommit)

	// Prepare the request payload
	requestBody := OpenRouterRequest{
		Model: p.model,
		Messages: []OpenRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("HTTP-Referer", "https://github.com/thanhphuchuynh/gocommit")

	// Send request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for API errors
	if openRouterResp.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %s)", openRouterResp.Error.Message, openRouterResp.Error.Code)
	}

	// Extract content from response
	if len(openRouterResp.Choices) == 0 || openRouterResp.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("no content generated")
	}

	content := openRouterResp.Choices[0].Message.Content

	// Parse the response based on mode
	messages, err := p.parseResponse(content, req.Detailed, req.UseIcons)
	if err != nil {
		return nil, err
	}

	return &GenerateResponse{
		Messages: messages,
		Provider: p.Name(),
	}, nil
}

// Name returns the provider name
func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

// ValidateConfig validates the OpenRouter configuration
func (p *OpenRouterProvider) ValidateConfig() error {
	if p.apiKey == "" {
		return fmt.Errorf("API key is required")
	}
	if p.model == "" {
		return fmt.Errorf("model is required")
	}
	// OpenRouter API keys typically start with "sk-or-"
	if !strings.HasPrefix(p.apiKey, "sk-or-") && !strings.HasPrefix(p.apiKey, "sk-") {
		return fmt.Errorf("invalid API key format: OpenRouter API keys typically start with 'sk-or-' or 'sk-'")
	}
	return nil
}

// parseResponse parses the AI response into commit messages
// This method follows the same logic as GeminiProvider for consistency
func (p *OpenRouterProvider) parseResponse(text string, detailed, useIcons bool) ([]string, error) {
	cleanText := strings.TrimSpace(text)
	var finalMessages []string

	if detailed && !useIcons {
		// Parse detailed format (plain text with --- separators) for non-icon mode only
		finalMessages = p.parseDetailedResponse(cleanText)
	} else {
		// Parse JSON format for regular mode
		finalMessages = p.parseJSONResponse(cleanText)
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
func (p *OpenRouterProvider) parseDetailedResponse(cleanText string) []string {
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

			// Check if line starts with a commit type (with or without emoji)
			if p.isCommitMessageStart(line) {
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

// isCommitMessageStart checks if a line starts with a commit type
func (p *OpenRouterProvider) isCommitMessageStart(line string) bool {
	// Check for emoji prefixes first
	emojis := []string{"✨", "🐛", "📖", "💄", "🛠", "⚡️", "✅", "📦", "⚙️", "🚀", "🗑", "🤞", "🎉"}
	for _, emoji := range emojis {
		if strings.HasPrefix(line, emoji) {
			return true
		}
	}

	// Check for standard commit types
	types := []string{"feat:", "fix:", "docs:", "style:", "refactor:", "perf:", "test:", "chore:",
		"feat(", "fix(", "docs(", "style(", "refactor(", "perf(", "test(", "chore("}
	for _, t := range types {
		if strings.HasPrefix(line, t) {
			return true
		}
	}

	return false
}

// parseJSONResponse parses JSON format responses
func (p *OpenRouterProvider) parseJSONResponse(cleanText string) []string {
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
