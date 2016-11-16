package goerrors

import (
    "reflect"
    "testing"
)

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
    return reflect.ValueOf(v1).Pointer() == reflect.ValueOf(v2).Pointer()
}

func TestConcat(t *testing.T) {
    if _concat(nil, nil) != nil {
        t.Error("Concat 2 nil lists should be nil")
    }

    lst0 := make([]string, 0)
    lst1 := []string{"A", "B"}

    if !same_ptr(_concat(lst0, nil), lst0) {
        t.Error("Concat with RHS nil list should return LHS list")
    }

    if !same_ptr(_concat(nil, lst0), lst0) {
        t.Error("Concat with LHS nil list should return RHS list")
    }

    if len(_concat(lst0, nil)) != 0 {
        t.Error("Concat with LHS empty list should have length as zero")
    }

    if len(_concat(nil, lst0)) != 0 {
        t.Error("Concat with RHS empty list should have length as zero")
    }

    if len(_concat(lst0, lst0)) != 0 {
        t.Error("Concat same empty list should have length as zero")
    }

    if len(_concat(lst0, make([]string, 0))) != 0 {
        t.Error("Concat two empty lists should have length as zero")
    }

    if !same_ptr(_concat(lst1, nil), lst1) {
        t.Error("Concat with LHS non-empty list and RHS nil list should return LHS list")
    }

    if !same_ptr(_concat(nil, lst1), lst1) {
        t.Error("Concat with RHS non-empty list and LHS nil list should return RHS list")
    }

    if !same_ptr(_concat(lst1, lst0), lst1) {
        t.Error("Concat with LHS non-empty list and RHS empty list should return LHS list")
    }

    if !same_ptr(_concat(lst0, lst1), lst1) {
        t.Error("Concat with RHS non-empty list and LHS empty list should return RHS list")
    }

    if len(_concat(lst1, lst1)) != (2 * len(lst0)) {
        t.Error("Concat same list should have length as twice the length of the list")
    }

    if len(_concat(lst1, []string{"X"})) != 3 {
        t.Error("Concat two list should have the good length")
    }
}

func TestStructHierarchy(t *testing.T) {
    type X struct {
    }

    type A struct {
    }

    type B struct {
        A
        X
    }

    type C struct {
        B
    }

    type D struct {
        A
        B
        X
    }

    type E struct {
        D
    }

    tA := reflect.TypeOf(A{})
    tB := reflect.TypeOf(B{})
    tC := reflect.TypeOf(C{})
    tD := reflect.TypeOf(C{})
    tE := reflect.TypeOf(C{})

    n1 := "github.com/corebreaker/goerrors.A"
    n2 := "github.com/corebreaker/goerrors.B"
    n3 := "github.com/corebreaker/goerrors.C"
    n4 := "github.com/corebreaker/goerrors.D"

    res := _get_type_hierarchy(tA, tA)

    if res == nil {
        t.Error("With zero level, it returns nil")
    }

    if len(res) != 0 {
        t.Error("With zero level, it should return an empty list", "; the return is:", res)
    }

    res = _get_type_hierarchy(tB, tA)

    if res == nil {
        t.Error("With one level, it returns nil")
    }

    if len(res) != 1 {
        t.Fatal("With one level, it should return a list with length 1", "; the return is:", res)
    }

    if res[0] != n1 {
        t.Error("With one level, it should return a list one type", n1, "; the return is:", res)
    }

    res = _get_type_hierarchy(tC, tA)

    if res == nil {
        t.Error("With two levels, it returns nil")
    }

    if len(res) != 2 {
        t.Fatal("With two levels, it should return a list with length 2", "; the return is:", res)
    }

    if (res[0] != n1) && (res[1] != n2) {
        t.Error("With two levels, it should return a list with types", n1, "and", n2, "; the return is:", res)
    }

    res = _get_type_hierarchy(tD, tA)

    if res == nil {
        t.Error("With two levels, it returns nil")
    }

    if len(res) != 2 {
        t.Fatal("With two levels, it should return a list with length 2", "; the return is:", res)
    }

    if (res[0] != n1) && (res[1] != n2) {
        t.Error("With two levels, it should return a list with types", n1, "and", n2, "; the return is:", res)
    }

    res = _get_type_hierarchy(tE, tA)

    if res == nil {
        t.Error("With two levels, it returns nil")
    }

    if len(res) != 2 {
        t.Fatal("With two levels, it should return a list with length 2", "; the return is:", res)
    }

    if (res[0] != n1) && (res[1] != n2) {
        t.Error("With two levels, it should return a list with types", n1, "and", n2, "; the return is:", res)
    }
}
