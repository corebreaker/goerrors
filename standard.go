package goerrors

import (
    "bytes"
    "fmt"
)

// Interface for a standard error which decorate another basic go error (`error` go interface)
// with an error code, message and other additionnal informations
type IStandardError interface {
    // Base interface
    IError

    // Add more informations on that error
    AddInfo(info string, args ...interface{}) IError

    // Get error code
    GetCode() int64
}

// Internal structure type for the standard error
type tStandardError struct {
    GoError

    code  int64        // Error code
    infos bytes.Buffer // Additionnal informations
}

// Add informations in standard error
func (self *tStandardError) AddInfo(info string, args ...interface{}) IStandardError {
    // Just prints into internal buffer the informations passed as parameters
    fmt.Fprintln(&self.infos, fmt.Sprintf(info, args...))

    return self
}

// «Decorate» the error passed as "err" parameter.
// The error returned will be an standard error with additionnal informations and stack trace.
func DecorateError(err error) IStandardError {
    res := new(tStandardError)
    res.Init(res, "", nil, err, -1)

    return res
}

// Like `DecorateError` with error code and custom data
func DecorateErrorWithInfos(err error, error_code int64, custom_data interface{}) IStandardError {
    res := &tStandardError{code: error_code}
    res.Init(res, "", nil, err, -1)

    return res
}

// Make an standard error from a message passed as "message" parameter
func MakeError(message string, args ...interface{}) IStandardError {
    res := new(tStandardError)
    res.Init(fmt.Sprintf(message, args...), "", nil, err, -1)

    return res
}

// Like `MakeError` with error code and custom data
func MakeErrorWithInfos(error_code int64, data interface{}, message string, args ...interface{}) IStandardError {
    res := &tStandardError{code: error_code}
    res.Init(fmt.Sprintf(message, args...), "", nil, err, -1)

    return res
}

// Global function to add information in an error whatever.
// This function just call the "AddInfo" method of an standard error.
func AddInfo(err error, info string, args ...interface{}) IStandardError {
    // Check if "err" is already an standard error
    go_err, ok := err.(IStandardError)
    if !ok {
        // Otherwise decorate that unknown error
        go_err = DecorateError(err)
    }

    // Delegate call to "AddInfo" method
    return go_err.AddInfo(info, args...)
}

func (self *GoError) Catch(err *error, catch, finally ErrorHandler) {
    recovered := recover()
    res_err, ok := recovered.(error)
    if !ok {
        panic(recovered)
    }

    defer func() {
        if finally != nil {
            ierr, _ := res_err.(IError)

            res_err = finally(ierr)
        }

        if err != nil {
            *err = res_err
        }
    }()

    this := self.get_reference()
    if !this.IsParentOf(res_err) {
        if err == nil {
            panic(recovered)
        }

        return
    }

    if catch != nil {
        res_err = catch(res_err.(IError))
    }
}

func Try(try, catch, finally ErrorHandler) error {
    defer func() {

    }

    return try(nil)
}

func Raise(message string, args ...interface{}) {
    panic(MakeError(message, args...))
}

func RaiseWithInfos(error_code int64, data interface{}, message string, args ...interface{}) {
    panic(MakeErrorWithInfos(error_code, data, message, args...))
}

func RaiseError(err error) {
    gerr, ok := err.(IError)
    if !ok {
        gerr = DecorateError(err)
    }

    panic(err)
}
