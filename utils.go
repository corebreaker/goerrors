package goerrors

import (
    "fmt"
    "reflect"
    "runtime"
)

func _list_len(l []string) int {
    if l == nil {
        return 0
    }

    return len(l)
}

// Concatenate 2 string lists
func _concat(l1, l2 []string) []string {
    sz1 := _list_len(l1)
    sz2 := _list_len(l2)

    if sz1 == 0 {
        return l2
    }

    if sz2 == 0 {
        return l1
    }

    res := make([]string, sz1+sz2)

    copy(res, l1)
    copy(res[sz1:], l2)

    return res
}

// Get parent types from an error type passed as `this_type` parameter.
// The `final_type` parameter represents type of `GoError` structure.
func _get_parents(this_type, final_type reflect.Type) []string {
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

        parents := _get_parents(field.Type, final_type)
        if len(parents) > 0 {
            list = _concat(list, parents)
        }
    }

    if len(list) == 0 {
        return list
    }

    return _concat(res, list)
}

// Build a stack trace
func _make_stacktrace(prune_levels int) []string {
    // Checks that pruned levels passed in parameted is really a positive value
    // A negative value means, no pruning
    if prune_levels < 0 {
        prune_levels = 0
    }

    // Program Counter initialisation
    var pc uintptr = 1

    // Resulting stack trace
    trace := make([]string, 0)

    // Populate the stack trace
    for i := prune_levels + 2; pc != 0; i++ {
        // Retreive runtime informations
        ptr, file, line, ok := runtime.Caller(i)
        pc = ptr

        // If there isn't significant information, go to next level
        if (pc == 0) || (!ok) {
            continue
        }

        // Retreive called function
        f := runtime.FuncForPC(pc)

        // Add stack trace entry
        trace = append(trace, fmt.Sprintf("%s (%s:%d)", f.Name(), file, line))
    }

    // Returns stack trace
    return trace
}
