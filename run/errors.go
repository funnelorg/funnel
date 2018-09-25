// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run

import (
	"github.com/funnelorg/funnel/parse"
	"strconv"
)

// ErrorStack implements error with call stack
type ErrorStack struct {
	Message string
	File    string
	Offset  int
	Inner   error
}

// Error returns the inner message if one exists. Otherwise it returns
// the local message.  Use Stack() to get the full stack
func (es *ErrorStack) Error() string {
	if es.Inner != nil {
		return es.Inner.Error()
	}
	return es.format()
}

// Stack returns the error stack, formatted as a string
func (es *ErrorStack) Stack() string {
	stack := ""
	if es.Inner != nil {
		if s, ok := es.Inner.(stackable); ok {
			stack = s.Stack()
		} else {
			stack = es.Inner.Error() + "\n"
		}
	}
	return stack + es.format() + "\n"
}

func (es *ErrorStack) format() string {
	suffix := ""

	if es.File != "" || es.Offset != 0 {
		suffix = " at " + es.File + ":" + strconv.Itoa(es.Offset)
	}
	return es.Message + suffix
}

type stackable interface {
	Stack() string
}

// WrapError checks for an error and if so, either updates the
// location of the error (if it has no location) or adds another entry
// to the stack with the provided location
func WrapError(x interface{}, l *parse.Loc) interface{} {
	if es, ok := x.(*ErrorStack); ok {
		if es.Offset == 0 && es.File == "" {
			es.File, es.Offset = l.File, l.Offset
		}
		return es
	}
	if e, ok := x.(error); ok {
		return &ErrorStack{File: l.File, Offset: l.Offset, Inner: e}
	}
	return x
}
