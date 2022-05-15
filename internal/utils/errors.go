package utils

import (
	"fmt"

	emperror "emperror.dev/errors"
	"github.com/joomcode/errorx"
	"github.com/rotisserie/eris"
)

var (
	ErrEris     = &errEris{}
	ErrEmperror = &errEmperror{}
	ErrErrorx   = &errErrorx{}

	allErrImpl         = []errImpl{ErrEris, ErrEmperror, ErrErrorx}
	curErrImpl errImpl = ErrErrorx
)

func ErrFull(err error) string {
	return curErrImpl.Full(err)
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return curErrImpl.Wrap(err, fmt.Sprintf(format, args...))
}

func Errf(format string, args ...any) error {
	return curErrImpl.New(fmt.Sprintf(format, args...))
}

func WrapExtendf(err error, format string, args ...any) error {
	// TODO: add to unitest
	if err == nil {
		return nil
	}
	return curErrImpl.WrapExtend(err, fmt.Sprintf(format, args...))
}

// errImpl provides adapter for particular errors implementation
type errImpl interface {
	// Name of implementation, in CamelCase
	Name() string

	// Wrap wraps error with message and adds stack frame
	Wrap(err error, msg string) error

	// WrapExtend wraps error with message and adds new stack trace.
	// Used for wrapping errors received from another gourtine
	WrapExtend(err error, msg string) error

	// New ceates new error from message and captures stack trace
	New(msg string) error

	// Full formats error message with stack trace
	Full(err error) string
}

type errEris struct{}

func (*errEris) Name() string {
	return "Eris"
}

func (*errEris) Wrap(err error, msg string) error {
	return eris.Wrap(err, msg)
}

func (e *errEris) WrapExtend(err error, msg string) error {
	return e.Wrap(err, msg)
}

func (*errEris) New(msg string) error {
	return eris.New(msg)
}

func (*errEris) Full(err error) string {
	return fmt.Sprintf("%+v", err)
}

type errErrorx struct{}

var (
	errorxNamespace = errorx.NewNamespace("")
	errorxType      = errorxNamespace.NewType("")
)

func (*errErrorx) Name() string {
	return "Errorx"
}

func (*errErrorx) Wrap(err error, msg string) error {
	return errorx.Decorate(err, msg)
}

func (*errErrorx) WrapExtend(err error, msg string) error {
	return errorx.EnhanceStackTrace(err, msg)
}

func (*errErrorx) New(msg string) error {
	return errorxType.New(msg)
}

func (*errErrorx) Full(err error) string {
	return fmt.Sprintf("%+v", err)
}

type errEmperror struct{}

func (*errEmperror) Name() string {
	return "Emprerror"
}

func (*errEmperror) Wrap(err error, msg string) error {
	return emperror.WrapIf(err, msg)
}

func (e *errEmperror) WrapExtend(err error, msg string) error {
	return emperror.Wrap(err, msg)
}

func (*errEmperror) New(msg string) error {
	return emperror.New(msg)
}

func (*errEmperror) Full(err error) string {
	return fmt.Sprintf("%+v", err)
}
