package utils

import (
	"context"
	"fmt"
	"log"
	"os"
)

const (
	Fatal = iota
	Error
	Warning
	Info
	Debug
	Trace
)

var Level = Info

type Logger interface {
	Logf(level int, format string, args ...any)
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
	Warningf(format string, args ...any)
	Infof(format string, args ...any)
	Debugf(format string, args ...any)
	Tracef(format string, args ...any)
}

func NewLogger(out *os.File, prefix string, showTimestamp bool) Logger {
	flag := log.Lshortfile | log.Lmsgprefix
	if showTimestamp {
		flag |= log.LstdFlags
	}
	prefix = fmt.Sprintf("[%s] ", prefix)
	return &logWrapper{out: out, logger: log.New(out, prefix, flag)}
}

type logWrapper struct {
	out    *os.File
	logger *log.Logger
}

func (l *logWrapper) Logf(level int, format string, args ...any) {
	if level > Level {
		return
	}
	l.logger.Printf(format, args...)
	err := l.out.Sync()
	// ignore Sync error, as stderr doesn't support it
	IgnoreErr(err)
}

func (l *logWrapper) Fatalf(format string, args ...any) {
	l.Logf(Fatal, format, args...)
}

func (l logWrapper) Errorf(format string, args ...any) {
	l.Logf(Error, format, args...)
}

func (l logWrapper) Warningf(format string, args ...any) {
	l.Logf(Warning, format, args...)
}

func (l logWrapper) Infof(format string, args ...any) {
	l.Logf(Info, format, args...)
}

func (l logWrapper) Debugf(format string, args ...any) {
	l.Logf(Debug, format, args...)
}

func (l logWrapper) Tracef(format string, args ...any) {
	l.Logf(Trace, format, args...)
}

func Logf(ctx context.Context, level int, format string, args ...any) {
	val := ctx.Value("logger")
	if logger, ok := val.(Logger); ok {
		logger.Logf(level, format, args...)
	}
}
