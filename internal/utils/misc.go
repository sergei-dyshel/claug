package utils

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/tymonx/go-formatter/formatter"

	"github.com/mitchellh/go-homedir"
)

func shellExpand(cmd string) string {
	exp, err := homedir.Expand(cmd)
	AssertErr(err)
	return os.ExpandEnv(exp)
}

func HomeDir() string {
	// TODO: consider `homedir.Dir()` as more versatile
	dirname, err := os.UserHomeDir()
	AssertErr(err)
	return dirname
}

func Doc(text string) string {
	text = strings.ReplaceAll(text, "''", "`")
	return strings.Join(Map(strings.Split(strings.TrimSpace(text), "\n"), strings.TrimSpace), " ")
}

func Docf(format string, args ...any) string {
	return Doc(fmt.Sprintf(format, args...))
}

func DocFmt(text string, args ...any) string {
	return Doc(formatter.MustFormat(text, args...))
}

func Indent(text string, indent string) {
}

func If[T any](cond bool, then, else_ T) T {
	if cond {
		return then
	}
	return else_
}
