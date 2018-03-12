package goerrors

import "testing"

func TestSetUncatchedErrorHandler(t *testing.T) {
	SetUncatchedErrorHandler(func(err IError) error { return err })
}

func TestCheckedMain(t *testing.T) {
	CheckedMain(func() error {
		return nil
	})
}
