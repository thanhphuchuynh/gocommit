package ui

import (
	"github.com/nsf/termbox-go"
)

// Common UI constants
const (
	// Colors
	ColorTitle        = termbox.ColorYellow
	ColorInstructions = termbox.ColorCyan
	ColorSelected     = termbox.ColorBlack
	ColorSelectedBG   = termbox.ColorGreen
	ColorDefault      = termbox.ColorDefault
	ColorError        = termbox.ColorRed

	// Characters
	BulletDefault  = "•"
	BulletSelected = "→"

	// Default terminal dimensions
	DefaultWidth  = 80
	DefaultHeight = 24
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	termbox.Clear(ColorDefault, ColorDefault)
}

// DrawTitle draws a title at the specified position
func DrawTitle(text string, x, y int) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, ColorTitle, ColorDefault)
	}
}

// DrawInstructions draws instruction text at the specified position
func DrawInstructions(text string, x, y int) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, ColorInstructions, ColorDefault)
	}
}

// DrawText draws text at the specified position with given colors
func DrawText(text string, x, y int, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

// DrawBullet draws a bullet point at the specified position
func DrawBullet(x, y int, isSelected bool, fg, bg termbox.Attribute) {
	bullet := BulletDefault
	if isSelected {
		bullet = BulletSelected
	}
	termbox.SetCell(x, y, []rune(bullet)[0], fg, bg)
}

// GetTerminalSize returns the terminal width and height, with defaults if too small
func GetTerminalSize() (width, height int) {
	width, height = termbox.Size()
	if width < 10 {
		width = DefaultWidth
	}
	if height < 10 {
		height = DefaultHeight
	}
	return width, height
}

// ClearLine clears a line at the specified y position
func ClearLine(y int, width int) {
	for i := 0; i < width; i++ {
		termbox.SetCell(i, y, ' ', ColorDefault, ColorDefault)
	}
}

// TruncateText truncates text to fit within maxWidth, adding "..." if needed
func TruncateText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	if maxWidth <= 3 {
		return text[:maxWidth]
	}
	return text[:maxWidth-3] + "..."
}
