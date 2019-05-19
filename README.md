# GoErrors - Easy error informations and stack traces for Go with rich informations and a Try/Catch/Finally mechanism.
[![Report](https://goreportcard.com/badge/github.com/corebreaker/goerrors?style=plastic)](https://goreportcard.com/report/github.com/corebreaker/goerrors)
[![Build Status](https://img.shields.io/travis/com/corebreaker/goerrors/master.svg?style=plastic)](https://travis-ci.com/corebreaker/goerrors)
[![Coverage Status](https://img.shields.io/coveralls/github/corebreaker/goerrors/master.svg?style=plastic)](https://coveralls.io/github/corebreaker/goerrors)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=plastic)](https://godoc.org/github.com/corebreaker/goerrors)
[![GitHub license](https://img.shields.io/github/license/corebreaker/goerrors.svg?color=blue&style=plastic)](https://github.com/corebreaker/goerrors/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/corebreaker/goerrors.svg?style=plastic)](https://github.com/corebreaker/goerrors/releases)
[![Dependency status](https://img.shields.io/librariesio/github/corebreaker/goerrors.svg?style=plastic)](https://libraries.io/github/corebreaker/goerrors)

It's a package that's allow you to make Go errors more comprehensive, more featured and easier to use.
This project is the realisation of the idea introduced by the GIST [try-catch-finally.go](https://gist.github.com/corebreaker/6a93c8210425e96dc1bcbb157f0270b0).

## Features
- Verbose with stack trace
- Ready for Try/Catch/Finally paradigm
- Extensible with your own error
- Error hierarchy
- Entirely customizable


## Installation
go get github.com/corebreaker/goerrors


## How is it work ?
A normal error is just an interface. But here we added an extended interface IError which can transport other infomations.
Therefore to add informations, the standard error (`error` interface) is &laquo;decorated&raquo; with informations and so
we obtain a new error value. These informations are:
- stack trace informations (files + lines)
- additionnal customized string
- optionnal error code

When the `Error` method is called, all these informations will be formated in a string.

This type of error gives you rich informations on an error without using `panic` function (even if with a panic, you
could show your additionnal infomations).


## How to use it ?
### Decorate error to add a stack trace
```go
// for example, let's open a file
func OpenMyFile() (*os.File, error) {
    file, err := os.Open("MyFile.txt")
    if err != nil {
        // Here, we decorate the error
        return nil, goerrors.DecorateError(err)
    }

    return file, nil
}
```

Then, we can `panic` this decorated error or simply print it, like this:
```go
func main() {
    // First we must enable the debug mode to activate the stacktrace
    goerrors.SetDebug(true)

    // Use a function that use the `goerrors` package
    file, err := OpenMyFile()

    // Check the error
    if err != nil {
        // Show the error
        fmt.Println(err)

        // Terminate
        return
    }

    // …
}
```

You will see some thing like this:
```
github.com/corebreaker/goerrors.tStandardError: open MyFile.txt: no such file or directory
    github.com/corebreaker/goerrors.(*GoError).Init (/home/frederic/go/src/github.com/corebreaker/goerrors/errors.go:219)
    github.com/corebreaker/goerrors.DecorateError (/home/frederic/go/src/github.com/corebreaker/goerrors/standard.go:51)
    main.OpenMyFile (/home/frederic/.local/share/data/liteide/liteide/goplay.go:13)
    main.main (/home/frederic/.local/share/data/liteide/liteide/goplay.go:24)
------------------------------------------------------------------------------
```

### A Try/Catch/Finally mechanism
Plus, this library uses the `panic()` function, the `recover()` function and the `defer` instruction,
as a Throw and a Try/Catch/Finally mechanisms and can be used like this:
```go
goerrors.Try(func(err goerrors.IError) error {
    // Try block
}, func(err goerrors.IError) error {
    // Catch block
}, func(err goerrors.IError) error {
    // Finally block
})
```

The error passed in `Try` block is the error which called the `Try` method:
```go
theError := goerrors.MakeError("the error in `Try` block")
theError.Try(func(err goerrors.IError) error {
	// Here err === theError
	
	return nil
})
```

In the case in using the `Try` function in the GoError package, the error passed as argument is an error created by the
call the `Try` function. Then, that error can be customized with the GoError API.

#### An example with a throw, called "raise" here

Actually, returning an error in the `Try` block is a `Throw`, and a Go `panic` call is too like a throw but there is
a `panic`-like function for keeping Try/Catch formalism, the `Raise` function, used like that:
```go
goerrors.Try(func(err goerrors.IError) error {
	if aCondition {
        // `Raise` call
        goerrors.Raise("an error")
    }
    
    // Do something
    
    return nil
}, func(err goerrors.IError) error {
    // Catch block
}, func(err goerrors.IError) error {
    // Finally block
})
```

At last, all errors generated by GoError have a `Raise` method. So, you can throw an error like that:
```go
goerrors.Try(func(err goerrors.IError) error {
    if aCondition {
        // `Raise` method call
        goerrors.MakeError("an error").Raise()
    }
    
    // Do something
    
    return nil
}, func(err goerrors.IError) error {
    // Catch block
}, func(err goerrors.IError) error {
    // Finally block
})
```

## A simple example
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
StandardError: Division by zero
    main.my_func (/projects/go/prototype/main.go:11)
    main.main (/projects/go/prototype/main.go:23)
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
func openFile(name string) (*os.File, error) {
    f, err := os.Open(name)

    // Decorate the opening error
    if err != nil {
        return nil, gerr.DecorateError(err)
    }

    return f, nil
}

// A function which read one byte in the opened file
func readFile(f *os.File) (byte, error) {
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
    f, err := openFile("a_file.txt")
    if err != nil {
        fmt.Println(err)

        return
    }

    // Here, in this example, this code won't never be executed if the file can't be opened
    defer f.Close()

    _, err = readFile(f)
}
```

This will show:
```
StandardError: open a_file.txt: no such file or directory
    main.open_file (/projects/go/src/github.com/corebreaker/goerrors.go:15)
    main.main (/projects/go/src/github.com/corebreaker/goerrors.go:46)
------------------------------------------------------------------------------
```


### A Try/Catch example with error inheritance
```go
package main

import (
    "fmt"
    gerr "github.com/corebreaker/goerrors"
)

type ErrorBase struct{ gerr.GoError }

type ErrorA struct{ ErrorBase }
type ErrorB struct{ ErrorBase }

type ChildError struct{ ErrorA }

func newErrorBase() gerr.IError {
    err := &ErrorBase{}

    return err.Init(err, "message from Error base", nil, nil, 0)
}

func newErrorB() gerr.IError {
    err := &ErrorB{}

    return err.Init(err, "message from Error B", nil, nil, 0)
}

func newChildError() gerr.IError {
    err := &ChildError{}

    return err.Init(err, "message from Child Error", nil, nil, 0)
}

// A function which raise and try to catch the error which is not in the same hierarchy
func myFunc() () {
    _ = newErrorB().Try(func(err gerr.IError) error {
        newChildError().Raise()

        return nil
    }, func(err gerr.IError) error {
        // This catch block will not called because ErrorB is not in the same error hierarchy of ChildError
        return nil
    }, nil)
}

// Main function
func main() {
    _ = newErrorBase().Try(func(err gerr.IError) error {
        myFunc()

        return nil
    }, func(err gerr.IError) error {
        fmt.Println("Catched error:", err)

        return nil
    }, nil)
}
```

This will show:
```
Catched error: main.ChildError: message from Child Error
```
