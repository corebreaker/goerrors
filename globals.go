package goerrors

import (
	"fmt"
	"log"
)

var (
	logFatal                           = log.Fatal
	uncatchedErrorHandler ErrorHandler = func(err IError) error {
		if err != nil {
			logFatal(err)
		}

		return nil
	}
)

// MainHandler is the handler for main function used with the CheckedMain function
type MainHandler func() error

// CheckedMain should be called in the main function body.
// It assumes the uncatched error to see errors with stacktraces.
func CheckedMain(handler MainHandler) {
	defer func() {
		recovered := recover()
		if (recovered == nil) || (uncatchedErrorHandler == nil) {
			return
		}

		var err error

		recoveredError, ok := recovered.(error)
		if ok {
			err = recoveredError
		} else {
			err = fmt.Errorf("error: %s", recovered)
		}

		ierr, ok := err.(IError)
		if !ok {
			ierr = DecorateError(err)
		}

		cerr := uncatchedErrorHandler(ierr)
		if cerr != nil {
			logFatal(cerr)
		}
	}()

	err := handler()

	if err != nil {
		logFatal(err)
	}
}

// SetUncatchedErrorHandler defines the handler called when an error is not catched.
//  But the handled is called only if the CheckedMain function is called (indeed, the handler is called by CheckedMain).
func SetUncatchedErrorHandler(handler ErrorHandler) ErrorHandler {
	oldHandler := uncatchedErrorHandler

	uncatchedErrorHandler = handler

	return oldHandler
}

// DiscardPanic discard a panic.
//  This is an helper function that should be called with a defer instruction to blocks a panic and permits to
//  continuing the excution instead of exiting the program.
func DiscardPanic() {
	recover()
}
