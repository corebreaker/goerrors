package goerrors

import (
    "bytes"
    "fmt"
    "runtime"
)

// Interface for extended Go errors
type GoError interface {
    // This is an error
    error

    // Add more informations on that error
    AddInfo(info string, args ...interface{}) GoError
}

// Internal error structure
type tError struct {
    source error        // Cause or original error
    infos  bytes.Buffer // Additionnal informations
    trace  []string     // Stack trace
}

// Standard method of `error` interface
// @see error.Error
func (self *tError) Error() string {
    var out bytes.Buffer

    // Prints error informations
    fmt.Fprintln(&out, self.source)
    fmt.Fprint(&out, &self.infos)

    // Prints stack trace
    for _, line := range self.trace {
        fmt.Fprintln(&out, "   ", line)
    }

    // Prints a separator if stack trace is not empty
    if len(self.trace) > 0 {
        fmt.Fprintln(&out, "------------------------------------------------------------------------------")
    }

    // Return content of the buffer resulting from printing theses informations
    return out.String()
}

// Add informations in extended Go error
func (self *tError) AddInfo(info string, args ...interface{}) GoError {
    // Just prints into internal buffer the informations passed as parameters
    fmt.Fprintln(&self.infos, fmt.Sprintf(info, args...))

    return self
}

// Error decorator
// Constructs an extended Go error from an existing error
func decorate_error(err error, prune_levels int) GoError {
    // If error is nil, therefore returns nil
    if err == nil {
        return nil
    }

    // Checks that pruned levels passed in parameted is really a positive value
    // A negative value means, no pruning
    if prune_levels < 0 {
        prune_levels = 0
    }

    // Program Counter initialisation
    var pc uintptr = 1

    // Resulting stack trace
    stack := make([]string, 0)

    // If we are in debugging mode,
    if err_debug {
        // Populate the stack trace
        for i := prune_levels + 2; pc != 0; i++ {
            // Retreive runtime informations
            ptr, file, line, ok := runtime.Caller(i)
            pc = ptr

            // If there isn't significant information, go to next level
            if (pc == 0) || (!ok) {
                continue
            }

            // Retreive called function
            f := runtime.FuncForPC(pc)

            // Add stack trace entry
            stack = append(stack, fmt.Sprintf("%s (%s:%d)", f.Name(), file, line))
        }
    }

    // Finalize by constructing the resulting extended Go error
    return &tError{
        source: err,
        trace:  stack,
    }
}

// «Decorate» the error passed as "err" parameter.
// The error returned will be an extended Go error with additionnal informations and stack trace.
func DecorateError(err error) GoError {
    return decorate_error(err, -1)
}

// Make an extended Go error from a message passed as "message" parameter
func MakeError(message string, args ...interface{}) GoError {
    // make a standard error with fmt.Errorf, then decorate it with pruning one level in stack trace
    // (to eliminate this function calling)
    return decorate_error(fmt.Errorf(message, args...), -1)
}

// Global function to add information in an error whatever.
// This function just call the "AddInfo" method of an extended Go error.
func AddInfo(err error, info string, args ...interface{}) GoError {
    // Check if "err" is already an extended Go error
    go_err, ok := err.(GoError)
    if !ok {
        // Otherwise decorate that unknown error
        go_err = DecorateError(err)
    }

    // Delegate call to "AddInfo" method
    return go_err.AddInfo(info, args...)
}
