package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/thanhphuchuynh/internal/timeutil"
)

// ShowTimeSelector displays interactive UI for time selection
// Returns selected time string, custom flag (bool), and error
func ShowTimeSelector(suggestedTimes []string) (string, bool, error) {
	err := termbox.Init()
	if err != nil {
		return "", false, err
	}
	defer termbox.Close()

	selected := 0
	customOption := len(suggestedTimes) // Index for "Enter custom time" option

	// Draw the UI
	drawTimeSelection := func() {
		ClearScreen()

		// Get terminal dimensions
		_, _ = GetTerminalSize()

		// Draw title
		title := "Commit during restricted hours detected"
		DrawTitle(title, 0, 0)

		// Draw subtitle
		subtitle := "Select commit time:"
		DrawInstructions(subtitle, 0, 2)

		// Draw suggested times
		for i, timeStr := range suggestedTimes {
			fg := ColorDefault
			bg := ColorDefault
			if i == selected {
				fg = ColorSelected
				bg = ColorSelectedBG
			}

			// Draw bullet point
			DrawBullet(2, i+4, i == selected, fg, bg)

			// Draw time
			displayText := fmt.Sprintf("%s", timeStr)
			DrawText(displayText, 4, i+4, fg, bg)
		}

		// Draw custom time option
		fg := ColorDefault
		bg := ColorDefault
		if selected == customOption {
			fg = ColorSelected
			bg = ColorSelectedBG
		}

		DrawBullet(2, customOption+4, selected == customOption, fg, bg)

		customMsg := "Enter custom time (c)"
		DrawText(customMsg, 4, customOption+4, fg, bg)

		// Draw instructions
		instructions := "↑↓: Move  Enter: Select  c: Custom  Esc: Cancel"
		DrawInstructions(instructions, 0, customOption+6)

		termbox.Flush()
	}

	drawTimeSelection()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if selected > 0 {
					selected--
					drawTimeSelection()
				}
			case termbox.KeyArrowDown:
				if selected < customOption {
					selected++
					drawTimeSelection()
				}
			case termbox.KeyEnter:
				termbox.Close()
				if selected == customOption {
					// User wants to enter custom time
					return "", true, nil
				}
				// Return selected time
				return suggestedTimes[selected], false, nil
			case termbox.KeyEsc:
				return "", false, fmt.Errorf("time selection cancelled")
			default:
				// Handle 'c' or 'C' key for custom time
				if ev.Ch == 'c' || ev.Ch == 'C' {
					termbox.Close()
					return "", true, nil
				}
			}
		case termbox.EventError:
			return "", false, ev.Err
		}
	}
}

// ShowCustomTimeInput creates interactive prompt for custom time entry
// Accepts formats: "HH:MM", "HH:MM AM/PM", "HHhMM"
// Returns formatted time string or error
func ShowCustomTimeInput() (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	ClearScreen()

	// Draw title
	title := "Enter custom commit time"
	DrawTitle(title, 0, 0)

	// Draw format instructions
	instructions := []string{
		"",
		"Accepted formats:",
		"  - HH:MM (24-hour): 18:45, 09:30",
		"  - HHhMM: 18h45, 9h30",
		"  - HH:MM AM/PM: 6:45 PM, 9:30 AM",
		"",
		"Time: ",
	}

	for i, line := range instructions {
		DrawInstructions(line, 0, i+2)
	}

	// Draw bottom instructions
	bottomInst := "Enter: Confirm  Esc: Cancel"
	DrawInstructions(bottomInst, 0, len(instructions)+4)

	// Input buffer
	input := []rune{}
	cursorPos := 0

	redraw := func(errorMsg string) {
		// Redraw prompt area
		ClearLine(len(instructions)+1, 80)

		prompt := "Time: "
		DrawInstructions(prompt, 0, len(instructions)+1)

		// Draw input
		DrawText(string(input), len(prompt), len(instructions)+1, ColorDefault, ColorDefault)

		// Draw error if any
		if errorMsg != "" {
			ClearLine(len(instructions)+3, 80)
			DrawText(errorMsg, 0, len(instructions)+3, ColorError, ColorDefault)
		}

		termbox.SetCursor(len(prompt)+cursorPos, len(instructions)+1)
		termbox.Flush()
	}

	redraw("")

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				// Validate and return input
				timeStr := string(input)
				if timeStr == "" {
					redraw("Error: Time cannot be empty")
					continue
				}

				hour, minute, err := timeutil.ParseTimeString(timeStr)
				if err != nil {
					redraw(fmt.Sprintf("Error: %v", err))
					continue
				}

				// Return formatted time
				return fmt.Sprintf("%02d:%02d", hour, minute), nil

			case termbox.KeyEsc:
				return "", fmt.Errorf("custom time input cancelled")

			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorPos > 0 {
					input = append(input[:cursorPos-1], input[cursorPos:]...)
					cursorPos--
					redraw("")
				}

			case termbox.KeyDelete:
				if cursorPos < len(input) {
					input = append(input[:cursorPos], input[cursorPos+1:]...)
					redraw("")
				}

			case termbox.KeyArrowLeft:
				if cursorPos > 0 {
					cursorPos--
					redraw("")
				}

			case termbox.KeyArrowRight:
				if cursorPos < len(input) {
					cursorPos++
					redraw("")
				}

			case termbox.KeySpace:
				input = append(input[:cursorPos], append([]rune{' '}, input[cursorPos:]...)...)
				cursorPos++
				redraw("")

			default:
				if ev.Ch != 0 {
					input = append(input[:cursorPos], append([]rune{ev.Ch}, input[cursorPos:]...)...)
					cursorPos++
					redraw("")
				}
			}

		case termbox.EventError:
			return "", ev.Err
		}
	}
}
