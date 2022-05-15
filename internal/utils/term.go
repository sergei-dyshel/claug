package utils

import (
	"os"

	"golang.org/x/term"
)

func IsTerminal(file *os.File) bool {
	return term.IsTerminal(int(file.Fd()))
}

func TerminalWidth(file *os.File) int {
	width, _, err := term.GetSize(int(file.Fd()))
	AssertErr(err)
	return width
}
