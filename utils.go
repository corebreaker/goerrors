package goerrors

import (
	"reflect"
)

// Concatenate 2 string lists
func _concat(l1, l2 []string) []string {
	if (l1 == nil) && (l2 == nil) {
		return nil
	}

	sz1 := len(l1)
	sz2 := len(l2)
	res := make([]string, sz1+sz2)

	if sz1 != 0 {
		copy(res, l1)
	}

	if sz2 != 0 {
		copy(res[sz1:], l2)
	}

	return res
}

// Get type hierarchy from an error type passed as `this_type` parameter.
// The `final_type` parameter represents type of `GoError` structure.
func _getTypeHierarchy(this_type, final_type reflect.Type) []string {
	if (this_type == nil) || (this_type.Kind() != reflect.Struct) {
		return []string{}
	}

	res := []string{this_type.PkgPath() + "." + this_type.Name()}

	if this_type == final_type {
		return res
	}

	list := make([]string, 0)

	n := this_type.NumField()
	for i := 0; i < n; i++ {
		field := this_type.Field(i)
		if !field.Anonymous {
			continue
		}

		parents := _getTypeHierarchy(field.Type, final_type)
		if len(parents) > 0 {
			list = _concat(list, parents)
		}
	}

	if len(list) == 0 {
		return list
	}

	return _concat(res, list)
}
