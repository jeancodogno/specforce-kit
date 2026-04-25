package tui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsTTY returns true if the standard output is a terminal.
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// GetTerminalWidth returns the current width of the terminal.
func GetTerminalWidth() int {
	// Simplified implementation, could use golang.org/x/term
	return 80 // Default fallback
}
