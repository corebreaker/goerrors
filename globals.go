package goerrors

import "log"

var (
    uncatched_error_handler ErrorHandler = func(err IError) error {
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

        err.Catch(&sent_err, uncatched_error_handler, nil)

        if sent_err != nil {
            log.Fatal(sent_err)
        }
    }()

    res := handler()

    if res != nil {
        log.Fatal(err)
    }
}

func SetUncatchedErrorHandler(handler ErrorHandler) ErrorHandler {
    old_handler := uncatched_error_handler

    uncatched_error_handler = handler

    return old_handler
}
