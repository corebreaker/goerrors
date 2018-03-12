// +build !go1.9

package goerrors

import (
	"fmt"
	"runtime"
	"strings"
)

// Construct formated stack trace.
func getTrace(start uint) []string {
	// Program Counter initialisation.
	var pc uintptr = 1

	// The resulting stack trace.
	var trace []string

	// Populates the stack trace.
	for i := 1 + start; pc != 0; i++ {
		// Retreives runtime informations.
		ptr, file, line, ok := runtime.Caller(int(i))
		pc = ptr

		// If there isn't significant information, go to next level.
		if (pc == 0) || (!ok) {
			continue
		}

		// Retreives the called function.
		f := runtime.FuncForPC(pc)

		// If the file is from `runtime` package, so go to the next frame.
		if strings.Contains(file, "runtime/") {
			continue
		}

		// Adds the stack trace entry.
		trace = append(trace, fmt.Sprintf("%s (%s:%d)", f.Name(), file, line))
	}

	// Returns stack trace.
	return trace
}
