// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"strconv"
)

// Number represents a generic number
type Number struct {
	F float64
}

// Num converts a number-like item to a number
func Num(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("num: incorrect number of args")
	}

	if args[0].Token != nil {
		ff, err := strconv.ParseFloat(args[0].Token.S, 64)
		if err == nil {
			return Number{ff}
		}
		return err
	}

	result := (&run.Runner{}).Run(s, args[0])
	switch f := result.(type) {
	case float64:
		return Number{f}
	case Number, error:
		return result
	}
	return errors.New("num: invalid arg type")
}

// Sum adds Number-like items together
func Sum(s run.Scope, args []parse.Node) interface{} {
	var sum float64
	for _, n := range args {
		result := (&run.Runner{}).Run(s, n)
		switch f := result.(type) {
		case float64:
			sum += f
		case Number:
			sum += f.F
		case error:
			return result
		default:
			return errors.New("sum: not a number")
		}
	}
	return Number{sum}
}

// Diff substracts Number-like items together
func Diff(s run.Scope, args []parse.Node) interface{} {
	var sum float64
	factor := 1.0

	for _, n := range args {
		result := (&run.Runner{}).Run(s, n)
		switch f := result.(type) {
		case float64:
			sum += factor * f
		case Number:
			sum += factor * f.F
		case error:
			return result
		default:
			return errors.New("diff: not a number")
		}
		factor = -1.0
	}
	return Number{sum}
}

// Multiply multiplie the Number-like items together
func Multiply(s run.Scope, args []parse.Node) interface{} {
	product := 1.0
	for _, n := range args {
		result := (&run.Runner{}).Run(s, n)
		switch f := result.(type) {
		case float64:
			product = product * f
		case Number:
			product = product * f.F
		case error:
			return result
		default:
			return errors.New("mult: not a number")
		}
	}
	return Number{product}
}

// Divide substracts Number-like items together
func Divide(s run.Scope, args []parse.Node) interface{} {
	product := 1.0
	for kk, n := range args {
		result := (&run.Runner{}).Run(s, n)
		item := 1.0
		switch f := result.(type) {
		case float64:
			item = f
		case Number:
			item = f.F
		case error:
			return result
		default:
			return errors.New("div: not a number")
		}
		if kk == 0 {
			product = item
		} else {
			product = product / item
		}

	}
	return Number{product}
}
