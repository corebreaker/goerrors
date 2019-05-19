package goerrors

import (
	"fmt"
	"runtime"
	"strings"
)

// STACKTRACE_MAXLEN this version of stack trace asks to have a limit which arbitrary set.
const STACKTRACE_MAXLEN = 65536

// Construct formated stack trace.
func getTrace(start uint) []string {
	// The resulting stack trace.
	var trace []string

	// The caller list.
	callers := make([]uintptr, STACKTRACE_MAXLEN)

	// Gets the caller list, and returns an empty stack trace if there is no caller.
	n := runtime.Callers(int(start+2), callers)
	if n == 0 {
		return trace
	}

	// Get frames from callers.
	frames := runtime.CallersFrames(callers[:n])

	// Populates the stack trace.
	for hasMore := true; hasMore; {
		// A stack frame.
		var frame runtime.Frame

		// Gets the next frame/
		frame, hasMore = frames.Next()

		// If the file is from `runtime` package, so go to the next frame.
		if strings.Contains(frame.File, "runtime/") {
			continue
		}

		// Adds the stack trace entry.
		trace = append(trace, fmt.Sprintf("%s (%s:%d)", frame.Function, frame.File, frame.Line))
	}

	// Returns stack trace.
	return trace
}
