package goerrors

import (
    "bytes"
    "fmt"
    "reflect"
    "unsafe"
)

var (
    // Inheritance cache
    error_hierarchies = make(map[string][]string)
)

// Handler for executing Try, Catch, or Finally block
type ErrorHandler func(err IError) error

// Basic error structure
type GoError struct {
    source   error        // Cause or original error
    message  string       // Error message
    trace    []string     // Stack trace
    data     interface{}  // Custom data
    err_type reflect.Type // Type of this error
}

// Standard method of `error` interface
func (self *GoError) Error() string {
    var out bytes.Buffer

    err := self.get_reference()

    // Prints error name
    fmt.Fprintf(&out, "%s: ", err.GetName())

    // Get informations
    message := err.GetMessage()
    source := err.GetSource()
    data := err.GetData()

    // Prints error informations
    if message != "" {
        fmt.Fprintln(&out, message)

        if data != nil {
            fmt.Fprintln(&out, data)
        }

        if source != nil {
            fmt.Fprintln(&out)
            fmt.Fprintln(&out, "Source:", source)
        }
    } else {
        if source != nil {
            fmt.Fprintln(&out, source)
        }

        if data != nil {
            fmt.Fprintln(&out, data)
        }
    }

    // Prints stack trace
    for _, entry := range self.trace {
        fmt.Fprintln(&out, "   ", entry)
    }

    // Prints a separator if stack trace is not empty
    if len(self.trace) > 0 {
        fmt.Fprintln(&out, "------------------------------------------------------------------------------")
    }

    // Return content of the buffer resulting from printing theses informations
    return out.String()
}

// Get error name
func (self *GoError) GetName() string {
    return self.err_type.PkgPath() + "." + self.err_type.Name()
}

// Get cause error (parent error)
func (self *GoError) GetSource() error {
    return self.source
}

// Get error message
func (self *GoError) GetMessage() string {
    return self.message
}

// Get custon data
func (self *GoError) GetData() interface{} {
    return self.data
}

// Complete try/catch/finally block
func (self *GoError) Try(try, catch, finally ErrorHandler) (err error) {
    defer self.Catch(&err, catch, finally)

    return try(self.get_reference())
}

// Catch error (used as a defered call)
func (self *GoError) Catch(err *error, catch, finally ErrorHandler) {
    var res_err error = nil

    defer func() {
        if finally != nil {
            ierr, _ := res_err.(IError)

            res_err = finally(ierr)
        }

        if err != nil {
            *err = res_err
        }
    }()

    recovered := recover()
    if recovered == nil {
        return
    }

    var ok bool

    if res_err, ok = recovered.(error); !ok {
        panic(recovered)
    }

    if this := self.get_reference(); !this.IsParentOf(res_err) {
        if err == nil {
            panic(recovered)
        }

        return
    }

    if catch != nil {
        res_err = catch(res_err.(IError))
    }
}

// Raise error
func (self *GoError) Raise() {
    self.raise(-1)
}

// Test if this error is one of parents of error `err` passed in parameter
func (self *GoError) IsParentOf(err error) bool {
    gerr, ok := err.(IError)
    if !ok {
        return false
    }

    name := self.GetName()

    for _, parent := range gerr.get_parents() {
        if parent == name {
            return true
        }
    }

    return false
}

func (self *GoError) Init(value interface{}, message string, data interface{}, source error, prune_levels int) IError {
    self.set_type(value)

    self.message = message
    self.data = data
    self.source = source

    self.populate_stack_trace(prune_levels)

    return self
}

func (self *GoError) raise(prune_levels int) {
    if prune_levels < 0 {
        prune_levels = 0
    }

    res := self.get_reference()
    res.populate_stack_trace(prune_levels + 1)

    panic(res)
}

func (self *GoError) set_type(value interface{}) {
    err_type := reflect.ValueOf(value).Type()
    if err_type.Kind() == reflect.Ptr {
        err_type = err_type.Elem()
    }

    self.err_type = err_type
}

// Get the real reference on this error
func (self *GoError) get_reference() IError {
    if self.err_type == nil {
        self.set_type(self)
    }

    ptr := unsafe.Pointer(reflect.ValueOf(self).Pointer())

    return reflect.NewAt(self.err_type, ptr).Interface().(IError)
}

// This method construct the stack trace only in 'Debug' Mode
func (self *GoError) populate_stack_trace(prune_levels int) {
    // If we aren't in debugging mode,
    if err_debug {
        // Do nothing
        return
    }

    self.trace = _make_stacktrace(prune_levels)
}

// Get type of this error
func (self *GoError) get_parents() []string {
    name := self.GetName()
    res, ok := error_hierarchies[name]

    if !ok {
        res = _get_parents(self.err_type, reflect.TypeOf(self).Elem())
        error_hierarchies[name] = res
    }

    return res
}

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
    get_reference() IError

    // This method construct the stack trace only in 'Debug' Mode
    populate_stack_trace(prune_levels int)

    // Get type of this error
    get_parents() []string

    // Raise error with pruned levels
    raise(prune_levels int)
}
