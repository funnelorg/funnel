// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package builtin has the set of builtins for the funnel runtime.
package builtin

import "github.com/funnelorg/funnel/run"

// Scope contains the definitions for all the builtin functions.
var Scope = run.NewScope([]map[interface{}]interface{}{Map}, nil)

// Map contains the map of builtins
var Map = map[interface{}]interface{}{
	"!":  errorInternal,
	"+":  add,
	"-":  sub,
	"*":  mult,
	"/":  div,
	">":  more,
	">=": gte,
	"<":  less,
	"<=": lte,
	"==": equals,
	"!=": notEquals,
	"|":  or,
	"&":  and,
	"if": iff,
	".":  dot,

	"builtin:string": stringf,
	"builtin:number": numberf,
	"or":             or,
	"and":            and,

	"error": errorf,

	"fun": fun,
}
