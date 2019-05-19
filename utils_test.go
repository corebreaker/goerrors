package goerrors

import (
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	if _concat(nil, nil) != nil {
		t.Error("Concatenation of 2 nil lists should be nil")
	}

	lst := _concat([]string{}, nil)
	if (lst == nil) || (len(lst) > 0) {
		t.Error("Concatenation of 2 empty list should be should: ", lst)
	}

	lst0 := make([]string, 0)
	lst1 := []string{"A", "B"}

	lst = _concat(lst0, nil)
	if (lst == nil) || (len(lst) > 0) {
		t.Error("Concat with RHS nil list should not return nil and non empty list (LHS list)")
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

	if len(_concat(lst1, lst1)) != (2 * len(lst1)) {
		t.Error("Concat same list should have length as twice the length of the list")
	}

	if len(_concat(lst1, []string{"X"})) != 3 {
		t.Error("Concat two list should have the good length")
	}
}

func TestStructHierarchy(t *testing.T) {
	type X struct {
		int
		j float32
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
		C
	}

	tA := reflect.TypeOf(A{})
	tB := reflect.TypeOf(B{})
	tC := reflect.TypeOf(C{})
	tD := reflect.TypeOf(D{})
	tE := reflect.TypeOf(E{})

	n1 := "github.com/corebreaker/goerrors.A"
	n2 := "github.com/corebreaker/goerrors.B"
	n3 := "github.com/corebreaker/goerrors.C"
	n4 := "github.com/corebreaker/goerrors.D"
	n5 := "github.com/corebreaker/goerrors.E"

	res := _getTypeHierarchy(tA, tA)

	if res == nil {
		t.Error("With zero level, it returns nil")
	}

	if len(res) != 1 {
		t.Fatal("With zero level, it should return a list with length 1; the return is:", res)
	}

	if res[0] != n1 {
		t.Errorf("With zero level, it should return a list one type %v; the return is: %v", n1, res)
	}

	res = _getTypeHierarchy(tB, tA)

	if res == nil {
		t.Error("With one level, it returns nil")
	}

	if len(res) != 2 {
		t.Fatal("With one level, it should return a list with length 2; the return is:", res)
	}

	if (res[0] != n2) || (res[1] != n1) {
		t.Errorf("With one levels, it should return a list with types %v and %v; the return is: %v", n2, n1, res)
	}

	res = _getTypeHierarchy(tC, tA)

	if res == nil {
		t.Error("With two levels, it returns nil")
	}

	if len(res) != 3 {
		t.Fatal("With two levels, it should return a list with length 3; the return is:", res)
	}

	if (res[0] != n3) || (res[1] != n2) || (res[2] != n1) {
		t.Errorf("With two levels, it should return a list with types %v, %v and %v; the return is: %v",
			n3,
			n2,
			n1,
			res)
	}

	res = _getTypeHierarchy(tD, tA)

	if res == nil {
		t.Error("With three levels, it returns nil")
	}

	if len(res) != 4 {
		t.Fatal("With three levels, it should return a list with length 4; the return is:", res)
	}

	if (res[0] != n4) || (res[1] != n1) || (res[2] != n2) || (res[3] != n1) {
		t.Errorf("With two levels, it should return a list with types %v, %v, %v and %v; the return is: %v",
			n4,
			n1,
			n2,
			n1,
			res)
	}

	res = _getTypeHierarchy(tE, tA)

	if res == nil {
		t.Error("With four levels, it returns nil")
	}

	if len(res) != 8 {
		t.Fatal("With four levels, it should return a list with length 8; the return is:", res)
	}

	if (res[0] != n5) ||
		(res[1] != n4) ||
		(res[2] != n1) ||
		(res[3] != n2) ||
		(res[4] != n1) ||
		(res[5] != n3) ||
		(res[6] != n2) ||
		(res[7] != n1) {
		t.Errorf("With two levels, it should return a list with types %v, %v, %v, %v, %v, %v, %v and %v; the return is: %v",
			n5,
			n4,
			n1,
			n2,
			n1,
			n3,
			n2,
			n1,
			res)
	}
}
