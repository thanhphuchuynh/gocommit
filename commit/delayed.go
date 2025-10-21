package commit

import (
	"fmt"
	"time"

	"github.com/thanhphuchuynh/config"
	"github.com/thanhphuchuynh/git"
	"github.com/thanhphuchuynh/internal/timeutil"
	"github.com/thanhphuchuynh/ui"
)

// HandleDelayedCommit checks if delayed commit should be used and handles it
// Returns true if delayed commit was used, false otherwise
func HandleDelayedCommit(message string) (bool, error) {
	// Check if delayed commit feature is enabled
	enabled, err := config.IsDelayedCommitEnabled()
	if err != nil || !enabled {
		return false, nil
	}

	now := time.Now()

	// Get delayed commit configuration
	delayedConfig, err := config.GetDelayedCommitConfig()
	if err != nil {
		return false, fmt.Errorf("failed to get delayed commit config: %w", err)
	}

	// Check if current time is in restricted range
	if !timeutil.IsInRestrictedRange(now, delayedConfig.RestrictedStartHour, delayedConfig.RestrictedEndHour) {
		return false, nil
	}

	// We're in restricted hours - show delayed commit UI
	fmt.Println()
	fmt.Println("⏰ Commit during restricted hours detected!")
	fmt.Printf("Restricted hours: %02d:00 - %02d:00\n", delayedConfig.RestrictedStartHour, delayedConfig.RestrictedEndHour)
	fmt.Println()

	// Generate suggested times
	suggestedTimes := timeutil.GenerateSuggestedTimes(delayedConfig.RestrictedEndHour, delayedConfig.SuggestionInterval)

	// Show time selection UI
	selectedTime, isCustom, err := ui.ShowTimeSelector(suggestedTimes)
	if err != nil {
		if err.Error() == "time selection cancelled" {
			return false, fmt.Errorf("commit cancelled by user")
		}
		return false, fmt.Errorf("error in time selection: %w", err)
	}

	// Handle custom time input if requested
	if isCustom {
		selectedTime, err = ui.ShowCustomTimeInput()
		if err != nil {
			if err.Error() == "custom time input cancelled" {
				return false, fmt.Errorf("commit cancelled by user")
			}
			return false, fmt.Errorf("error in custom time input: %w", err)
		}
	}

	// Execute delayed commit
	fmt.Printf("\nScheduling commit for: %s\n", selectedTime)
	if err := git.CommitWithDate(message, selectedTime); err != nil {
		return false, fmt.Errorf("failed to create delayed commit: %w", err)
	}

	fmt.Printf("✓ Successfully created commit with message: %s\n", message)
	return true, nil
}
