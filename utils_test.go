package goerrors

import "testing"

func TestListLen(t *testing.T) {
    if _list_len(nil) != 0 {
        t.Error("Nil list length is not zero")
    }

    if _list_len(make([]string, 0, 0)) != 0 {
        t.Error("Empty list length is not zero")
    }

    l := _list_len(make([]string, 10, 10))

    if l == 0 {
        t.Error("Any list length is zero")
    }

    if l < 0 {
        t.Error("Any list length is not positive")
    }

    if l != 10 {
        t.Error("Any list length has not the good value")
    }
}

func same_ptr(v1, v2 interface{}) bool {
    return v1 == v2
}

func TestConcat(t *testing.T) {
    if _concat(nil, nil) != nil {
        t.Error("Concat 2 nil lists should be nil")
    }

    lst := make([]string, 0)

    if !same_ptr(_concat(lst, nil), lst) {
        t.Error("Concat with RHS nil list should return LHS list")
    }

    if !same_ptr(_concat(nil, lst), lst) {
        t.Error("Concat with LHS nil list should return RHS list")
    }

    if len(_concat(lst, nil)) != 0 {
        t.Error("Concat with LHS empty list should have length as zero")
    }

    if len(_concat(nil, lst)) != 0 {
        t.Error("Concat with RHS empty list should have length as zero")
    }

    if len(_concat(lst, lst)) != 0 {
        t.Error("Concat same empty list should have length as zero")
    }

    if len(_concat(lst, make([]string, 0))) != 0 {
        t.Error("Concat two empty lists should have length as zero")
    }

    lst = []string{"A", "B"}

    if len(_concat(lst, lst)) != (2 * len(lst)) {
        t.Error("Concat same list should have length as twice the length of the list")
    }

    if len(_concat(lst, []string{"X"})) != 3 {
        t.Error("Concat two list should have the good length")
    }
}
