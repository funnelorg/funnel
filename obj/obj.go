// Copyright (C) 2019 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package obj defines the basic semantics of an object and implements
// some simple objects
package obj

// O represents the interface to any object.
type O interface {
	// Invoke treats the object as a function and calls it with
	// the provided parameters. The caller context provides info
	// on where the call originated and is used for generating
	// the call stack with errors
	Invoke(caller O, args ...O) O

	// Get returns the named field. The caller context is used for
	// reporting errors
	Get(caller O, name string) O
}

// Nil returns a nil object which does not have any fields and cannot
// be invoked.  Attempting to access a field or invoke a method
// results in an error.
func Nil() O {
	return nilobj{}
}

type nilobj struct{}

func (n nilobj) Invoke(caller O, args ...O) O {
	return Error(caller, "<nil> cannot be invoked")
}

func (n nilobj) Get(caller O, key string) O {
	return Error(caller, "cannot access field "+key+" of <nil>")
}

// Error returns an error object.  Error objects support only the
// "stack" field.  Accessing anything else simply causes a chained
// error.
func Error(caller O, message string) O {
	return errobj{caller, nil, message}
}

type errobj struct {
	caller  O
	prev    *errobj
	message string
}

func (e errobj) Invoke(caller O, args ...O) O {
	return errobj{caller, &e, e.message}
}

func (e errobj) Get(caller O, key string) O {
	if key == "stack" {
		stack := ""
		if e.prev != nil {
			stack = string(e.prev.Get(caller, key).(S)) + "\n"
		}
		return S(stack + e.message)
	}

	return errobj{caller, &e, "cannot access field " + key + " of <error>"}
}

// Tuple creates a tuple with the provided fields.  Invoking the
// tuple effectively invokes the "()" field of the tuple
func Tuple(fields map[string]O) O {
	return tuple(fields)
}

type tuple map[string]O

func (t tuple) Invoke(caller O, args ...O) O {
	if o, ok := t["()"]; ok {
		return o.Invoke(caller, args...)
	}

	return Error(caller, "tuple is not callable")
}

func (t tuple) Get(caller O, key string) O {
	if o, ok := t[key]; ok {
		return o
	}

	return Error(caller, "tuple does not have field "+key)
}

// S wraps string with an empty implementation of O
type S string

// Invoke just panics.
func (s S) Invoke(caller O, args ...O) O {
	return Error(caller, "string is not callable")
}

// Get just panics. TODO: implement string methods?
func (s S) Get(caller O, name string) O {
	return Error(caller, "cannot access field "+name+" of <string>")
}
