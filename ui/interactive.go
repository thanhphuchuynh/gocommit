package ui

import (
	"strings"

	"github.com/nsf/termbox-go"
)

// InteractiveAction represents the user's choice in interactive mode
type InteractiveAction int

const (
	ActionAccept InteractiveAction = iota
	ActionEdit
	ActionRegenerate
	ActionShowDiff
	ActionQuit
)

// InteractivePromptResult contains the action and any modified message
type InteractivePromptResult struct {
	Action  InteractiveAction
	Message string // Updated message if edited
}

// ShowInteractivePrompt displays an interactive prompt after message selection
// Returns the action the user wants to take
func ShowInteractivePrompt(message, diff string) (*InteractivePromptResult, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	defer termbox.Close()

	showingDiff := false
	drawInteractivePrompt(message, diff, showingDiff)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'a', 'A':
				return &InteractivePromptResult{
					Action:  ActionAccept,
					Message: message,
				}, nil
			case 'e', 'E':
				termbox.Close()
				editedMsg, err := ShowMessageEditor(message)
				if err != nil {
					if err.Error() == "edit cancelled" {
						// Reopen interactive prompt if edit was cancelled
						err = termbox.Init()
						if err != nil {
							return nil, err
						}
						drawInteractivePrompt(message, diff, showingDiff)
						continue
					}
					return nil, err
				}
				return &InteractivePromptResult{
					Action:  ActionAccept,
					Message: editedMsg,
				}, nil
			case 'r', 'R':
				return &InteractivePromptResult{
					Action:  ActionRegenerate,
					Message: message,
				}, nil
			case 'd', 'D':
				showingDiff = !showingDiff
				drawInteractivePrompt(message, diff, showingDiff)
			case 'q', 'Q':
				return &InteractivePromptResult{
					Action:  ActionQuit,
					Message: message,
				}, nil
			}

			switch ev.Key {
			case termbox.KeyEsc:
				return &InteractivePromptResult{
					Action:  ActionQuit,
					Message: message,
				}, nil
			case termbox.KeyEnter:
				// Enter key accepts by default
				return &InteractivePromptResult{
					Action:  ActionAccept,
					Message: message,
				}, nil
			}
		case termbox.EventError:
			return nil, ev.Err
		}
	}
}

// drawInteractivePrompt renders the interactive prompt UI
func drawInteractivePrompt(message, diff string, showingDiff bool) {
	ClearScreen()

	width, height := GetTerminalSize()
	y := 0

	// Draw title
	DrawTitle("🤖 Generated Commit Message", 0, y)
	y += 2

	if showingDiff {
		// Split screen: diff on left, message on right
		splitX := width / 2

		// Draw diff on the left
		DrawText("Diff:", 2, y, ColorDefault, ColorDefault)
		y++

		diffLines := strings.Split(diff, "\n")
		maxDiffLines := height - y - 6 // Leave space for instructions
		if len(diffLines) > maxDiffLines {
			diffLines = diffLines[:maxDiffLines]
		}

		for _, line := range diffLines {
			// Truncate diff line if too long
			truncated := TruncateText(line, splitX-4)

			// Color diff lines
			fg := ColorDefault
			if strings.HasPrefix(line, "+") {
				fg = termbox.ColorGreen
			} else if strings.HasPrefix(line, "-") {
				fg = termbox.ColorRed
			} else if strings.HasPrefix(line, "@@") {
				fg = termbox.ColorCyan
			}

			DrawText(truncated, 2, y, fg, ColorDefault)
			y++
		}

		// Draw message on the right
		y = 3 // Reset y for message column
		DrawText("Message:", splitX+2, y, ColorDefault, ColorDefault)
		y++

		messageLines := strings.Split(message, "\n")
		for _, line := range messageLines {
			truncated := TruncateText(line, width-splitX-4)
			DrawText(truncated, splitX+2, y, ColorSelected, ColorDefault)
			y++
		}
	} else {
		// Show message only (full width)
		messageLines := strings.Split(message, "\n")

		// Draw message box
		boxWidth := width - 4
		DrawText("┌"+strings.Repeat("─", boxWidth-2)+"┐", 2, y, ColorDefault, ColorDefault)
		y++

		for _, line := range messageLines {
			if line == "" {
				DrawText("│"+strings.Repeat(" ", boxWidth-2)+"│", 2, y, ColorDefault, ColorDefault)
			} else {
				// Pad line to box width
				paddedLine := line
				if len(line) < boxWidth-4 {
					paddedLine = line + strings.Repeat(" ", boxWidth-4-len(line))
				} else {
					paddedLine = TruncateText(line, boxWidth-4)
				}
				DrawText("│ "+paddedLine+" │", 2, y, ColorSelected, ColorDefault)
			}
			y++
		}

		DrawText("└"+strings.Repeat("─", boxWidth-2)+"┘", 2, y, ColorDefault, ColorDefault)
		y += 2
	}

	// Position instructions at the bottom
	instructionY := height - 4

	// Draw menu options
	DrawText("What would you like to do?", 2, instructionY, ColorDefault, ColorDefault)
	instructionY++

	options := "[A]ccept  [E]dit  [R]egenerate  [D]iff  [Q]uit"
	DrawText(options, 2, instructionY, termbox.ColorCyan|termbox.AttrBold, ColorDefault)

	termbox.Flush()
}
