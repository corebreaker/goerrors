package goerrors

import (
	"testing"
)

func TestNoTrace(t *testing.T) {
	trace := getTrace(0)
	if len(getTrace(uint(len(trace))+1)) > 0 {
		t.Fail()
	}
}
