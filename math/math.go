// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package math implements some common math runtime routines for the
// funnel language.
package math

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/run"
	"math"
)

// Scope defines a new scope for math based on the provided scope
func Scope(base run.Scope) run.Scope {
	return run.NewScope([]map[interface{}]interface{}{Map}, base)
}

// Map provides the set of functions that math defines
var Map = map[interface{}]interface{}{
	"math:square": run.ArgsResolver(square),
	"math:root":   run.ArgsResolver(root),
}

func square(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "math:square: must have exactly 1 arg"}
	}

	switch f := args[0].(type) {
	case builtin.Number:
		return builtin.Number{f.F * f.F}
	case error:
		return f
	}
	return &run.ErrorStack{Message: "math:square: not a number"}
}

func root(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "math:root: must have exactly 1 arg"}
	}

	switch f := args[0].(type) {
	case builtin.Number:
		return builtin.Number{math.Sqrt(f.F)}
	case error:
		return f
	}
	return &run.ErrorStack{Message: "math:root: not a number"}
}
