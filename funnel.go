// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package funnel implements the parser and interpreter for
// the funnel programming language.
//
// The funnel language is basically simple expressions:
//
//     x + 2
//     (2*x + 5)*y
//
// It also has functions and a dot notation:
//
//     math.square(x)
//
// It also has let/expressions for creating more complicated expressions:
//
//    {l = 5, h = 10, hypo = math.root(l*l + h*h)}.hypo
//
// Let expressions can be given in any order.
//
// Functions can be defined via fun (which is just a regular function)
//
//    {f = fun(x, y, x+y), z = f(2, 3)}.z
package funnel

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/runtime"
)

// Scope is the interface that defines the "scope" of a particular
// evaluation. This is a key-value store of the global identifiers
type Scope interface {
	Get(key interface{}) interface{}
	Value() interface{}
}

// Eval evaluates code using the provided scope. The filename is used
// for reporting errors.
//
// The default scope uses  the runtime which only defines "." (object
// access), "!" (error function) and "+" (for arithmetic).
//
// Custom scopes can be defined like so:
//
//      m := map[interface{}]interface{} {
//           "x": 5.0,
//           "square": func(n runtime.Number) float64 {
//               return n.F * n.F
//            },
//      }
//      myScope := runtime.NewScope(m, runtime.DefaultScope)
//
//      x := funnel.Eval(myScope, "myfile.go", "square(x + 2)")
//      // x == 49
//
func Eval(s Scope, filename, code string) interface{} {
	r := run.Runner{}
	if s == nil {
		s = runtime.DefaultScope
	}
	return r.Run(s, parse.Parse(filename, code))
}
