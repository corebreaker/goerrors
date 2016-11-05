package goerrors

import (
    "bytes"
    "fmt"
    "runtime"
)

type GoError interface {
    error

    AddInfo(info string, args ...interface{}) GoError
}

type tError struct {
    source error
    infos  bytes.Buffer
    trace  []string
}

func (self *tError) Error() string {
    var out bytes.Buffer

    fmt.Fprintln(&out, self.source)
    fmt.Fprint(&out, &self.infos)

    for _, line := range self.trace {
        fmt.Fprintln(&out, "   ", line)
    }

    if len(self.trace) > 0 {
        fmt.Fprintln(&out, "------------------------------------------------------------------------------")
    }

    return out.String()
}

func (self *tError) AddInfo(info string, args ...interface{}) GoError {
    fmt.Fprintln(&self.infos, fmt.Sprintf(info, args...))

    return self
}

func DecorateError(err error) GoError {
    if err == nil {
        return nil
    }

    var pc uintptr = 1

    stack := make([]string, 0)
    if err_debug {
        for i := 1; pc != 0; i++ {
            ptr, file, line, ok := runtime.Caller(i)
            pc = ptr
            if (pc == 0) || (!ok) {
                continue
            }

            f := runtime.FuncForPC(pc)

            stack = append(stack, fmt.Sprintf("%s (%s:%d)", f.Name(), file, line))
        }
    }

    return &tError{
        source: err,
        trace:  stack,
    }
}

func MakeError(message string, args ...interface{}) GoError {
    return DecorateError(fmt.Errorf(message, args...))
}
