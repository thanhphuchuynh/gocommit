package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/nsf/termbox-go"
	"github.com/thanhphuchuynh/config"
	"github.com/thanhphuchuynh/logger"
	"google.golang.org/api/option"
)

// CommitResponse represents the JSON response from AI
type CommitResponse struct {
	Messages []string `json:"messages"`
}

const promptTemplate = `You are an expert at writing conventional commit messages. Analyze the git diff and generate 3 high-quality commit messages.

## Conventional Commit Format
type(optional-scope): description

## Available Types
- feat: new feature for users
- fix: bug fix for users
- docs: documentation changes
- style: formatting, missing semicolons, etc (no code change)
- refactor: code change that neither fixes bug nor adds feature
- perf: code change that improves performance
- test: adding/updating tests
- chore: updating build tasks, package manager configs, etc

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end
4. Scope should reflect the actual changed component
5. Description should explain WHAT changed, not WHY
6. Be specific about the actual changes in the diff

## Examples
- feat(auth): add OAuth2 login support
- fix(api): handle null response in user endpoint
- refactor(utils): extract validation logic to separate module
- docs(readme): update installation instructions
- style(components): fix indentation in header component

Generate exactly 3 different commit messages in JSON format:

{
  "messages": [
    "feat(scope): commit message 1",
    "fix(scope): commit message 2",
    "refactor(scope): commit message 3"
  ]
}

Return only valid JSON with no additional text.`

const detailedPromptTemplate = `You are an expert at writing conventional commit messages. Analyze the git diff and generate 3 high-quality, detailed commit messages.

## Conventional Commit Format
type(optional-scope): description

Detailed description explaining what was changed and why.

## Available Types
- feat: new feature for users
- fix: bug fix for users
- docs: documentation changes
- style: formatting, missing semicolons, etc (no code change)
- refactor: code change that neither fixes bug nor adds feature
- perf: code change that improves performance
- test: adding/updating tests
- chore: updating build tasks, package manager configs, etc

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Subject line: max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end of subject line
4. Scope should reflect the actual changed component
5. Subject should explain WHAT changed
6. Body should explain WHY and HOW in more detail
7. Be specific about the actual changes in the diff
8. Include technical details and context in the body

## Examples
feat(auth): add OAuth2 login support

Implements OAuth2 authentication flow with Google and GitHub providers.
Adds token validation, refresh mechanism, and user profile fetching.
Includes comprehensive error handling for failed authentication attempts.

fix(api): handle null response in user endpoint

Prevents application crash when user data is missing from database.
Adds null checks and default values for required user fields.
Improves error messaging for better debugging experience.

Generate exactly 3 different detailed commit messages with body text. Format them as:

feat: enhance commit message generation with JSON output

Refactors the commit message generation process to return responses in JSON format.
This change ensures a structured and parsable output, improving integration with other tools.
Updates prompt templates to explicitly request JSON formatted messages and removes parsing logic.

---

fix: correct JSON parsing in commit message generation

Addresses an issue where the JSON output from the AI model was not correctly parsed.
Improves JSON extraction from the response by handling potential code blocks.
Adds more robust error handling and logging for debugging JSON parsing failures.

---

chore: update dependencies and improve error handling

Updates the go.mod and go.sum files to include the latest dependencies.
Improves error handling throughout the application, providing more informative error messages.
Includes changes to gracefully handle API request failures.

Return only the 3 commit messages in the format shown above with no additional text.`

const iconPromptTemplate = `You are an expert at writing conventional commit messages with emoji icons. Analyze the git diff and generate 3 high-quality commit messages using emoji icons.

## Conventional Commit Format with Icons
emoji type(optional-scope): description

## Available Types with Icons
- ✨ feat: new feature for users
- 🐛 fix: bug fix for users
- 📖 docs: documentation changes
- 💄 style: formatting, missing semicolons, etc (no code change)
- 🛠 refactor: code change that neither fixes bug nor adds feature
- ⚡️ perf: code change that improves performance
- ✅ test: adding/updating tests
- 📦 build: changes that affect the build system or external dependencies
- ⚙️ ci: changes to CI configuration files and scripts
- 🚀 chore: other changes that don't modify src or test files
- 🗑 revert: reverts a previous commit
- 🤞 try: add untested to production
- 🎉 init: project init

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end
4. Scope should reflect the actual changed component
5. Description should explain WHAT changed, not WHY
6. Be specific about the actual changes in the diff
7. Always start with the appropriate emoji icon

## Examples
- ✨ feat(auth): add OAuth2 login support
- 🐛 fix(api): handle null response in user endpoint
- 🛠 refactor(utils): extract validation logic to separate module
- 📖 docs(readme): update installation instructions
- 💄 style(components): fix indentation in header component

Generate exactly 3 different commit messages in JSON format:

{
  "messages": [
    "✨ feat(scope): commit message 1",
    "🐛 fix(scope): commit message 2",
    "🛠 refactor(scope): commit message 3"
  ]
}

Return only valid JSON with no additional text.`

const iconDetailedPromptTemplate = `You are an expert at writing conventional commit messages with emoji icons. Analyze the git diff and generate 3 high-quality, detailed commit messages using emoji icons.

## Conventional Commit Format with Icons
emoji type(optional-scope): description

Detailed description explaining what was changed and why.

## Available Types with Icons
- ✨ feat: new feature for users
- 🐛 fix: bug fix for users
- 📖 docs: documentation changes
- 💄 style: formatting, missing semicolons, etc (no code change)
- 🛠 refactor: code change that neither fixes bug nor adds feature
- ⚡️ perf: code change that improves performance
- ✅ test: adding/updating tests
- 📦 build: changes that affect the build system or external dependencies
- ⚙️ ci: changes to CI configuration files and scripts
- 🚀 chore: other changes that don't modify src or test files
- 🗑 revert: reverts a previous commit
- 🤞 try: add untested to production
- 🎉 init: project init

## Scope Guidelines
- Use specific component/module names when applicable
- Examples: api, ui, auth, config, database, utils
- Omit scope if change affects multiple areas

## Git Diff
%s

## Last Commit (for style reference)
%s

## Requirements
1. Subject line: max 200 characters (industry standard)
2. Use imperative mood ("add" not "added" or "adds")
3. No period at the end of subject line
4. Scope should reflect the actual changed component
5. Subject should explain WHAT changed
6. Body should explain WHY and HOW in more detail
7. Be specific about the actual changes in the diff
8. Include technical details and context in the body
9. Always start with the appropriate emoji icon

## Examples
✨ feat(auth): add OAuth2 login support

Implements OAuth2 authentication flow with Google and GitHub providers.
Adds token validation, refresh mechanism, and user profile fetching.
Includes comprehensive error handling for failed authentication attempts.

🐛 fix(api): handle null response in user endpoint

Prevents application crash when user data is missing from database.
Adds null checks and default values for required user fields.
Improves error messaging for better debugging experience.

Generate exactly 3 different detailed commit messages with body text in JSON format:

{
  "messages": [
    "✨ feat: enhance commit message generation with JSON output\n\nRefactors the commit message generation process to return responses in JSON format.\nThis change ensures a structured and parsable output, improving integration with other tools.\nUpdates prompt templates to explicitly request JSON formatted messages and removes parsing logic.",
    "🐛 fix: correct JSON parsing in commit message generation\n\nAddresses an issue where the JSON output from the AI model was not correctly parsed.\nImproves JSON extraction from the response by handling potential code blocks.\nAdds more robust error handling and logging for debugging JSON parsing failures.",
    "🚀 chore: update dependencies and improve error handling\n\nUpdates the go.mod and go.sum files to include the latest dependencies.\nImproves error handling throughout the application, providing more informative error messages.\nIncludes changes to gracefully handle API request failures."
  ]
}

Return only valid JSON with no additional text.`

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

func generateCommitMessages(diff string, apiKey string, detailed bool, iconMode bool) ([]string, string, string, string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, "", "", "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Get last commit message for reference
	lastCommitMsg, err := getLastCommitMessage()
	if err != nil {
		log.Printf("Warning: Could not get last commit message: %v", err)
		lastCommitMsg = ""
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	var prompt string
	if iconMode {
		if detailed {
			prompt = fmt.Sprintf(iconDetailedPromptTemplate, diff, lastCommitMsg)
		} else {
			prompt = fmt.Sprintf(iconPromptTemplate, diff, lastCommitMsg)
		}
	} else {
		if detailed {
			prompt = fmt.Sprintf(detailedPromptTemplate, diff, lastCommitMsg)
		} else {
			prompt = fmt.Sprintf(promptTemplate, diff, lastCommitMsg)
		}
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, diff, lastCommitMsg, prompt, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, diff, lastCommitMsg, prompt, fmt.Errorf("no content generated")
	}

	// Get the text content from the response
	text := ""
	// fmt.Printf("Generated response: %+v \n", resp.Candidates[0].Content.Parts)

	for _, part := range resp.Candidates[0].Content.Parts {
		if str, ok := part.(genai.Text); ok {
			text += string(str)
		}
	}

	// Clean the response text - remove markdown code blocks if present
	cleanText := strings.TrimSpace(text)

	// Log the raw response for debugging
	// log.Printf("Raw AI response: %s", text)
	var finalMessages []string

	if detailed && !iconMode {
		// Parse detailed format (plain text with --- separators) for non-icon mode only
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
			} else {
				finalMessages = messages
			}
		}

		if len(finalMessages) == 0 {
			return nil, diff, lastCommitMsg, prompt, fmt.Errorf("no valid commit messages found in detailed response: %s", cleanText)
		}
	} else {
		// Parse JSON format for regular mode
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
			// Try to find JSON in the original text if cleaning failed
			jsonStart = strings.Index(text, "{")
			if jsonStart == -1 {
				return nil, diff, lastCommitMsg, prompt, fmt.Errorf("no valid JSON found in response: %s", text)
			}
			cleanText = text
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
			return nil, diff, lastCommitMsg, prompt, fmt.Errorf("no complete JSON object found in response: %s", cleanText)
		}

		// log.Printf("Extracted JSON: %s", jsonText)

		// Parse JSON response
		var commitResponse CommitResponse
		if err := json.Unmarshal([]byte(jsonText), &commitResponse); err != nil {
			return nil, diff, lastCommitMsg, prompt, fmt.Errorf("failed to parse JSON response: %v\nJSON: %s", err, jsonText)
		}

		if len(commitResponse.Messages) < 3 {
			return nil, diff, lastCommitMsg, prompt, fmt.Errorf("expected 3 messages, got %d", len(commitResponse.Messages))
		}

		// Take only the first 3 messages
		finalMessages = make([]string, 3)
		copy(finalMessages, commitResponse.Messages[:3])
	}

	// Ensure we have at least 3 messages
	for len(finalMessages) < 3 {
		// If we don't have enough messages, duplicate the last one with slight variation
		lastMsg := finalMessages[len(finalMessages)-1]
		finalMessages = append(finalMessages, lastMsg)
	}

	return finalMessages, diff, lastCommitMsg, prompt, nil
}

func getLastCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting last commit message: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func drawMessages(messages []string, selected int, showEditPrompt bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Get terminal width
	width, _ := termbox.Size()
	if width < 10 {
		width = 80 // Default width if terminal is too small
	}

	// Draw title
	title := "Select a commit message:"
	for i, ch := range title {
		termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
	}

	// Draw messages
	for i, msg := range messages {
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault
		if i == selected {
			fg = termbox.ColorBlack
			bg = termbox.ColorGreen
		}
		// Add bullet point
		bullet := "•"
		if i == selected {
			bullet = "→"
		}
		termbox.SetCell(2, i+2, []rune(bullet)[0], fg, bg)

		// For detailed messages, show only the subject line (first line)
		displayMsg := msg
		if strings.Contains(msg, "\n") {
			// Extract just the subject line for display
			lines := strings.Split(msg, "\n")
			displayMsg = lines[0]
		}

		// Truncate message if it's too long for the terminal
		maxMsgWidth := width - 6 // Account for bullet point and spacing
		if len(displayMsg) > maxMsgWidth {
			displayMsg = displayMsg[:maxMsgWidth-3] + "..."
		}

		// Draw message
		for j, ch := range displayMsg {
			termbox.SetCell(j+4, i+2, ch, fg, bg)
		}
	}

	// Draw custom message option
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if selected == len(messages) {
		fg = termbox.ColorBlack
		bg = termbox.ColorGreen
	}
	customMsg := "Edit custom message"
	bullet := "•"
	if selected == len(messages) {
		bullet = "→"
	}
	termbox.SetCell(2, len(messages)+2, []rune(bullet)[0], fg, bg)
	for j, ch := range customMsg {
		termbox.SetCell(j+4, len(messages)+2, ch, fg, bg)
	}

	// Draw instructions
	instructions := "↑↓: Move  Enter: Select  Esc: Cancel"
	for i, ch := range instructions {
		termbox.SetCell(i, len(messages)+4, ch, termbox.ColorCyan, termbox.ColorDefault)
	}

	// Draw edit prompt if needed
	if showEditPrompt {
		prompt := "Edit message (press Enter to confirm):"
		for i, ch := range prompt {
			termbox.SetCell(i, len(messages)+6, ch, termbox.ColorYellow, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func editMessage(initialMsg string) (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	// Get terminal width
	width, _ := termbox.Size()
	if width < 10 {
		width = 80 // Default width if terminal is too small
	}

	// Clear screen and show edit prompt
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw title
	title := "Edit commit message (Enter to confirm, Shift+Enter for new line, Esc to cancel):"
	for i, ch := range title {
		termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
	}

	// Draw instructions
	instructions := "Use arrow keys to move, Shift+Enter for new line, backspace/delete to edit"
	for i, ch := range instructions {
		termbox.SetCell(i, 1, ch, termbox.ColorCyan, termbox.ColorDefault)
	}

	// Initialize message buffer and cursor
	editedMsg := []rune(initialMsg)
	cursorPos := len(editedMsg)
	scrollX := 0
	maxScroll := 0
	currentLine := 0

	// Function to redraw the message with wrapping and scrolling
	redraw := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Get terminal height for scrolling
		_, height := termbox.Size()
		if height < 10 {
			height = 24 // Default height
		}

		// Draw title and instructions
		for i, ch := range title {
			termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
		}
		for i, ch := range instructions {
			termbox.SetCell(i, 1, ch, termbox.ColorCyan, termbox.ColorDefault)
		}

		// Split message into display lines (considering word wrapping)
		msgLines := strings.Split(string(editedMsg), "\n")
		displayLines := []string{}
		lineToOriginal := []int{} // Maps display line to original line number

		for lineIdx, msgLine := range msgLines {
			if len(msgLine) <= width-1 {
				// Line fits in one display line
				displayLines = append(displayLines, msgLine)
				lineToOriginal = append(lineToOriginal, lineIdx)
			} else {
				// Wrap long lines
				for i := 0; i < len(msgLine); i += width - 1 {
					end := i + width - 1
					if end > len(msgLine) {
						end = len(msgLine)
					}
					displayLines = append(displayLines, msgLine[i:end])
					lineToOriginal = append(lineToOriginal, lineIdx)
				}
			}
		}

		// Calculate cursor position in display coordinates
		linesBeforeCursor := strings.Split(string(editedMsg[:cursorPos]), "\n")
		cursorOriginalLine := len(linesBeforeCursor) - 1
		lastLine := linesBeforeCursor[len(linesBeforeCursor)-1]
		cursorColInOriginal := len(lastLine)

		// Find cursor position in display coordinates
		cursorDisplayLine := 0
		cursorDisplayCol := 0
		for i, origLine := range lineToOriginal {
			if origLine == cursorOriginalLine {
				lineStart := 0
				for j := 0; j < i; j++ {
					if lineToOriginal[j] == cursorOriginalLine {
						lineStart += len(displayLines[j])
					}
				}
				if cursorColInOriginal >= lineStart && cursorColInOriginal <= lineStart+len(displayLines[i]) {
					cursorDisplayLine = i
					cursorDisplayCol = cursorColInOriginal - lineStart
					break
				}
			}
		}

		// Calculate scroll offset to keep cursor visible
		maxDisplayLines := height - 3 // Account for title and instructions
		scrollY := 0
		if cursorDisplayLine >= maxDisplayLines {
			scrollY = cursorDisplayLine - maxDisplayLines + 1
		}

		// Draw visible lines
		for i := scrollY; i < len(displayLines) && i-scrollY < maxDisplayLines; i++ {
			line := displayLines[i]
			y := 2 + (i - scrollY)

			// Draw each character of the line
			for j, ch := range line {
				if j < width-1 {
					termbox.SetCell(j, y, ch, termbox.ColorDefault, termbox.ColorDefault)
				}
			}
		}

		// Set cursor position (adjusted for scroll)
		if cursorDisplayLine >= scrollY && cursorDisplayLine < scrollY+maxDisplayLines {
			cursorY := 2 + (cursorDisplayLine - scrollY)
			if cursorDisplayCol < width-1 {
				termbox.SetCursor(cursorDisplayCol, cursorY)
			}
		}

		termbox.Flush()
	}

	redraw()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				if ev.Mod&termbox.ModAlt != 0 {
					// Insert newline when Shift+Enter is pressed
					editedMsg = append(editedMsg[:cursorPos], append([]rune{'\n'}, editedMsg[cursorPos:]...)...)
					cursorPos++
					redraw()
				} else {
					// Regular Enter confirms the edit
					return string(editedMsg), nil
				}
			case termbox.KeyEsc:
				return "", fmt.Errorf("edit cancelled")
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorPos > 0 {
					editedMsg = append(editedMsg[:cursorPos-1], editedMsg[cursorPos:]...)
					cursorPos--
					if cursorPos < scrollX {
						scrollX = cursorPos
					}
				}
			case termbox.KeyDelete:
				if cursorPos < len(editedMsg) {
					editedMsg = append(editedMsg[:cursorPos], editedMsg[cursorPos+1:]...)
				}
			case termbox.KeyArrowLeft:
				if cursorPos > 0 {
					cursorPos--
					if cursorPos < scrollX {
						scrollX = cursorPos
					}
				}
			case termbox.KeyArrowRight:
				if cursorPos < len(editedMsg) {
					cursorPos++
					if cursorPos >= scrollX+width-1 {
						scrollX = cursorPos - width + 2
					}
				}
			case termbox.KeyArrowUp:
				// Move cursor up one line
				lines := strings.Split(string(editedMsg[:cursorPos]), "\n")
				if len(lines) > 1 {
					currentLine = len(lines) - 2
					prevLine := lines[currentLine]
					if len(prevLine) < cursorPos-len(lines[len(lines)-1])-1 {
						cursorPos = len(strings.Join(lines[:currentLine+1], "\n")) + 1
					} else {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1 + len(prevLine)
					}
					redraw()
				}
			case termbox.KeyArrowDown:
				// Move cursor down one line
				lines := strings.Split(string(editedMsg[:cursorPos]), "\n")
				if currentLine < len(lines)-1 {
					currentLine++
					nextLine := lines[currentLine]
					if len(nextLine) < cursorPos-len(strings.Join(lines[:currentLine], "\n"))-1 {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1 + len(nextLine)
					} else {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1
					}
					redraw()
				}
			case termbox.KeySpace:
				editedMsg = append(editedMsg[:cursorPos], append([]rune{' '}, editedMsg[cursorPos:]...)...)
				cursorPos++
				if cursorPos >= scrollX+width-1 {
					scrollX = cursorPos - width + 2
				}
			default:
				if ev.Ch != 0 {
					editedMsg = append(editedMsg[:cursorPos], append([]rune{ev.Ch}, editedMsg[cursorPos:]...)...)
					cursorPos++
					if cursorPos >= scrollX+width-1 {
						scrollX = cursorPos - width + 2
					}
				}
			}

			// Update max scroll if needed
			if len(editedMsg) > width-1 {
				maxScroll = len(editedMsg) - width + 1
			} else {
				maxScroll = 0
			}
			if scrollX > maxScroll {
				scrollX = maxScroll
			}

			redraw()

		case termbox.EventError:
			return "", ev.Err
		}
	}
}

func getUserChoice(messages []string) (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	selected := 0
	drawMessages(messages, selected, false)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if selected > 0 {
					selected--
					drawMessages(messages, selected, false)
				}
			case termbox.KeyArrowDown:
				if selected < len(messages) {
					selected++
					drawMessages(messages, selected, false)
				}
			case termbox.KeyEnter:
				termbox.Close()
				if selected == len(messages) {
					// Custom message input
					msg, err := editMessage("")
					if err != nil {
						if err.Error() == "edit cancelled" {
							// Reopen the selection screen if edit was cancelled
							err = termbox.Init()
							if err != nil {
								return "", err
							}
							drawMessages(messages, selected, false)
							continue
						}
						return "", err
					}
					return msg, nil
				}

				// Edit selected message
				editedMsg, err := editMessage(messages[selected])
				if err != nil {
					if err.Error() == "edit cancelled" {
						// Reopen the selection screen if edit was cancelled
						err = termbox.Init()
						if err != nil {
							return "", err
						}
						drawMessages(messages, selected, false)
						continue
					}
					return "", err
				}
				return editedMsg, nil
			case termbox.KeyEsc:
				termbox.Close()
				return "", fmt.Errorf("selection cancelled")
			}
		case termbox.EventError:
			return "", ev.Err
		}
	}
}

// ============================================================================
// Delayed Commit Feature - Time/Date Selection Logic and Utilities
// ============================================================================

// isTimeInRestrictedRange checks if the given time falls within the restricted hour range
// Returns true if time is restricted, false otherwise
// Handles edge case where end < start (overnight range)
func isTimeInRestrictedRange(now time.Time, startHour, endHour int) bool {
	currentHour := now.Hour()

	// Handle normal range (e.g., 9-17)
	if startHour < endHour {
		return currentHour >= startHour && currentHour < endHour
	}

	// Handle overnight range (e.g., 22-6)
	if startHour > endHour {
		return currentHour >= startHour || currentHour < endHour
	}

	// If startHour == endHour, no restriction
	return false
}

// parseTimeString parses time strings in multiple formats
// Supports: "HH:MM", "HH:MM AM/PM", "HHhMM"
// Returns hour (0-23) and minute (0-59)
func parseTimeString(timeStr string) (hour, minute int, err error) {
	timeStr = strings.TrimSpace(timeStr)

	// Try parsing "HH:MM AM/PM" format
	if strings.Contains(strings.ToUpper(timeStr), "AM") || strings.Contains(strings.ToUpper(timeStr), "PM") {
		isPM := strings.Contains(strings.ToUpper(timeStr), "PM")
		// Remove AM/PM
		timeStr = strings.TrimSuffix(strings.TrimSuffix(strings.ToUpper(timeStr), "PM"), "AM")
		timeStr = strings.TrimSpace(timeStr)

		// Parse HH:MM
		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid time format. Use HH:MM AM/PM (e.g., 6:30 PM)")
		}

		h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour: %v", err)
		}

		m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute: %v", err)
		}

		// Convert 12-hour to 24-hour format
		if isPM && h != 12 {
			h += 12
		} else if !isPM && h == 12 {
			h = 0
		}

		if h < 0 || h > 23 {
			return 0, 0, fmt.Errorf("hour must be between 0 and 23")
		}
		if m < 0 || m > 59 {
			return 0, 0, fmt.Errorf("minute must be between 0 and 59")
		}

		return h, m, nil
	}

	// Try parsing "HHhMM" format
	if strings.Contains(timeStr, "h") || strings.Contains(timeStr, "H") {
		parts := strings.FieldsFunc(timeStr, func(r rune) bool {
			return r == 'h' || r == 'H'
		})

		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid time format. Use HHhMM (e.g., 18h30)")
		}

		h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour: %v", err)
		}

		m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute: %v", err)
		}

		if h < 0 || h > 23 {
			return 0, 0, fmt.Errorf("hour must be between 0 and 23")
		}
		if m < 0 || m > 59 {
			return 0, 0, fmt.Errorf("minute must be between 0 and 59")
		}

		return h, m, nil
	}

	// Try parsing "HH:MM" format (24-hour)
	if strings.Contains(timeStr, ":") {
		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid time format. Use HH:MM (e.g., 18:30)")
		}

		h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour: %v", err)
		}

		m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute: %v", err)
		}

		if h < 0 || h > 23 {
			return 0, 0, fmt.Errorf("hour must be between 0 and 23")
		}
		if m < 0 || m > 59 {
			return 0, 0, fmt.Errorf("minute must be between 0 and 59")
		}

		return h, m, nil
	}

	return 0, 0, fmt.Errorf("invalid time format. Use HH:MM, HH:MM AM/PM, or HHhMM")
}

// validateHour validates that hour is in range 0-23
func validateHour(hour int) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23, got %d", hour)
	}
	return nil
}

// generateSuggestedTimes generates suggested times starting from the end of restricted hours
// Uses the configured suggestion interval (e.g., 20 minutes)
// Generates 6-8 suggestions, only for today (doesn't go past 23:59)
func generateSuggestedTimes(restrictedEndHour, interval int) []string {
	suggestions := []string{}

	// Start from the end of restricted hours
	currentHour := restrictedEndHour
	currentMinute := 0

	// If interval doesn't evenly divide 60, start with first interval after restriction end
	if interval > 0 && interval < 60 {
		// Calculate first suggestion time
		currentMinute = interval
		if currentMinute >= 60 {
			currentHour++
			currentMinute = 0
		}
	}

	// Generate 6-8 suggestions
	maxSuggestions := 8
	for i := 0; i < maxSuggestions; i++ {
		// Stop if we've gone past 23:59
		if currentHour >= 24 {
			break
		}

		// Format time as HH:MM
		timeStr := fmt.Sprintf("%02d:%02d", currentHour, currentMinute)
		suggestions = append(suggestions, timeStr)

		// Increment by interval
		currentMinute += interval
		if currentMinute >= 60 {
			currentHour += currentMinute / 60
			currentMinute = currentMinute % 60
		}
	}

	return suggestions
}

// getNextAvailableTime calculates the next available time after restrictions end
// Returns the time for today, or current time if no time available today
func getNextAvailableTime(now time.Time, restrictedEndHour int) time.Time {
	// Create time for end of restriction today
	endTime := time.Date(now.Year(), now.Month(), now.Day(),
		restrictedEndHour, 0, 0, 0, now.Location())

	// If end time is still in the future today, return it
	if endTime.After(now) {
		return endTime
	}

	// Otherwise, no available time today, return current time as fallback
	return now
}

// showTimeSelectionUI displays interactive UI for time selection
// Similar to getUserChoice() using termbox
// Returns selected time string, custom flag (bool), and error
func showTimeSelectionUI(suggestedTimes []string) (string, bool, error) {
	err := termbox.Init()
	if err != nil {
		return "", false, err
	}
	defer termbox.Close()

	selected := 0
	customOption := len(suggestedTimes) // Index for "Enter custom time" option

	// Draw the UI
	drawTimeSelection := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Get terminal width
		width, _ := termbox.Size()
		if width < 10 {
			width = 80
		}

		// Draw title
		title := "Commit during restricted hours detected"
		for i, ch := range title {
			termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
		}

		// Draw subtitle
		subtitle := "Select commit time:"
		for i, ch := range subtitle {
			termbox.SetCell(i, 2, ch, termbox.ColorCyan, termbox.ColorDefault)
		}

		// Draw suggested times
		for i, timeStr := range suggestedTimes {
			fg := termbox.ColorDefault
			bg := termbox.ColorDefault
			if i == selected {
				fg = termbox.ColorBlack
				bg = termbox.ColorGreen
			}

			// Add bullet point
			bullet := "•"
			if i == selected {
				bullet = "→"
			}
			termbox.SetCell(2, i+4, []rune(bullet)[0], fg, bg)

			// Draw time
			displayText := fmt.Sprintf("%s", timeStr)
			for j, ch := range displayText {
				termbox.SetCell(j+4, i+4, ch, fg, bg)
			}
		}

		// Draw custom time option
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault
		if selected == customOption {
			fg = termbox.ColorBlack
			bg = termbox.ColorGreen
		}

		bullet := "•"
		if selected == customOption {
			bullet = "→"
		}
		termbox.SetCell(2, customOption+4, []rune(bullet)[0], fg, bg)

		customMsg := "Enter custom time (c)"
		for j, ch := range customMsg {
			termbox.SetCell(j+4, customOption+4, ch, fg, bg)
		}

		// Draw instructions
		instructions := "↑↓: Move  Enter: Select  c: Custom  Esc: Cancel"
		for i, ch := range instructions {
			termbox.SetCell(i, customOption+6, ch, termbox.ColorCyan, termbox.ColorDefault)
		}

		termbox.Flush()
	}

	drawTimeSelection()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if selected > 0 {
					selected--
					drawTimeSelection()
				}
			case termbox.KeyArrowDown:
				if selected < customOption {
					selected++
					drawTimeSelection()
				}
			case termbox.KeyEnter:
				termbox.Close()
				if selected == customOption {
					// User wants to enter custom time
					return "", true, nil
				}
				// Return selected time
				return suggestedTimes[selected], false, nil
			case termbox.KeyEsc:
				return "", false, fmt.Errorf("time selection cancelled")
			default:
				// Handle 'c' or 'C' key for custom time
				if ev.Ch == 'c' || ev.Ch == 'C' {
					termbox.Close()
					return "", true, nil
				}
			}
		case termbox.EventError:
			return "", false, ev.Err
		}
	}
}

// handleCustomTimeInput creates interactive prompt for custom time entry
// Accepts formats: "HH:MM", "HH:MM AM/PM", "HHhMM"
// Returns formatted time string or error
func handleCustomTimeInput() (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw title
	title := "Enter custom commit time"
	for i, ch := range title {
		termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
	}

	// Draw format instructions
	instructions := []string{
		"",
		"Accepted formats:",
		"  - HH:MM (24-hour): 18:45, 09:30",
		"  - HHhMM: 18h45, 9h30",
		"  - HH:MM AM/PM: 6:45 PM, 9:30 AM",
		"",
		"Time: ",
	}

	for i, line := range instructions {
		for j, ch := range line {
			termbox.SetCell(j, i+2, ch, termbox.ColorCyan, termbox.ColorDefault)
		}
	}

	// Draw bottom instructions
	bottomInst := "Enter: Confirm  Esc: Cancel"
	for i, ch := range bottomInst {
		termbox.SetCell(i, len(instructions)+4, ch, termbox.ColorCyan, termbox.ColorDefault)
	}

	// Input buffer
	input := []rune{}
	cursorPos := 0

	redraw := func(errorMsg string) {
		// Redraw prompt area
		for i := 0; i < 80; i++ {
			termbox.SetCell(i, len(instructions)+1, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}

		prompt := "Time: "
		for j, ch := range prompt {
			termbox.SetCell(j, len(instructions)+1, ch, termbox.ColorCyan, termbox.ColorDefault)
		}

		// Draw input
		for j, ch := range input {
			termbox.SetCell(j+len(prompt), len(instructions)+1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}

		// Draw error if any
		if errorMsg != "" {
			for i := 0; i < 80; i++ {
				termbox.SetCell(i, len(instructions)+3, ' ', termbox.ColorDefault, termbox.ColorDefault)
			}
			for j, ch := range errorMsg {
				termbox.SetCell(j, len(instructions)+3, ch, termbox.ColorRed, termbox.ColorDefault)
			}
		}

		termbox.SetCursor(len(prompt)+cursorPos, len(instructions)+1)
		termbox.Flush()
	}

	redraw("")

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				// Validate and return input
				timeStr := string(input)
				if timeStr == "" {
					redraw("Error: Time cannot be empty")
					continue
				}

				hour, minute, err := parseTimeString(timeStr)
				if err != nil {
					redraw(fmt.Sprintf("Error: %v", err))
					continue
				}

				// Return formatted time
				return fmt.Sprintf("%02d:%02d", hour, minute), nil

			case termbox.KeyEsc:
				return "", fmt.Errorf("custom time input cancelled")

			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorPos > 0 {
					input = append(input[:cursorPos-1], input[cursorPos:]...)
					cursorPos--
					redraw("")
				}

			case termbox.KeyDelete:
				if cursorPos < len(input) {
					input = append(input[:cursorPos], input[cursorPos+1:]...)
					redraw("")
				}

			case termbox.KeyArrowLeft:
				if cursorPos > 0 {
					cursorPos--
					redraw("")
				}

			case termbox.KeyArrowRight:
				if cursorPos < len(input) {
					cursorPos++
					redraw("")
				}

			case termbox.KeySpace:
				input = append(input[:cursorPos], append([]rune{' '}, input[cursorPos:]...)...)
				cursorPos++
				redraw("")

			default:
				if ev.Ch != 0 {
					input = append(input[:cursorPos], append([]rune{ev.Ch}, input[cursorPos:]...)...)
					cursorPos++
					redraw("")
				}
			}

		case termbox.EventError:
			return "", ev.Err
		}
	}
}

// executeDelayedCommit executes git commit with custom timestamp
// Uses git commit --date flag for both author and committer date
func executeDelayedCommit(commitMsg, timeStr string) error {
	// Parse the time string
	hour, minute, err := parseTimeString(timeStr)
	if err != nil {
		return fmt.Errorf("invalid time string: %v", err)
	}

	// Construct timestamp for today
	now := time.Now()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(),
		hour, minute, 0, 0, now.Location())

	// Format timestamp for git (ISO 8601)
	dateStr := timestamp.Format(time.RFC3339)

	// Create command with --date flag
	cmd := exec.Command("git", "commit", "-m", commitMsg, "--date", dateStr)

	// Also set committer date via environment variable
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute commit
	if err := cmd.Run(); err != nil {
		// Fallback: try without custom date
		log.Printf("Warning: Failed to commit with custom date, using current time: %v", err)
		fallbackCmd := exec.Command("git", "commit", "-m", commitMsg)
		fallbackCmd.Stdout = os.Stdout
		fallbackCmd.Stderr = os.Stderr
		return fallbackCmd.Run()
	}

	fmt.Printf("Commit timestamp: %s\n", timestamp.Format("2006-01-02 15:04:05 -0700"))
	return nil
}

func isValidAPIKey(apiKey string) bool {
	// Google AI API keys typically start with "AIza" and are 39 characters long
	return len(apiKey) == 39 && strings.HasPrefix(apiKey, "AIza")
}

func main() {
	configFlag := flag.Bool("config", false, "Configure API key")
	detailedFlag := flag.Bool("d", false, "Generate detailed commit messages with body text")
	enableLoggingFlag := flag.Bool("enable-logging", false, "Enable request logging to file")
	disableLoggingFlag := flag.Bool("disable-logging", false, "Disable request logging to file")
	iconFlag := flag.Bool("icon", false, "Use emoji icons in commit messages")
	configDelayedFlag := flag.Bool("config-delayed", false, "Configure delayed commit settings")
	enableDelayedFlag := flag.Bool("enable-delayed", false, "Enable delayed commit feature")
	disableDelayedFlag := flag.Bool("disable-delayed", false, "Disable delayed commit feature")
	flag.Parse()

	// Handle logging configuration flags
	if *enableLoggingFlag {
		if err := config.SetLoggingEnabled(true); err != nil {
			log.Fatalf("Failed to enable logging: %v", err)
		}
		fmt.Println("Logging enabled successfully!")
		return
	}

	if *disableLoggingFlag {
		if err := config.SetLoggingEnabled(false); err != nil {
			log.Fatalf("Failed to disable logging: %v", err)
		}
		fmt.Println("Logging disabled successfully!")
		return
	}

	// Handle delayed commit configuration flags
	if *configDelayedFlag {
		fmt.Println("Configure Delayed Commit Settings")
		fmt.Println("===================================")
		fmt.Println()

		// Get restricted start hour
		fmt.Print("Enter restricted start hour (0-23, e.g., 9 for 9 AM): ")
		var startHourStr string
		fmt.Scanln(&startHourStr)
		startHour, err := strconv.Atoi(strings.TrimSpace(startHourStr))
		if err != nil || startHour < 0 || startHour > 23 {
			log.Fatal("Error: Invalid start hour. Must be between 0 and 23.")
		}

		// Get restricted end hour
		fmt.Print("Enter restricted end hour (0-23, e.g., 17 for 5 PM): ")
		var endHourStr string
		fmt.Scanln(&endHourStr)
		endHour, err := strconv.Atoi(strings.TrimSpace(endHourStr))
		if err != nil || endHour < 0 || endHour > 23 {
			log.Fatal("Error: Invalid end hour. Must be between 0 and 23.")
		}

		// Get suggestion interval
		fmt.Print("Enter suggestion interval in minutes (e.g., 20, 30, 60): ")
		var intervalStr string
		fmt.Scanln(&intervalStr)
		interval, err := strconv.Atoi(strings.TrimSpace(intervalStr))
		if err != nil || interval <= 0 || interval > 1440 {
			log.Fatal("Error: Invalid interval. Must be between 1 and 1440 minutes.")
		}

		// Save configuration
		if err := config.SetDelayedCommitConfig(startHour, endHour, interval); err != nil {
			log.Fatalf("Failed to save delayed commit configuration: %v", err)
		}

		fmt.Println()
		fmt.Println("Delayed commit settings configured successfully!")
		fmt.Printf("Restricted hours: %02d:00 - %02d:00\n", startHour, endHour)
		fmt.Printf("Suggestion interval: %d minutes\n", interval)
		return
	}

	if *enableDelayedFlag {
		if err := config.SetDelayedCommitEnabled(true); err != nil {
			log.Fatalf("Failed to enable delayed commit: %v", err)
		}
		fmt.Println("Delayed commit feature enabled successfully!")
		return
	}

	if *disableDelayedFlag {
		if err := config.SetDelayedCommitEnabled(false); err != nil {
			log.Fatalf("Failed to disable delayed commit: %v", err)
		}
		fmt.Println("Delayed commit feature disabled successfully!")
		return
	}

	if *configFlag {
		fmt.Print("Enter your Google AI API key: ")
		var apiKey string
		fmt.Scanln(&apiKey)

		// Trim whitespace and check if empty
		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			log.Fatal("Error: API key cannot be empty")
		}

		// Check if API key format is correct
		if !isValidAPIKey(apiKey) {
			log.Fatal("Error: Invalid API key format. Google AI API keys should start with 'AIza' and be 39 characters long.")
		}

		if err := config.SetAPIKey(apiKey); err != nil {
			log.Fatalf("Failed to save API key: %v", err)
		}
		fmt.Println("API key configured successfully!")
		return
	}

	apiKey, err := config.GetAPIKey()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	diff, err := getGitDiff()
	if err != nil {
		log.Fatal(err)
	}

	if diff == "" {
		log.Fatal("No staged changes found. Please stage your changes using 'git add' first.")
	}

	messages, gitDiff, lastCommitMsg, prompt, err := generateCommitMessages(diff, apiKey, *detailedFlag, *iconFlag)
	if err != nil {
		// Log error
		logger.LogError(gitDiff, lastCommitMsg, prompt, err.Error())
		log.Fatal(err)
	}

	commitMsg, err := getUserChoice(messages)
	if err != nil {
		log.Fatal(err)
	}

	// Get the full AI response for logging
	aiResponse := strings.Join(messages, "\n")

	// Log successful request
	logger.LogSuccess(gitDiff, lastCommitMsg, prompt, aiResponse, messages, commitMsg)

	// Check if delayed commit feature is enabled
	delayedCommitUsed := false
	enabled, err := config.IsDelayedCommitEnabled()
	if err == nil && enabled {
		now := time.Now()

		// Get delayed commit configuration
		delayedConfig, err := config.GetDelayedCommitConfig()
		if err == nil {
			// Check if current time is in restricted range
			if isTimeInRestrictedRange(now, delayedConfig.RestrictedStartHour, delayedConfig.RestrictedEndHour) {
				fmt.Println()
				fmt.Println("⏰ Commit during restricted hours detected!")
				fmt.Printf("Restricted hours: %02d:00 - %02d:00\n", delayedConfig.RestrictedStartHour, delayedConfig.RestrictedEndHour)
				fmt.Println()

				// Generate suggested times
				suggestedTimes := generateSuggestedTimes(delayedConfig.RestrictedEndHour, delayedConfig.SuggestionInterval)

				// Show time selection UI
				selectedTime, isCustom, err := showTimeSelectionUI(suggestedTimes)
				if err != nil {
					if err.Error() == "time selection cancelled" {
						fmt.Println("Commit cancelled.")
						return
					}
					log.Fatalf("Error in time selection: %v", err)
				}

				// Handle custom time input if requested
				if isCustom {
					selectedTime, err = handleCustomTimeInput()
					if err != nil {
						if err.Error() == "custom time input cancelled" {
							fmt.Println("Commit cancelled.")
							return
						}
						log.Fatalf("Error in custom time input: %v", err)
					}
				}

				// Execute delayed commit
				fmt.Printf("\nScheduling commit for: %s\n", selectedTime)
				if err := executeDelayedCommit(commitMsg, selectedTime); err != nil {
					log.Fatalf("Failed to create delayed commit: %v", err)
				}

				delayedCommitUsed = true
				fmt.Printf("✓ Successfully created commit with message: %s\n", commitMsg)
			}
		}
	}

	// Execute normal commit if delayed commit was not used
	if !delayedCommitUsed {
		// Create git commit with the chosen message
		cmd := exec.Command("git", "commit", "-m", commitMsg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to create commit:", err)
		}

		fmt.Printf("Successfully created commit with message: %s\n", commitMsg)
	}

	// Show log file location
	// logPath, _ := logger.GetLogPath()
	// fmt.Printf("Request logged to: %s\n", logPath)
}
