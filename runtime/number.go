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

	result := (&run.Runner{}).Run(s, args[0])
	switch f := result.(type) {
	case float64:
		return Number{f}
	case Number:
		return result
	case error:
		if args[0].Token != nil {
			ff, err := strconv.ParseFloat(args[0].Token.S, 64)
			if err == nil {
				return Number{ff}
			}
		}
		return f
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
			if n.Token != nil {
				ff, err := strconv.ParseFloat(n.Token.S, 64)
				if err == nil {
					sum += ff
					continue
				}
			}
			return result
		default:
			return errors.New("sum: not a number")
		}
	}
	return Number{sum}
}
