package goerrors

import "log"

var (
	uncatchedErrorHandler ErrorHandler = func(err IError) error {
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
)

type MainHandler func() error

func CheckedMain(handler MainHandler) {
	defer func() {
		var sent_err error = nil

		err := new(GoError)
		err.Init(err, "", nil, nil, -1)

		err.Catch(&sent_err, uncatchedErrorHandler, nil)

		if sent_err != nil {
			log.Fatal(sent_err)
		}
	}()

	err := handler()

	if err != nil {
		log.Fatal(err)
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
