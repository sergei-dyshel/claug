package utils

import (
	"testing"

	emperror "emperror.dev/errors"

	"github.com/rotisserie/eris"
)

func f1() error {
	return Errf("some error")
}

func f2() error {
	return Wrapf(f1(), "f1() failed")
}

func f3() error {
	return Wrapf(f2(), "f2() failed")
}

func f4() error {
	return Wrapf(f3(), "f3() failed")
}

func f5() error {
	ch := make(chan error)
	go func() {
		ch <- f4()
	}()
	return WrapExtendf(<-ch, "received err from chan")
}

func TestErrors(t *testing.T) {
	for _, impl := range allErrImpl {
		t.Run(impl.Name(), func(t *testing.T) {
			curErrImpl = impl
			err := f5()
			t.Logf("Error only: %s", err)
			t.Logf("Error with stack trace: \n%s", ErrFull(err))
		})
	}
}

func TestEris(t *testing.T) {
	err := eris.New("some error")
	err = eris.Wrapf(err, "wrapped error")
	err = eris.Wrapf(err, "another error")
	ch := make(chan error)
	go func() {
		ch <- err
	}()

	t.Logf("%+v", eris.Wrapf(<-ch, "err received from chan"))
}

func TestEmperror(t *testing.T) {
	err := emperror.New("some error")
	err = emperror.WrapIff(err, "wrapped error")
	err = emperror.WrapIff(err, "another error")
	t.Logf("%+v", err)
}
