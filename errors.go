package goerrors

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	// Inheritance cache
	errorHierarchies = make(map[string][]string)
)

// Interface for extended Go errors
type IError interface {
	// This is an error
	error

	// Get error name
	GetName() string

	// Get original error
	GetSource() error

	// Get error message
	GetMessage() string

	// Get custon data
	GetData() interface{}

	// Complete try/catch/finally block
	Try(try, catch, finally ErrorHandler) error

	// Catch error (used as a defered call)
	Catch(err *error, catch, finally ErrorHandler)

	// Raise error
	Raise()

	// Test if this error is one of parents of error `err` passed in parameter
	IsParentOf(err error) bool

	// Get the real reference on this error
	getReference() IError

	// This method construct the stack trace only in 'Debug' Mode
	populateStackTrace(pruneLevels uint)

	// Get type of this error
	getParents() []string

	// Raise error with pruned levels
	raise(pruneLevels uint)
}

// Handler for executing Try, Catch, or Finally block
type ErrorHandler func(err IError) error

// Basic error structure
type GoError struct {
	source  error        // Cause or original error
	message string       // Error message
	trace   []string     // Stack trace
	data    interface{}  // Custom data
	errType reflect.Type // Type of this error
}

// Standard method of `error` interface
func (goErr *GoError) Error() string {
	var out bytes.Buffer

	err := goErr.getReference()

	// Prints error name
	_, _ = fmt.Fprintf(&out, "%s: ", err.GetName())

	// Get informations
	message := err.GetMessage()
	source := err.GetSource()
	data := err.GetData()

	// Prints error informations
	if message != "" {
		_, _ = fmt.Fprintln(&out, message)

		if data != nil {
			_, _ = fmt.Fprintln(&out, data)
		}

		if source != nil {
			_, _ = fmt.Fprintln(&out)
			_, _ = fmt.Fprintln(&out, "Source:", source)
		}
	} else {
		if source != nil {
			_, _ = fmt.Fprintln(&out, source)
		}

		if data != nil {
			_, _ = fmt.Fprintln(&out, data)
		}
	}

	// Prints stack trace only in debug mode
	if errDebug {
		for _, entry := range goErr.trace {
			_, _ = fmt.Fprintln(&out, "   ", entry)
		}

		// Prints a separator if stack trace is not empty
		if len(goErr.trace) > 0 {
			const sep = "------------------------------------------------------------------------------"

			_, _ = fmt.Fprintln(&out, sep)
		}
	}

	// Return content of the buffer resulting from printing theses informations
	return out.String()
}

// Get error name
func (goErr *GoError) GetName() string {
	return goErr.errType.PkgPath() + "." + goErr.errType.Name()
}

// Get cause error (parent error)
func (goErr *GoError) GetSource() error {
	return goErr.source
}

// Get error message
func (goErr *GoError) GetMessage() string {
	return goErr.message
}

// Get custon data
func (goErr *GoError) GetData() interface{} {
	return goErr.data
}

// Complete try/catch/finally block
func (goErr *GoError) Try(try, catch, finally ErrorHandler) (err error) {
	defer goErr.Catch(&err, catch, finally)

	return try(goErr.getReference())
}

// Catch error (used as a defered call)
func (goErr *GoError) Catch(err *error, catch, finally ErrorHandler) {
	var resErr error = nil

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

	if resErr, ok = recovered.(error); !ok {
		panic(recovered)
	}

	if this := goErr.getReference(); !this.IsParentOf(resErr) {
		if err == nil {
			panic(recovered)
		}

		return
	}

	if catch != nil {
		resErr = catch(resErr.(IError))
	}
}

// Raise error
func (goErr *GoError) Raise() {
	goErr.raise(1)
}

// Test if this error is one of parents of error `err` passed in parameter
func (goErr *GoError) IsParentOf(err error) bool {
	gerr, ok := err.(IError)
	if !ok {
		return false
	}

	name := goErr.GetName()

	for _, parent := range gerr.getParents() {
		if parent == name {
			return true
		}
	}

	return false
}

func (goErr *GoError) Init(value interface{}, message string, data interface{}, source error, pruneLevels uint) IError {
	if goErr.errType == nil {
		goErr.setType(value)

		goErr.message = message
		goErr.data = data
		goErr.source = source

		goErr.populateStackTrace(pruneLevels + 1)
	}

	return goErr
}

func (goErr *GoError) raise(pruneLevels uint) {
	if pruneLevels < 0 {
		pruneLevels = 0
	}

	res := goErr.getReference()
	res.populateStackTrace(pruneLevels + 1)

	panic(res)
}

func (goErr *GoError) setType(value interface{}) {
	errType := reflect.ValueOf(value).Type()
	if errType.Kind() == reflect.Ptr {
		errType = errType.Elem()
	}

	goErr.errType = errType
}

// Get the real reference on this error
func (goErr *GoError) getReference() IError {
	if goErr.errType == nil {
		goErr.setType(goErr)
	}

	ptr := unsafe.Pointer(reflect.ValueOf(goErr).Pointer())

	return reflect.NewAt(goErr.errType, ptr).Interface().(IError)
}

// This method construct the stack trace only in 'Debug' Mode
func (goErr *GoError) populateStackTrace(pruneLevels uint) {
	// If we aren't in debugging mode,
	if !errDebug {
		// Do nothing
		return
	}

	goErr.trace = getTrace(pruneLevels + 1)
}

// Get type of this error
func (goErr *GoError) getParents() []string {
	name := goErr.GetName()
	res, ok := errorHierarchies[name]

	if !ok {
		res = _getTypeHierarchy(goErr.errType, reflect.TypeOf(goErr).Elem())
		errorHierarchies[name] = res
	}

	return res
}

// GetSource
func GetSource(err error) error {
	ierr, ok := err.(IError)
	if !ok {
		return nil
	}

	return ierr.GetSource()
}
