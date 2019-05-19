package goerrors

import (
	"fmt"
	"os"
)

func Example() {
	// A function which open a file
	openFile := func(name string) (*os.File, error) {
		f, err := os.Open(name)

		// Decorate the opening error
		if err != nil {
			return nil, DecorateError(err)
		}

		return f, nil
	}

	// A function which read one byte in the opened file
	readFile := func(f *os.File) (byte, error) {
		var b [1]byte

		n, err := f.Read(b[:])

		// Decorate the read error
		if err != nil {
			return 0, DecorateError(err)
		}

		// Return custom error
		if n == 0 {
			return 0, MakeError("No data to read")
		}

		return b[0], nil
	}

	// Deactivate stack trace
	// (cause stacktrace produced for testing package is specific to go installation and may change with Go version)
	SetDebug(false)

	// Make an unfindable filename
	const name = ".a_file_5123351069599224559.txt"

	// Call the checked open function
	f, err := openFile(name)
	if err != nil {
		fmt.Println(err)

		return
	}

	// Here, in this example, this code won't never be executed if the file can't be opened
	defer f.Close()

	_, err = readFile(f)
	if err != nil {
		fmt.Println(err)

		return
	}

	// Output:
	// StandardError: open .a_file_5123351069599224559.txt: no such file or directory
}
