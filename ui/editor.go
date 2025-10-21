package ui

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

// ShowMessageEditor displays an interactive editor for commit messages
// Returns the edited message or an error
func ShowMessageEditor(initialMsg string) (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	// Get terminal width
	width, _ := GetTerminalSize()

	// Clear screen and show edit prompt
	ClearScreen()

	// Draw title
	title := "Edit commit message (Enter to confirm, Shift+Enter for new line, Esc to cancel):"
	DrawTitle(title, 0, 0)

	// Draw instructions
	instructions := "Use arrow keys to move, Shift+Enter for new line, backspace/delete to edit"
	DrawInstructions(instructions, 0, 1)

	// Initialize message buffer and cursor
	editedMsg := []rune(initialMsg)
	cursorPos := len(editedMsg)
	scrollX := 0
	maxScroll := 0
	currentLine := 0

	// Function to redraw the message with wrapping and scrolling
	redraw := func() {
		ClearScreen()

		// Get terminal height for scrolling
		_, height := GetTerminalSize()

		// Draw title and instructions
		DrawTitle(title, 0, 0)
		DrawInstructions(instructions, 0, 1)

		// Split message into display lines (considering word wrapping)
		msgLines := strings.Split(string(editedMsg), "\n")
		displayLines := []string{}
		lineToOriginal := []int{} // Maps display line to original line number

		for lineIdx, msgLine := range msgLines {
			if len(msgLine) <= width-1 {
				// Line fits in one display line
				displayLines = append(displayLines, msgLine)
				lineToOriginal = append(lineToOriginal, lineIdx)
			} else {
				// Wrap long lines
				for i := 0; i < len(msgLine); i += width - 1 {
					end := i + width - 1
					if end > len(msgLine) {
						end = len(msgLine)
					}
					displayLines = append(displayLines, msgLine[i:end])
					lineToOriginal = append(lineToOriginal, lineIdx)
				}
			}
		}

		// Calculate cursor position in display coordinates
		linesBeforeCursor := strings.Split(string(editedMsg[:cursorPos]), "\n")
		cursorOriginalLine := len(linesBeforeCursor) - 1
		lastLine := linesBeforeCursor[len(linesBeforeCursor)-1]
		cursorColInOriginal := len(lastLine)

		// Find cursor position in display coordinates
		cursorDisplayLine := 0
		cursorDisplayCol := 0
		for i, origLine := range lineToOriginal {
			if origLine == cursorOriginalLine {
				lineStart := 0
				for j := 0; j < i; j++ {
					if lineToOriginal[j] == cursorOriginalLine {
						lineStart += len(displayLines[j])
					}
				}
				if cursorColInOriginal >= lineStart && cursorColInOriginal <= lineStart+len(displayLines[i]) {
					cursorDisplayLine = i
					cursorDisplayCol = cursorColInOriginal - lineStart
					break
				}
			}
		}

		// Calculate scroll offset to keep cursor visible
		maxDisplayLines := height - 3 // Account for title and instructions
		scrollY := 0
		if cursorDisplayLine >= maxDisplayLines {
			scrollY = cursorDisplayLine - maxDisplayLines + 1
		}

		// Draw visible lines
		for i := scrollY; i < len(displayLines) && i-scrollY < maxDisplayLines; i++ {
			line := displayLines[i]
			y := 2 + (i - scrollY)

			// Draw each character of the line
			for j, ch := range line {
				if j < width-1 {
					termbox.SetCell(j, y, ch, ColorDefault, ColorDefault)
				}
			}
		}

		// Set cursor position (adjusted for scroll)
		if cursorDisplayLine >= scrollY && cursorDisplayLine < scrollY+maxDisplayLines {
			cursorY := 2 + (cursorDisplayLine - scrollY)
			if cursorDisplayCol < width-1 {
				termbox.SetCursor(cursorDisplayCol, cursorY)
			}
		}

		termbox.Flush()
	}

	redraw()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				if ev.Mod&termbox.ModAlt != 0 {
					// Insert newline when Shift+Enter is pressed
					editedMsg = append(editedMsg[:cursorPos], append([]rune{'\n'}, editedMsg[cursorPos:]...)...)
					cursorPos++
					redraw()
				} else {
					// Regular Enter confirms the edit
					return string(editedMsg), nil
				}
			case termbox.KeyEsc:
				return "", fmt.Errorf("edit cancelled")
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorPos > 0 {
					editedMsg = append(editedMsg[:cursorPos-1], editedMsg[cursorPos:]...)
					cursorPos--
					if cursorPos < scrollX {
						scrollX = cursorPos
					}
				}
			case termbox.KeyDelete:
				if cursorPos < len(editedMsg) {
					editedMsg = append(editedMsg[:cursorPos], editedMsg[cursorPos+1:]...)
				}
			case termbox.KeyArrowLeft:
				if cursorPos > 0 {
					cursorPos--
					if cursorPos < scrollX {
						scrollX = cursorPos
					}
				}
			case termbox.KeyArrowRight:
				if cursorPos < len(editedMsg) {
					cursorPos++
					if cursorPos >= scrollX+width-1 {
						scrollX = cursorPos - width + 2
					}
				}
			case termbox.KeyArrowUp:
				// Move cursor up one line
				lines := strings.Split(string(editedMsg[:cursorPos]), "\n")
				if len(lines) > 1 {
					currentLine = len(lines) - 2
					prevLine := lines[currentLine]
					if len(prevLine) < cursorPos-len(lines[len(lines)-1])-1 {
						cursorPos = len(strings.Join(lines[:currentLine+1], "\n")) + 1
					} else {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1 + len(prevLine)
					}
					redraw()
				}
			case termbox.KeyArrowDown:
				// Move cursor down one line
				lines := strings.Split(string(editedMsg[:cursorPos]), "\n")
				if currentLine < len(lines)-1 {
					currentLine++
					nextLine := lines[currentLine]
					if len(nextLine) < cursorPos-len(strings.Join(lines[:currentLine], "\n"))-1 {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1 + len(nextLine)
					} else {
						cursorPos = len(strings.Join(lines[:currentLine], "\n")) + 1
					}
					redraw()
				}
			case termbox.KeySpace:
				editedMsg = append(editedMsg[:cursorPos], append([]rune{' '}, editedMsg[cursorPos:]...)...)
				cursorPos++
				if cursorPos >= scrollX+width-1 {
					scrollX = cursorPos - width + 2
				}
			default:
				if ev.Ch != 0 {
					editedMsg = append(editedMsg[:cursorPos], append([]rune{ev.Ch}, editedMsg[cursorPos:]...)...)
					cursorPos++
					if cursorPos >= scrollX+width-1 {
						scrollX = cursorPos - width + 2
					}
				}
			}

			// Update max scroll if needed
			if len(editedMsg) > width-1 {
				maxScroll = len(editedMsg) - width + 1
			} else {
				maxScroll = 0
			}
			if scrollX > maxScroll {
				scrollX = maxScroll
			}

			redraw()

		case termbox.EventError:
			return "", ev.Err
		}
	}
}
