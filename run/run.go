// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package run implements an interpreter for the funnel language.
package run

import "github.com/funnelorg/funnel/parse"

// Runner evaluates an expression
type Runner struct{}

// Run evaluates a parsed expression
func (r *Runner) Run(s Scope, expr parse.Node) interface{} {
	return unwrapValue(r.LazyRun(s, expr))
}

// LazyRun is like run but it can return "deferred" values which all
// implement the Value() method.
func (r *Runner) LazyRun(s Scope, expr parse.Node) interface{} {
	return r.run(s, expr)
}

func (r *Runner) run(s Scope, n parse.Node) interface{} {
	switch {
	case n.Token != nil:
		return s.Get(n.Token.S)
	case n.Map != nil:
		return newMapScope(*n.Map, s, r)
	}

	fn := r.Run(s, n.Call.Nodes[0])
	args := n.Call.Nodes[1:]

	switch fn := fn.(type) {
	case func(s Scope, args []parse.Node) interface{}:
		return fn(s, args)
	case error:
		return fn
	}
	return ErrNotFunction
}

// Eval is like Run except it parses the string as needed
func (r *Runner) Eval(s Scope, fname, str string) interface{} {
	return r.Run(s, parse.Parse(fname, str))
}

func unwrapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case Lazy:
		return unwrapValue(v.Value())
	case map[interface{}]interface{}:
		result := map[interface{}]interface{}{}
		for key, value := range v {
			result[key] = unwrapValue(value)
		}
		return result
	}
	return v
}
