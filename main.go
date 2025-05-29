package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/nsf/termbox-go"
	"github.com/tphuc/gocommit/config"
	"github.com/tphuc/gocommit/logger"
	"google.golang.org/api/option"
)

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

Generate exactly 3 different commit messages, numbered 1-3:`

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

func generateCommitMessages(diff string, apiKey string) ([]string, string, string, string, error) {
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
	prompt := fmt.Sprintf(promptTemplate, diff, lastCommitMsg)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, diff, lastCommitMsg, prompt, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, diff, lastCommitMsg, prompt, fmt.Errorf("no content generated")
	}

	// Get the text content from the response
	text := ""
	fmt.Printf("Generated response: %+v \n", resp.Candidates[0].Content.Parts)

	for _, part := range resp.Candidates[0].Content.Parts {
		if str, ok := part.(genai.Text); ok {
			text += string(str)
		}
	}

	// Split the response into individual messages
	messages := strings.Split(strings.TrimSpace(text), "\n")
	if len(messages) < 3 {
		return nil, diff, lastCommitMsg, prompt, fmt.Errorf("expected 3 messages, got %d", len(messages))
	}

	// Clean up the messages (remove numbers and extra spaces)
	cleanMessages := make([]string, 3)
	for i, msg := range messages[:3] {
		// Remove the number prefix if present
		parts := strings.SplitN(msg, " ", 2)
		if len(parts) > 1 {
			cleanMessages[i] = strings.TrimSpace(parts[1])
		} else {
			cleanMessages[i] = strings.TrimSpace(msg)
		}
	}

	return cleanMessages, diff, lastCommitMsg, prompt, nil
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
		// Draw message
		for j, ch := range msg {
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

	// Function to redraw the message with wrapping
	redraw := func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Draw title and instructions
		for i, ch := range title {
			termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
		}
		for i, ch := range instructions {
			termbox.SetCell(i, 1, ch, termbox.ColorCyan, termbox.ColorDefault)
		}

		// Split message into lines
		lines := strings.Split(string(editedMsg), "\n")
		if currentLine >= len(lines) {
			currentLine = len(lines) - 1
		}

		// Draw visible portion of the message
		line := 2
		for _, msgLine := range lines {
			// Calculate visible portion for this line
			visibleStart := scrollX
			visibleEnd := scrollX + width - 1
			if visibleEnd > len(msgLine) {
				visibleEnd = len(msgLine)
			}

			// Draw the line
			col := 0
			for j := visibleStart; j < visibleEnd; j++ {
				if col >= width-1 {
					line++
					col = 0
				}
				termbox.SetCell(col, line, rune(msgLine[j]), termbox.ColorDefault, termbox.ColorDefault)
				col++
			}
			line++
		}

		// Calculate and set cursor position
		linesBeforeCursor := strings.Split(string(editedMsg[:cursorPos]), "\n")
		cursorLine := 2 + len(linesBeforeCursor) - 1
		lastLine := linesBeforeCursor[len(linesBeforeCursor)-1]
		cursorCol := len(lastLine) % (width - 1)
		termbox.SetCursor(cursorCol, cursorLine)
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

func isValidAPIKey(apiKey string) bool {
	// Google AI API keys typically start with "AIza" and are 39 characters long
	return len(apiKey) == 39 && strings.HasPrefix(apiKey, "AIza")
}

func main() {
	configFlag := flag.Bool("config", false, "Configure API key")
	flag.Parse()

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

	messages, gitDiff, lastCommitMsg, prompt, err := generateCommitMessages(diff, apiKey)
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

	// Create git commit with the chosen message
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to create commit:", err)
	}

	fmt.Printf("Successfully created commit with message: %s\n", commitMsg)

	// Show log file location
	logPath, _ := logger.GetLogPath()
	fmt.Printf("Request logged to: %s\n", logPath)
}
