package ui

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

// ShowMessageSelector displays an interactive UI for selecting commit messages
// Returns the selected/edited message or an error
func ShowMessageSelector(messages []string, prompt string) (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	selected := 0
	drawMessageList(messages, selected, prompt, false)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if selected > 0 {
					selected--
					drawMessageList(messages, selected, prompt, false)
				}
			case termbox.KeyArrowDown:
				if selected < len(messages) {
					selected++
					drawMessageList(messages, selected, prompt, false)
				}
			case termbox.KeyEnter:
				termbox.Close()
				if selected == len(messages) {
					// Custom message input
					msg, err := ShowMessageEditor("")
					if err != nil {
						if err.Error() == "edit cancelled" {
							// Reopen the selection screen if edit was cancelled
							err = termbox.Init()
							if err != nil {
								return "", err
							}
							drawMessageList(messages, selected, prompt, false)
							continue
						}
						return "", err
					}
					return msg, nil
				}

				// Edit selected message
				editedMsg, err := ShowMessageEditor(messages[selected])
				if err != nil {
					if err.Error() == "edit cancelled" {
						// Reopen the selection screen if edit was cancelled
						err = termbox.Init()
						if err != nil {
							return "", err
						}
						drawMessageList(messages, selected, prompt, false)
						continue
					}
					return "", err
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

// drawMessageList renders the message selection UI
func drawMessageList(messages []string, selected int, prompt string, showEditPrompt bool) {
	ClearScreen()

	// Get terminal width
	width, _ := GetTerminalSize()

	// Draw title
	if prompt == "" {
		prompt = "Select a commit message:"
	}
	DrawTitle(prompt, 0, 0)

	// Draw messages
	for i, msg := range messages {
		fg := ColorDefault
		bg := ColorDefault
		if i == selected {
			fg = ColorSelected
			bg = ColorSelectedBG
		}

		// Draw bullet point
		DrawBullet(2, i+2, i == selected, fg, bg)

		// For detailed messages, show only the subject line (first line)
		displayMsg := msg
		if strings.Contains(msg, "\n") {
			// Extract just the subject line for display
			lines := strings.Split(msg, "\n")
			displayMsg = lines[0]
		}

		// Truncate message if it's too long for the terminal
		maxMsgWidth := width - 6 // Account for bullet point and spacing
		displayMsg = TruncateText(displayMsg, maxMsgWidth)

		// Draw message
		DrawText(displayMsg, 4, i+2, fg, bg)
	}

	// Draw custom message option
	fg := ColorDefault
	bg := ColorDefault
	if selected == len(messages) {
		fg = ColorSelected
		bg = ColorSelectedBG
	}
	customMsg := "Edit custom message"
	DrawBullet(2, len(messages)+2, selected == len(messages), fg, bg)
	DrawText(customMsg, 4, len(messages)+2, fg, bg)

	// Draw instructions
	instructions := "↑↓: Move  Enter: Select  Esc: Cancel"
	DrawInstructions(instructions, 0, len(messages)+4)

	// Draw edit prompt if needed
	if showEditPrompt {
		prompt := "Edit message (press Enter to confirm):"
		DrawTitle(prompt, 0, len(messages)+6)
	}

	termbox.Flush()
}
