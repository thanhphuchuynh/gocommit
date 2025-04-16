package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const promptTemplate = `As an AI commit message generator, analyze the following git diff and generate a concise, meaningful commit message following conventional commits format (type(scope): description). The commit types should be one of: feat, fix, docs, style, refactor, perf, test, or chore.

Git diff:
%s

Generate a commit message that is:
1. Concise (max 100 characters)
2. Descriptive
3. Following the format: type(scope): description
4. Based on the actual code changes shown in the diff

Commit message only, no explanation needed.`

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

func generateCommitMessage(diff string, apiKey string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	prompt := fmt.Sprintf(promptTemplate, diff)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	// Get the text content from the response
	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if str, ok := part.(genai.Text); ok {
			text += string(str)
		}
	}

	return strings.TrimSpace(text), nil
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

	commitMsg, err := generateCommitMessage(diff, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create git commit with the generated message
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to create commit:", err)
	}

	fmt.Printf("Successfully created commit with message: %s\n", commitMsg)
}
