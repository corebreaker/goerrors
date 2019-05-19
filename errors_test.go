package goerrors

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

type MyError struct {
	GoError
}

var (
	__errTest MyError
	__errType = reflect.TypeOf(__errTest)
)

func samePtr(v1, v2 interface{}) bool {
	return reflect.ValueOf(v1).Pointer() == reflect.ValueOf(v2).Pointer()
}

func TestSetType(t *testing.T) {
	gerr := &MyError{}
	gerr.setType(gerr)

	if gerr.errType != __errType {
		t.Fail()
	}
}

func TestInitError(t *testing.T) {
	src := fmt.Errorf("")
	data := new(int)

	gerr := MyError{}
	_ = gerr.Init(gerr, "--message--", data, src, 0)

	if gerr.errType != __errType {
		t.Error("Bad type")
	}

	if gerr.message != "--message--" {
		t.Error("Bad message")
	}

	if !samePtr(gerr.data, data) {
		t.Error("Bad data")
	}

	if !samePtr(gerr.source, src) {
		t.Error("Bad data")
	}
}

func TestShowErrorNoDebug(t *testing.T) {
	SetDebug(false)

	moment := time.Now()
	rand.Seed(moment.UnixNano())

	r := rand.Int63()
	id := uuid.New()

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

func TestErrorString(t *testing.T) {
	gerr := &MyError{}
	_ = gerr.Init(gerr, "--message--", 1234, errors.New("error"), 0)

	_ = gerr.Error()

	gerr.message = ""
	_ = gerr.Error()
}

func TestGetSource(t *testing.T) {
	_ = GetSource(nil)
	_ = GetSource(DecorateError(fmt.Errorf("error")))
}

func TestErrorTry(t *testing.T) {
	gerr := &MyError{}
	_ = gerr.Init(gerr, "--message--", nil, nil, 0)

	_ = gerr.Try(func(err IError) error {
		return err
	}, nil, func(err IError) error { return nil })

	_ = gerr.Try(func(err IError) error {
		err.Raise()

		return nil
	}, func(err IError) error {
		return err
	}, nil)
}

func TestErrorCatchWithPrimitiveError(t *testing.T) {
	defer DiscardPanic()

	gerr := &MyError{}
	_ = gerr.Init(gerr, "--message--", nil, nil, 0)

	_ = gerr.Try(func(err IError) error {
		panic("error")
	}, nil, nil)
}

func TestErrorCatchWithBasicError(t *testing.T) {
	defer DiscardPanic()

	gerr := &MyError{}
	_ = gerr.Init(gerr, "--message--", nil, nil, 0)

	_ = gerr.Try(func(err IError) error {
		panic(errors.New("error"))
	}, nil, nil)
}

func TestErrorCatchWithOtherError(t *testing.T) {
	defer DiscardPanic()

	gerr := &MyError{}
	_ = gerr.Init(gerr, "--message--", nil, nil, 0)

	_ = gerr.Try(func(err IError) error {
		otherError := &struct{ GoError }{}

		otherError.Init(otherError, "error", nil, nil, 0).Raise()

		return nil
	}, nil, nil)
}

func TestErrNakedRaise(t *testing.T) {
	defer DiscardPanic()

	(&MyError{}).raise(0)
}
