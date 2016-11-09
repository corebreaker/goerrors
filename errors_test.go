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

    raw_err, ok := err.(*GoError)
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

    err := MakeError("%s:%d", id, r).(*GoError)
    if err.source.Error() != fmt.Sprintf("%s:%d", id, r) {
        t.Error("Failed on MakeError:", err.source.Error(), "!=", fmt.Sprintf("%s:%d", id, r))
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

    err := MakeError("").(*GoError)
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
