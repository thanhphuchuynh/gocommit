// Package timeutil provides time-related utility functions for handling
// delayed commits and time-based restrictions.
package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// IsInRestrictedRange checks if the given time falls within the restricted hour range.
// Returns true if time is restricted, false otherwise.
// Handles edge case where end < start (overnight range).
//
// Parameters:
//   - now: The time to check
//   - startHour: The start of restricted hours (0-23)
//   - endHour: The end of restricted hours (0-23)
//
// Examples:
//   - Normal range (9-17): Returns true if hour is between 9 and 17
//   - Overnight range (22-6): Returns true if hour is >= 22 or < 6
//   - Same start/end: No restriction, returns false
func IsInRestrictedRange(now time.Time, startHour, endHour int) bool {
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

// ParseTimeString parses time strings in multiple formats.
// Supports: "HH:MM", "HH:MM AM/PM", "HHhMM"
// Returns hour (0-23) and minute (0-59), or an error if parsing fails.
//
// Supported formats:
//   - "HH:MM" (24-hour): "18:30", "09:15"
//   - "HH:MM AM/PM" (12-hour): "6:30 PM", "9:15 AM"
//   - "HHhMM": "18h30", "9h15"
//
// Parameters:
//   - timeStr: The time string to parse
//
// Returns:
//   - hour: The parsed hour (0-23)
//   - minute: The parsed minute (0-59)
//   - err: Error if parsing fails or values are invalid
func ParseTimeString(timeStr string) (hour, minute int, err error) {
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

// ValidateHour validates that hour is in valid range 0-23.
//
// Parameters:
//   - hour: The hour value to validate
//
// Returns:
//   - error: Non-nil if hour is invalid
func ValidateHour(hour int) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 0 and 23, got %d", hour)
	}
	return nil
}

// GenerateSuggestedTimes generates time suggestions starting from the end of restricted hours.
// Uses the configured suggestion interval (e.g., 20 minutes).
// Generates 6-8 suggestions, only for today (doesn't go past 23:59).
//
// Parameters:
//   - restrictedEndHour: The hour when restrictions end (0-23)
//   - intervalMinutes: The interval between suggestions in minutes
//
// Returns:
//   - []string: A slice of time suggestions in "HH:MM" format
//
// Example:
//   - restrictedEndHour=17, intervalMinutes=20 might return:
//     ["17:20", "17:40", "18:00", "18:20", "18:40", "19:00", "19:20", "19:40"]
func GenerateSuggestedTimes(restrictedEndHour, intervalMinutes int) []string {
	suggestions := []string{}

	// Start from the end of restricted hours
	currentHour := restrictedEndHour
	currentMinute := 0

	// If interval doesn't evenly divide 60, start with first interval after restriction end
	if intervalMinutes > 0 && intervalMinutes < 60 {
		// Calculate first suggestion time
		currentMinute = intervalMinutes
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
		currentMinute += intervalMinutes
		if currentMinute >= 60 {
			currentHour += currentMinute / 60
			currentMinute = currentMinute % 60
		}
	}

	return suggestions
}

// GetNextAvailableTime calculates the next available time after restrictions end.
// Returns the time for today, or current time if no time available today.
//
// Parameters:
//   - now: The current time
//   - restrictedEndHour: The hour when restrictions end (0-23)
//
// Returns:
//   - time.Time: The next available time after restrictions
//
// Example:
//   - If now is 10:00 and restrictedEndHour is 17, returns 17:00 today
//   - If now is 18:00 and restrictedEndHour is 17, returns 18:00 (current time)
func GetNextAvailableTime(now time.Time, restrictedEndHour int) time.Time {
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
