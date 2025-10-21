// Package git provides Git operations for the gocommit tool.
// This package encapsulates all Git command execution and provides
// a clean interface for interacting with Git repositories.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GetStagedDiff returns the diff of staged changes in the Git repository.
// Returns an empty string if no changes are staged.
// Returns an error if the git diff command fails.
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting git diff: %v", err)
	}
	return string(output), nil
}

// GetLastCommitMessage returns the last commit message from the Git repository.
// Returns an empty string if there are no commits yet (new repository).
// This function is designed to never return an error for the "no commits" case,
// as it's a valid state for a new repository.
func GetLastCommitMessage() string {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.Output()
	if err != nil {
		// If there are no commits yet, git log will fail
		// This is not an error condition - just return empty string
		return ""
	}
	return strings.TrimSpace(string(output))
}

// Commit executes a standard git commit with the given message.
// The commit is created with the current timestamp.
// Stdout and stderr are connected to the parent process for user feedback.
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %v", err)
	}

	return nil
}

// CommitWithDate executes a git commit with a custom date/time.
// The timeStr parameter should be in "HH:MM" format (e.g., "14:30").
// The commit will be created with the specified time on the current date.
// Both author date and committer date are set to the specified time.
//
// If the custom date commit fails, it will automatically fall back to
// a standard commit with the current timestamp.
func CommitWithDate(message, timeStr string) error {
	// Parse the time string (format: "HH:MM")
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid time format: expected HH:MM, got %s", timeStr)
	}

	// Parse hour and minute
	var hour, minute int
	if _, err := fmt.Sscanf(parts[0], "%d", &hour); err != nil {
		return fmt.Errorf("invalid hour in time string: %v", err)
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &minute); err != nil {
		return fmt.Errorf("invalid minute in time string: %v", err)
	}

	// Validate hour and minute ranges
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23, got %d", hour)
	}
	if minute < 0 || minute > 59 {
		return fmt.Errorf("minute must be between 0 and 59, got %d", minute)
	}

	// Construct timestamp for today with specified time
	now := time.Now()
	timestamp := time.Date(now.Year(), now.Month(), now.Day(),
		hour, minute, 0, 0, now.Location())

	// Format timestamp for git (ISO 8601)
	dateStr := timestamp.Format(time.RFC3339)

	// Create command with --date flag for author date
	cmd := exec.Command("git", "commit", "-m", message, "--date", dateStr)

	// Set committer date via environment variable
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GIT_COMMITTER_DATE=%s", dateStr))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute commit with custom date
	if err := cmd.Run(); err != nil {
		// Fallback: try standard commit without custom date
		fmt.Printf("Warning: Failed to commit with custom date, using current time: %v\n", err)
		return Commit(message)
	}

	fmt.Printf("Commit timestamp: %s\n", timestamp.Format("2006-01-02 15:04:05 -0700"))
	return nil
}
