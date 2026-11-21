package tui

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
)

// IsTTY returns true if the standard output is a terminal.
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// defaultTerminalWidth is used when the width of stdout cannot be determined.
const defaultTerminalWidth = 80

// GetTerminalWidth returns the current width of the terminal, or 80 when it
// cannot be determined (for example when stdout is not a terminal).
func GetTerminalWidth() int {
	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || width <= 0 {
		return defaultTerminalWidth
	}
	return width
}
