package commit

import (
	"context"
	"fmt"
	"strings"

	"github.com/thanhphuchuynh/ai"
	"github.com/thanhphuchuynh/git"
	"github.com/thanhphuchuynh/logger"
	"github.com/thanhphuchuynh/ui"
)

// WorkflowConfig contains configuration for the commit workflow
type WorkflowConfig struct {
	Detailed bool
	UseIcons bool
}

// Workflow orchestrates the commit process
type Workflow struct {
	config   WorkflowConfig
	provider ai.Provider
}

// NewWorkflow creates a new commit workflow
func NewWorkflow(cfg WorkflowConfig, provider ai.Provider) *Workflow {
	return &Workflow{
		config:   cfg,
		provider: provider,
	}
}

// Execute runs the complete commit workflow
func (w *Workflow) Execute(ctx context.Context) error {
	// 1. Get staged diff
	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged diff: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("no staged changes found. Please stage your changes using 'git add' first")
	}

	// 2. Get last commit message for reference
	lastCommitMsg := git.GetLastCommitMessage()

	// 3. Generate commit messages using AI provider
	req := &ai.GenerateRequest{
		Diff:       diff,
		LastCommit: lastCommitMsg,
		Detailed:   w.config.Detailed,
		UseIcons:   w.config.UseIcons,
	}

	resp, err := w.provider.GenerateMessages(ctx, req)
	if err != nil {
		// Log error (note: we don't have the prompt directly, so we'll pass empty string)
		logger.LogError(diff, lastCommitMsg, "", err.Error())
		return fmt.Errorf("failed to generate commit messages: %w", err)
	}

	messages := resp.Messages

	// 4. Show message selector UI
	commitMsg, err := ui.ShowMessageSelector(messages, "Select a commit message:")
	if err != nil {
		return fmt.Errorf("failed to select commit message: %w", err)
	}

	// Get the full AI response for logging
	aiResponse := strings.Join(messages, "\n")

	// Log successful request (note: we don't have the prompt directly, so we'll pass empty string)
	logger.LogSuccess(diff, lastCommitMsg, "", aiResponse, messages, commitMsg)

	// 5. Handle delayed commit if enabled
	delayedCommitUsed, err := HandleDelayedCommit(commitMsg)
	if err != nil {
		return fmt.Errorf("failed to handle delayed commit: %w", err)
	}

	// 6. Execute normal commit if delayed commit was not used
	if !delayedCommitUsed {
		if err := git.Commit(commitMsg); err != nil {
			return fmt.Errorf("failed to create commit: %w", err)
		}

		fmt.Printf("Successfully created commit with message: %s\n", commitMsg)
	}

	return nil
}
