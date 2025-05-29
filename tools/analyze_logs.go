package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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

func getLogPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(homeDir, ".gocommit", "gocommit_requests.log"), nil
}

func loadLogs() ([]LogEntry, error) {
	logPath, err := getLogPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []LogEntry{}, nil
		}
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			log.Printf("Warning: Failed to parse line %d: %v", lineNum, err)
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %v", err)
	}

	return entries, nil
}

func analyzeCommitTypes(entries []LogEntry) {
	typeCount := make(map[string]int)

	for _, entry := range entries {
		if !entry.Success {
			continue
		}

		// Extract commit type from selected message
		msg := entry.SelectedMsg
		if idx := strings.Index(msg, ":"); idx > 0 {
			commitType := strings.TrimSpace(msg[:idx])
			if strings.Contains(commitType, "(") {
				// Remove scope: feat(scope) -> feat
				if parenIdx := strings.Index(commitType, "("); parenIdx > 0 {
					commitType = commitType[:parenIdx]
				}
			}
			typeCount[commitType]++
		}
	}

	fmt.Println("## Commit Type Usage")

	// Sort by frequency
	type kv struct {
		Key   string
		Value int
	}

	var sorted []kv
	for k, v := range typeCount {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	for _, item := range sorted {
		fmt.Printf("- %s: %d times\n", item.Key, item.Value)
	}
	fmt.Println()
}

func analyzeSuccess(entries []LogEntry) {
	successful := 0
	failed := 0

	for _, entry := range entries {
		if entry.Success {
			successful++
		} else {
			failed++
		}
	}

	total := successful + failed
	if total == 0 {
		fmt.Println("## Success Rate: No entries found")
		return
	}

	successRate := float64(successful) / float64(total) * 100
	fmt.Printf("## Success Rate\n")
	fmt.Printf("- Total requests: %d\n", total)
	fmt.Printf("- Successful: %d (%.1f%%)\n", successful, successRate)
	fmt.Printf("- Failed: %d (%.1f%%)\n", failed, 100-successRate)
	fmt.Println()
}

func analyzeUserChoices(entries []LogEntry) {
	choiceIndex := make(map[int]int) // 0=first option, 1=second, 2=third, 3=custom

	for _, entry := range entries {
		if !entry.Success || len(entry.GeneratedMsgs) < 3 {
			continue
		}

		selected := entry.SelectedMsg
		found := false

		// Check which generated message was selected
		for i, msg := range entry.GeneratedMsgs {
			if msg == selected {
				choiceIndex[i]++
				found = true
				break
			}
		}

		// If not found in generated messages, it was custom/edited
		if !found {
			choiceIndex[3]++
		}
	}

	fmt.Println("## User Choice Patterns")
	total := 0
	for _, count := range choiceIndex {
		total += count
	}

	if total == 0 {
		fmt.Println("No choice data available")
		return
	}

	labels := []string{"First option", "Second option", "Third option", "Custom/Edited"}
	for i := 0; i < 4; i++ {
		count := choiceIndex[i]
		percentage := float64(count) / float64(total) * 100
		fmt.Printf("- %s: %d times (%.1f%%)\n", labels[i], count, percentage)
	}
	fmt.Println()
}

func showRecentActivity(entries []LogEntry, days int) {
	cutoff := time.Now().AddDate(0, 0, -days)
	recent := make([]LogEntry, 0)

	for _, entry := range entries {
		if entry.Timestamp.After(cutoff) {
			recent = append(recent, entry)
		}
	}

	fmt.Printf("## Recent Activity (Last %d days)\n", days)
	fmt.Printf("- Total requests: %d\n", len(recent))

	if len(recent) == 0 {
		fmt.Println("No recent activity")
		return
	}

	// Group by day
	dayCount := make(map[string]int)
	for _, entry := range recent {
		day := entry.Timestamp.Format("2006-01-02")
		dayCount[day]++
	}

	// Sort days
	var days_list []string
	for day := range dayCount {
		days_list = append(days_list, day)
	}
	sort.Strings(days_list)

	fmt.Println("- Daily breakdown:")
	for _, day := range days_list {
		fmt.Printf("  - %s: %d requests\n", day, dayCount[day])
	}
	fmt.Println()
}

func main() {
	entries, err := loadLogs()
	if err != nil {
		log.Fatalf("Failed to load logs: %v", err)
	}

	if len(entries) == 0 {
		fmt.Println("No log entries found. Use gocommit to generate some data first!")
		return
	}

	fmt.Printf("# gocommit Log Analysis\n\n")
	fmt.Printf("Analyzed %d log entries\n\n", len(entries))

	analyzeSuccess(entries)
	analyzeCommitTypes(entries)
	analyzeUserChoices(entries)
	showRecentActivity(entries, 7)

	logPath, _ := getLogPath()
	fmt.Printf("Log file: %s\n", logPath)
}
