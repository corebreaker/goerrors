package goerrors

import (
	"bytes"
	"fmt"
)

// IStandardError Interface for a standard error which decorate another basic go error (`error` go interface)
// with an error code, message and other additionnal informations
type IStandardError interface {
	// Base interface
	IError

	// Add more informations on that error
	AddInfo(info string, args ...interface{}) IStandardError

	// Get error code
	GetCode() int64
}

// Internal structure type for the standard error
type tStandardError struct {
	GoError

	code  int64        // Error code
	infos bytes.Buffer // Additionnal informations
}

// GetName gets standard error name
func (se *tStandardError) GetName() string {
	return "StandardError"
}

// AddInfo adds informations in standard error
func (se *tStandardError) AddInfo(info string, args ...interface{}) IStandardError {
	// Just prints into internal buffer the informations passed as parameters
	_, _ = fmt.Fprintln(&se.infos, fmt.Sprintf(info, args...))

	return se
}

// GetCode gets error code
func (se *tStandardError) GetCode() int64 {
	return se.code
}

// DecorateError «decorates» the error passed as "err" parameter.
// The error returned will be an standard error with additionnal informations and stack trace.
func DecorateError(err error) IStandardError {
	if err == nil {
		return nil
	}

	ierr, ok := err.(IStandardError)
	if !ok {
		res := new(tStandardError)
		_ = res.Init(res, "", nil, err, 1)

		ierr = res
	}

	return ierr
}

// DecorateErrorWithDatas is like `DecorateError` with error code and custom data
func DecorateErrorWithDatas(err error, code int64, data interface{}, msg string, args ...interface{}) IStandardError {
	if err == nil {
		return nil
	}

	ierr, ok := err.(IStandardError)
	if ok {
		_ = ierr.AddInfo("Recorate for code=%d and message=%s", code, fmt.Sprintf(msg, args...))
	} else {
		res := &tStandardError{code: code}
		_ = res.Init(res, fmt.Sprintf(msg, args...), data, err, 1)

		ierr = res
	}

	return ierr
}

// MakeError makes an standard error from a message passed as "message" parameter
func MakeError(message string, args ...interface{}) IStandardError {
	res := new(tStandardError)
	_ = res.Init(res, fmt.Sprintf(message, args...), nil, nil, 1)

	return res
}

// MakeErrorWithDatas is like `MakeError` with error code and custom data
func MakeErrorWithDatas(code int64, data interface{}, message string, args ...interface{}) IStandardError {
	res := &tStandardError{code: code}
	_ = res.Init(res, fmt.Sprintf(message, args...), data, nil, 1)

	return res
}

// AddInfo is the global function to add information in an error whatever.
// This function just call the "AddInfo" method of an standard error.
func AddInfo(err error, info string, args ...interface{}) IStandardError {
	// Check if "err" is already an standard error
	goErr, ok := err.(IStandardError)
	if !ok {
		// Otherwise decorate that unknown error
		goErr = DecorateError(err)
	}

	// Delegate call to "AddInfo" method
	return goErr.AddInfo(info, args...)
}

// Catch is the global function to catch an error
func Catch(err *error, catch, finally ErrorHandler) {
	var resErr error

	defer func() {
		if finally != nil {
			ierr, _ := resErr.(IError)

			resErr = finally(ierr)
		}

		if err != nil {
			*err = resErr
		}
	}()

	recovered := recover()
	if recovered == nil {
		return
	}

	var ok bool

	resErr, ok = recovered.(error)
	if !ok {
		panic(recovered)
	}

	if catch == nil {
		return
	}

	ierr, ok := resErr.(IError)
	if !ok {
		ierr = DecorateError(resErr)
	}

	cerr := catch(ierr)
	if cerr != nil {
		resErr = cerr
	}
}

// Try is the global function to call a try block
func Try(try, catch, finally ErrorHandler) (err error) {
	defer func() {
		if finally == nil {
			return
		}

		ierr, _ := err.(IError)

		ferr := finally(ierr)
		if ferr != nil {
			err = ferr
		}
	}()

	defer func() {
		recovered := recover()
		if ((recovered == nil) || (catch == nil)) && (err == nil) {
			return
		}

		if err == nil {
			recoveredError, ok := recovered.(error)
			if ok {
				err = recoveredError
			} else {
				err = fmt.Errorf("error: %s", recovered)
			}
		}

		ierr, ok := err.(IError)
		if !ok {
			ierr = DecorateError(err)
		}

		cerr := catch(ierr)
		if cerr != nil {
			err = cerr
		}
	}()

	return try(nil)
}

// Raise is the global function to raise an anonymous error
func Raise(message string, args ...interface{}) {
	MakeError(message, args...).raise(1)
}

// RaiseWithInfos is like Raise with error code and custom data
func RaiseWithInfos(errorCode int64, data interface{}, message string, args ...interface{}) {
	MakeErrorWithDatas(errorCode, data, message, args...).raise(1)
}

// RaiseError is the global funtion to raise the error passed in argument
func RaiseError(err error) {
	gerr, ok := err.(IError)
	if !ok {
		gerr = DecorateError(err)
	}

	gerr.raise(1)
}
