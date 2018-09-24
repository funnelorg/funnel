// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"reflect"
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
	r := &run.Runner{}
	fnv := reflect.ValueOf(fn)
	return func(s run.Scope, args []parse.Node) (output interface{}) {
		// TODO: replace panic/recover with proper type checks
		defer func() {
			if r := recover(); r != nil {
				output = errors.New("invalid function call")
			}
		}()

		v := make([]reflect.Value, len(args))
		for kk := range args {
			v[kk] = reflect.ValueOf(r.Run(s, args[kk]))
		}
		result := fnv.Call(v)
		switch len(result) {
		case 0:
			return nil
		case 1:
			return result[0].Interface()
		}
		return errors.New("invalid function")
	}
}

// NewScope takes a map of values and creates a scope out of it. If
// any of them are functions, it wraps them with Function if
// necessary.
func NewScope(m map[interface{}]interface{}, base run.Scope) run.Scope {
	wrapped := map[interface{}]interface{}{}
	for k, v := range m {
		wrapped[k] = v
		if reflect.ValueOf(v).Kind() == reflect.Func {
			_, ok := v.(func(s run.Scope, args []parse.Node) interface{})
			if !ok {
				wrapped[k] = Function(v)
			}
		}
	}
	return &rtscope{wrapped, base}
}
