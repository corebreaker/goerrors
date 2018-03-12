package goerrors

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-xweb/uuid"
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
	id := uuid.NewUUID()

	err = AddInfo(err, "%s:%d", id, r)

	raw_err, ok := err.(*tStandardError)
	if !ok {
		t.Error("Failed on convesion")
	}

	if raw_err.infos.Len() == 0 {
		t.Error("Failed on has-infos")
	}

	if raw_err.infos.String() != fmt.Sprintf("%s:%d\n", id, r) {
		t.Error("Failed on set-infos:", raw_err.infos.String(), "!=", fmt.Sprintf("%s:%d\n", id, r))
	}
}

func TestMakeErrorNoDebug(t *testing.T) {
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

		panic(fmt.Errorf("Error"))
	}()

	func() {
		defer DiscardPanic()
		defer Catch(&err, nil, nil)

		panic("Error")
	}()

	func() {
		defer Catch(&err, nil, nil)

		panic(fmt.Errorf("Error"))
	}()

	Catch(&err, func(err IError) error {
		return nil
	}, func(err IError) error {
		return nil
	})
}

func TestTry(t *testing.T) {
	Try(func(err IError) error {
		return nil
	}, nil, nil)

	Try(func(err IError) error {
		return fmt.Errorf("Error")
	}, func(err IError) error {
		return err
	}, func(err IError) error {
		return err
	})

	Try(func(err IError) error {
		panic(fmt.Errorf("Error"))
	}, func(err IError) error {
		return err
	}, nil)

	Try(func(err IError) error {
		panic("Error")
	}, func(err IError) error {
		return err
	}, nil)
}

func TestRaise(t *testing.T) {
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

		RaiseError(fmt.Errorf("Error"))
	}()
}
