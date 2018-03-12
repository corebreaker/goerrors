package goerrors

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/go-xweb/uuid"
)

type MyError struct {
	GoError
}

var (
	__err_test MyError
	__err_type = reflect.TypeOf(__err_test)
)

func same_ptr(v1, v2 interface{}) bool {
	return reflect.ValueOf(v1).Pointer() == reflect.ValueOf(v2).Pointer()
}

func TestSetType(t *testing.T) {
	gerr := new(MyError)
	gerr.set_type(gerr)

	if gerr.err_type != __err_type {
		t.Fail()
	}
}

func TestInitError(t *testing.T) {
	src := fmt.Errorf("")
	data := new(int)

	gerr := new(MyError)
	gerr.Init(gerr, "--message--", data, src, -1)

	if gerr.err_type != __err_type {
		t.Error("Bad type")
	}

	if gerr.message != "--message--" {
		t.Error("Bad message")
	}

	if !same_ptr(gerr.data, data) {
		t.Error("Bad data")
	}

	if !same_ptr(gerr.source, src) {
		t.Error("Bad data")
	}
}

func TestShowErrorNoDebug(t *testing.T) {
	SetDebug(false)

	moment := time.Now()
	rand.Seed(moment.UnixNano())

	r := rand.Int63()
	id := uuid.NewUUID()

	err := MakeError("%s:%d", id, r).(*tStandardError)
	if err.message != fmt.Sprintf("%s:%d", id, r) {
		t.Error("Failed on MakeError:", err.message, "!=", fmt.Sprintf("%s:%d", id, r))
	}

	if len(err.trace) != 0 {
		t.Error("Failed on empty trace: ", err.trace)
	}

	if err.infos.Len() != 0 {
		t.Error("Failed on no-infos")
	}
}

func TestShowErrorDebug(t *testing.T) {
	SetDebug(true)

	err := MakeError("").(*tStandardError)
	if len(err.trace) == 0 {
		t.Error("Failed on stack trace")
	}
}

func TestErrorMethods(t *testing.T) {
	SetDebug(true)

	err := MakeError("")
	if len(err.Error()) == 0 {
		t.Error("Failed on infos")
	}
}
