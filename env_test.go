package goerrors

import (
	"testing"
)

func TestDebug(t *testing.T) {
	SetDebug(true)
	if !GetDebug() {
		t.Errorf("Error on true (first)")
	}

	SetDebug(false)
	if GetDebug() {
		t.Errorf("Error on false (first)")
	}

	SetDebug(true)
	if !GetDebug() {
		t.Errorf("Error on true (second)")
	}

	SetDebug(false)
	if GetDebug() {
		t.Errorf("Error on false (second)")
	}
}
