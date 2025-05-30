package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thanhphuchuynh/config"
)

type LogEntry struct {
	Timestamp     time.Time `json:"timestamp"`
	GitDiff       string    `json:"git_diff"`
	LastCommitMsg string    `json:"last_commit_msg"`
	PromptSent    string    `json:"prompt_sent"`
	AIResponse    string    `json:"ai_response"`
	GeneratedMsgs []string  `json:"generated_messages"`
	SelectedMsg   string    `json:"selected_message"`
	Success       bool      `json:"success"`
	ErrorMsg      string    `json:"error_msg,omitempty"`
}

const logFileName = "gocommit_requests.log"

func getLogPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(homeDir, ".gocommit", logFileName), nil
}

func ensureLogDir() error {
	logPath, err := getLogPath()
	if err != nil {
		return err
	}

	logDir := filepath.Dir(logPath)
	return os.MkdirAll(logDir, 0755)
}

func LogRequest(entry LogEntry) error {
	if err := ensureLogDir(); err != nil {
		return fmt.Errorf("failed to ensure log directory: %v", err)
	}

	logPath, err := getLogPath()
	if err != nil {
		return err
	}

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Set timestamp if not already set
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Marshal to JSON and write to file
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	// Write JSON entry followed by newline for easier parsing
	if _, err := file.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %v", err)
	}

	return nil
}

func LogSuccess(gitDiff, lastCommitMsg, prompt, aiResponse string, generatedMsgs []string, selectedMsg string) {
	// Check if logging is enabled
	loggingEnabled, err := config.IsLoggingEnabled()
	if err != nil {
		fmt.Printf("Warning: Failed to check logging config: %v\n", err)
		return
	}
	if !loggingEnabled {
		return
	}

	entry := LogEntry{
		GitDiff:       gitDiff,
		LastCommitMsg: lastCommitMsg,
		PromptSent:    prompt,
		AIResponse:    aiResponse,
		GeneratedMsgs: generatedMsgs,
		SelectedMsg:   selectedMsg,
		Success:       true,
	}

	if err := LogRequest(entry); err != nil {
		fmt.Printf("Warning: Failed to log request: %v\n", err)
	}
}

func LogError(gitDiff, lastCommitMsg, prompt string, errorMsg string) {
	// Check if logging is enabled
	loggingEnabled, err := config.IsLoggingEnabled()
	if err != nil {
		fmt.Printf("Warning: Failed to check logging config: %v\n", err)
		return
	}
	if !loggingEnabled {
		return
	}

	entry := LogEntry{
		GitDiff:       gitDiff,
		LastCommitMsg: lastCommitMsg,
		PromptSent:    prompt,
		Success:       false,
		ErrorMsg:      errorMsg,
	}

	if err := LogRequest(entry); err != nil {
		fmt.Printf("Warning: Failed to log error: %v\n", err)
	}
}

func GetLogPath() (string, error) {
	return getLogPath()
}
