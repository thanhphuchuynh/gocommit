package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/nsf/termbox-go"
	"google.golang.org/api/option"
)

const promptTemplate = `As an AI commit message generator, analyze the following git diff and generate 3 different concise, meaningful commit messages following conventional commits format (type(scope): description). The commit types should be one of: feat, fix, docs, style, refactor, perf, test, or chore.

Git diff:
%s

Last commit message format (for reference):
%s

Generate 3 different commit messages that are:
1. Concise (max 100 characters)
2. Descriptive
3. Following the format: type(scope): description
4. Based on the actual code changes shown in the diff
5. Following a similar style to the last commit message if applicable

Return each commit message on a new line, numbered 1-3. No explanation needed.`

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

func generateCommitMessages(diff string, apiKey string) ([]string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
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

	// Split the response into individual messages
	messages := strings.Split(strings.TrimSpace(text), "\n")
	if len(messages) < 3 {
		return nil, fmt.Errorf("expected 3 messages, got %d", len(messages))
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

	return cleanMessages, nil
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
					fmt.Print("\033[33mMessage: \033[0m")
					var customMsg string
					fmt.Scanln(&customMsg)
					return customMsg, nil
				}

				// Show edit prompt for selected message
				fmt.Printf("\033[33mEdit message (press Enter to keep as is): \033[0m%s", messages[selected])
				var editedMsg string
				fmt.Scanln(&editedMsg)
				if editedMsg == "" {
					return messages[selected], nil
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

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	diff, err := getGitDiff()
	if err != nil {
		log.Fatal(err)
	}

	if diff == "" {
		log.Fatal("No staged changes found. Please stage your changes using 'git add' first.")
	}

	messages, err := generateCommitMessages(diff, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	commitMsg, err := getUserChoice(messages)
	if err != nil {
		log.Fatal(err)
	}

	// Create git commit with the chosen message
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to create commit:", err)
	}

	fmt.Printf("Successfully created commit with message: %s\n", commitMsg)
}
