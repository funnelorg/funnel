// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package math implements some common math runtime routines for the
// funnel language.
package math

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/runtime"
	"math"
	"strconv"
)

// New returns a scope with the default math routines
func New() run.Scope {
	m := map[interface{}]interface{}{
		"square": Square,
		"root":   Root,
	}

	math := map[interface{}]interface{}{
		"math": runtime.NewScope(m, runtime.DefaultScope),
	}

	return runtime.NewScope(math, runtime.DefaultScope)
}

// Square calculates the square of a number-like item
func Square(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("square: must have exactly 1 arg")
	}

	result := (&run.Runner{}).Run(s, args[0])
	switch f := result.(type) {
	case float64:
		return runtime.Number{f * f}
	case runtime.Number:
		return runtime.Number{f.F * f.F}
	case error:
		if args[0].Token != nil {
			ff, err := strconv.ParseFloat(args[0].Token.S, 64)
			if err == nil {
				return runtime.Number{ff * ff}
			}
		}
		return result
	}
	return errors.New("square: not a number")
}

// Root calculates the square of a number-like item
func Root(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("root: must have exactly 1 arg")
	}

	result := (&run.Runner{}).Run(s, args[0])
	switch f := result.(type) {
	case float64:
		return runtime.Number{math.Sqrt(f)}
	case runtime.Number:
		return runtime.Number{math.Sqrt(f.F)}
	case error:
		if args[0].Token != nil {
			ff, err := strconv.ParseFloat(args[0].Token.S, 64)
			if err == nil {
				return runtime.Number{math.Sqrt(ff)}
			}
		}
		return result
	}
	return errors.New("root: not a number")
}
