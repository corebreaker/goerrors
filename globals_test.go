package goerrors

import (
	"testing"
)

func TestUncatchedErrorHandlerCall(t *testing.T) {
	logFatal = func(v ...interface{}) {}

	_ = uncatchedErrorHandler(nil)
	_ = uncatchedErrorHandler(MakeError("error"))
}

func TestSetUncatchedErrorHandler(t *testing.T) {
	SetUncatchedErrorHandler(func(err IError) error { return err })
}

func TestCheckedMain(t *testing.T) {
	logFatal = func(v ...interface{}) {}

	CheckedMain(func() error {
		return nil
	})

	CheckedMain(func() error {
		return MakeError("error")
	})

	SetUncatchedErrorHandler(func(err IError) error { return nil })
	CheckedMain(func() error {
		Raise("error")

		return nil
	})
}
