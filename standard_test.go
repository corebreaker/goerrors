package goerrors

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestDecorateError(t *testing.T) {
	if DecorateError(nil) != nil {
		t.Fail()
	}
}

func TestAddInfo(t *testing.T) {
	err := fmt.Errorf("")

	moment := time.Now()
	rand.Seed(moment.UnixNano())

	r := rand.Int63()
	id := uuid.New()

	err = AddInfo(err, "%s:%d", id, r)

	rawErr, ok := err.(*tStandardError)
	if !ok {
		t.Error("Failed on convesion")
	}

	if rawErr.infos.Len() == 0 {
		t.Error("Failed on has-infos")
	}

	if rawErr.infos.String() != fmt.Sprintf("%s:%d\n", id, r) {
		t.Error("Failed on set-infos:", rawErr.infos.String(), "!=", fmt.Sprintf("%s:%d\n", id, r))
	}
}

func TestDecorateNilErrorWithDatas(t *testing.T) {
	if DecorateErrorWithDatas(nil, 0, nil, "") != nil {
		t.Fail()
	}
}

func TestDecorateErrorWithDatas(t *testing.T) {
	err := DecorateErrorWithDatas(errors.New("error"), 1234, "data", "hello %d", 1)
	if err == nil {
		t.Fatal("Error should not be nil")
	}

	err = DecorateErrorWithDatas(err, 1234, "data", "hello %d", 1)
	if err == nil {
		t.Fatal("Error should not be nil")
	}
}

func TestErrorDecoration(t *testing.T) {
	moment := time.Now()
	rand.Seed(moment.UnixNano())

	code := rand.Int63()

	r := rand.Int63()
	id := uuid.New()

	err := MakeErrorWithDatas(code, nil, "%s:%d", id, r)
	if err.GetCode() != code {
		t.Fail()
	}
}

func TestMakeErrorNoDebug(t *testing.T) {
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
		t.Error("Failed on empty trace")
	}

	if err.infos.Len() != 0 {
		t.Error("Failed on no-infos")
	}
}

func TestMakeErrorDebug(t *testing.T) {
	SetDebug(true)

	err := MakeError("").(*tStandardError)
	if len(err.trace) == 0 {
		t.Error("Failed on stack trace")
	}
}

func TestError(t *testing.T) {
	SetDebug(true)

	err := MakeError("")
	if len(err.Error()) == 0 {
		t.Error("Failed on infos")
	}
}

func TestCatch(t *testing.T) {
	var err error

	func() {
		Catch(&err, func(err IError) error {
			return err
		}, func(err IError) error {
			return err
		})
	}()

	func() {
		defer Catch(&err, func(err IError) error {
			return err
		}, func(err IError) error {
			return err
		})

		panic(fmt.Errorf("error"))
	}()

	func() {
		defer DiscardPanic()
		defer Catch(&err, nil, nil)

		panic("Error")
	}()

	func() {
		defer Catch(&err, nil, nil)

		panic(fmt.Errorf("error"))
	}()

	Catch(&err, func(err IError) error {
		return nil
	}, func(err IError) error {
		return nil
	})
}

func TestStdTry(t *testing.T) {
	err := Try(func(err IError) error {
		return nil
	}, nil, nil)

	if err != nil {
		t.Fail()
	}

	err = Try(func(err IError) error {
		return fmt.Errorf("error")
	}, func(err IError) error {
		return err
	}, func(err IError) error {
		return err
	})

	err = Try(func(err IError) error {
		panic(fmt.Errorf("error"))
	}, func(err IError) error {
		return err
	}, nil)

	err = Try(func(err IError) error {
		panic("Error")
	}, func(err IError) error {
		return err
	}, nil)
}

func TestStdRaise(t *testing.T) {
	func() {
		defer DiscardPanic()

		Raise("")
	}()
}

func TestRaiseWithInfos(t *testing.T) {
	func() {
		defer DiscardPanic()

		RaiseWithInfos(0, nil, "")
	}()
}

func TestRaiseError(t *testing.T) {
	func() {
		defer DiscardPanic()

		RaiseError(fmt.Errorf("error"))
	}()
}
