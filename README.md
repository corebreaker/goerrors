# goerrors
Easy error stack trace for Go.

It's a package that's allow you to make Go errors more comprehensive and more usable with a cost of little effort to specify your errors.

## An simple example
```go
package main

import (
    "fmt"
    err "github.com/corebreaker/goerrors"
)

// A function which return checked quotient
func my_func(i, j int) (int, error) {
    if j == 0 {
        return 0, err.MakeError("Division by zero")
    }
    
    return i / j, nil
}

// Main function
func main() {
    // Activate stack trace
    err.SetDebug(true)
    
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
    github.com/corebreaker/goerrors.MakeError (/data/development/projects/go/src/github.com/corebreaker/goerrors/errors.go:73)
    main.my_func (/home/frederic/.local/share/data/liteide/liteide/goplay.go:11)
    main.main (/home/frederic/.local/share/data/liteide/liteide/goplay.go:23)
    runtime.main (/opt/go/src/runtime/proc.go:183)
    runtime.goexit (/opt/go/src/runtime/asm_amd64.s:2086)
------------------------------------------------------------------------------
```
