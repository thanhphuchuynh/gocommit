package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/rivo/tview"
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

func getUserChoice(messages []string) (string, error) {
	app := tview.NewApplication()

	// Create a simple list for the messages
	list := tview.NewList().
		ShowSecondaryText(false)

	// Add messages to the list
	for _, msg := range messages {
		list.AddItem(msg, "", 0, nil)
	}

	// Add custom edit option
	list.AddItem("Edit custom message", "", 0, nil)

	var selectedMessage string

	// Set up the list's selected function
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index == len(messages) { // Custom edit option
			app.Stop()
			// Create a simple input field for custom message
			inputField := tview.NewInputField().
				SetLabel("Message: ")

			form := tview.NewForm().
				AddFormItem(inputField).
				AddButton("OK", func() {
					selectedMessage = inputField.GetText()
					app.Stop()
				})

			if err := tview.NewApplication().SetRoot(form, true).Run(); err != nil {
				log.Fatal(err)
			}
		} else {
			selectedMessage = messages[index]
			app.Stop()
		}
	})

	// Run the application
	if err := app.SetRoot(list, true).Run(); err != nil {
		return "", err
	}

	if selectedMessage == "" {
		return "", fmt.Errorf("no message selected")
	}

	return selectedMessage, nil
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
