# GoErrors - Easy extended Go error with rich informations.
It's a package that's allow you to make Go errors more comprehensive and easier to use.

- Verbose with stack trace
- Ready for Try/Catch/Finally paradigm
- Extensible with your own error
- Entirely customizable

[![GoDoc](https://godoc.org/github.com/corebreaker/goerrors?status.svg)](https://godoc.org/github.com/corebreaker/goerrors)


## Installation
go get github.com/corebreaker/goerrors


## How is it work ?
An normal error is &laquo;decorated&raquo; with informations to obtain a new type of error. These informations are:
- stack trace informations (files + lines)
- additionnal customized string

When the `Error` method is called, all these information will be formated in a string.

This type of error gives you rich informations on an error without using `panic` function.


## An simple example
```go
package main

import (
    "fmt"
    gerr "github.com/corebreaker/goerrors"
)

// A function which return checked quotient
func my_func(i, j int) (int, error) {
    if j == 0 {
        return 0, gerr.MakeError("Division by zero")
    }

    return i / j, nil
}

// Main function
func main() {
    // Activate stack trace
    gerr.SetDebug(true)

    i, j := 1, 0

    // Call checked function
    q, err := my_func(i, j)
    if err != nil {
        fmt.Println(err)

        return
    }

    // Here, in this example, this code won't never be executed
    fmt.Print(i, "/", j, "=", q)
}
```

This will show:
```
Division by zero
    main.my_func (/projects/go/prototype/main.go:11)
    main.main (/projects/go/prototype/main.go:23)
    runtime.main (/opt/go/src/runtime/proc.go:183)
    runtime.goexit (/opt/go/src/runtime/asm_amd64.s:2086)
------------------------------------------------------------------------------
```


## Another example with existing error
```go
package main

import (
    "fmt"
    "os"
    gerr "github.com/corebreaker/goerrors"
)

// A function which open a file
func open_file(name string) (*os.File, error) {
    f, err := os.Open(name)

    // Decorate the opening error
    if err != nil {
        return nil, gerr.DecorateError(err)
    }

    return f, nil
}

// A function which read one byte in the opened file
func read_file(f *os.File) (byte, error) {
    var b [1]byte

    n, err := f.Read(b[:])

    // Decorate the read error
    if err != nil {
        return 0, gerr.DecorateError(err)
    }

    // Return custom error
    if n == 0 {
        return 0, gerr.MakeError("No data to read")
    }

    return b[0], nil
}

// Main function
func main() {
    // Activate stack trace
    gerr.SetDebug(true)

    // Call the checked open function
    f, err := open_file("a_file.txt")
    if err != nil {
        fmt.Println(err)

        return
    }

    // Here, in this example, this code won't never be executed if the file can't be opened
    defer f.Close()

    _, err = read_file(f)
}
```

This will show:
```
open a_file.txt: no such file or directory
    main.open_file (/projects/go/src/github.com/corebreaker/goerrors.go:15)
    main.main (/projects/go/src/github.com/corebreaker/goerrors.go:46)
    runtime.main (/opt/go/src/runtime/proc.go:183)
    runtime.goexit (/opt/go/src/runtime/asm_amd64.s:2086)
------------------------------------------------------------------------------
```


## API Documentation
https://godoc.org/github.com/corebreaker/goerrors
