// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build gopherjs,!jsreflect

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

// Function converts any typed function into the type required by the
// runner.  On GopherJS builds, this does not use reflection and is
// restricted to only two types of functions:
//      func(s Scope, args []parse.Node) interface{}
//      func(args ...interface{}) interface{}
//
// Full support can be enabled for GopherJS builds by using the build
// tag jsreflect
func Function(fn interface{}) func(s run.Scope, args []parse.Node) interface{} {
	if f, ok := fn.(func(s run.Scope, args []parse.Node) interface{}); ok {
		return f
	}

	r := &run.Runner{}
	return func(s run.Scope, args []parse.Node) (output interface{}) {
		// TODO: replace panic/recover with proper type checks
		defer func() {
			if r := recover(); r != nil {
				output = errors.New("invalid function call")
			}
		}()

		v := make([]interface{}, len(args))
		for kk := range args {
			v[kk] = r.LazyRun(s, args[kk])
		}
		return fn.(func(args ...interface{}) interface{})(v...)
	}
}

// NewScope takes a map of values and creates a scope out of it. If
// any of them are functions, it wraps them with Function if
// necessary.
func NewScope(m map[interface{}]interface{}, base run.Scope) run.Scope {
	wrapped := map[interface{}]interface{}{}
	for k, v := range m {
		wrapped[k] = v
		if _, ok := v.(func(args ...interface{}) interface{}); ok {
			wrapped[k] = Function(v)
		}
	}
	return &rtscope{wrapped, base}
}
