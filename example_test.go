package goerrors_test

import (
    "fmt"
    "os"

    "github.com/corebreaker/goerrors"
)

func Example() {
    // A function which open a file
    open_file := func(name string) (*os.File, error) {
        f, err := os.Open(name)

        // Decorate the opening error
        if err != nil {
            return nil, goerrors.DecorateError(err)
        }

        return f, nil
    }

    // A function which read one byte in the opened file
    read_file := func(f *os.File) (byte, error) {
        var b [1]byte

        n, err := f.Read(b[:])

        // Decorate the read error
        if err != nil {
            return 0, goerrors.DecorateError(err)
        }

        // Return custom error
        if n == 0 {
            return 0, goerrors.MakeError("No data to read")
        }

        return b[0], nil
    }

    // Activate stack trace
    goerrors.SetDebug(true)

    // Call the checked open function
    f, err := open_file("a_file.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)

        return
    }

    // Here, in this example, this code won't never be executed if the file can't be opened
    defer f.Close()

    _, err = read_file(f)

    // Output:
}
