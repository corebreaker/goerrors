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

type MainHandler func() error

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

func SetUncatchedErrorHandler(handler ErrorHandler) ErrorHandler {
	old_handler := uncatchedErrorHandler

	uncatchedErrorHandler = handler

	return old_handler
}

func DiscardPanic() {
	recover()
}
