package utils

import "fmt"

func Assertf(condition bool, format string, args ...any) {
	if !condition {
		Panicf(format, args...)
	}
}

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func AssertErrf(err error, format string, args ...any) {
	if err != nil {
		Panicf(format, args...)
	}
}

func AssertErr(err error) {
	AssertErrf(err, "Unexpected error: %v", err)
}

func IgnoreErr(_ error) {
	// do nothing
	// errcheck should be pleased
}
