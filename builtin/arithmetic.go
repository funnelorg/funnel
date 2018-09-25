// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"strconv"
)

// Number is the builtin numeric type for funnel
type Number struct {
	F float64
}

func numberf(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "builtin:number: must have 1 arg"}
	}

	if args[0].Token == nil {
		return &run.ErrorStack{Message: "builtin:number: must be a number"}
	}

	ff, err := strconv.ParseFloat(args[0].Token.S, 64)
	if err != nil {
		return err
	}
	return Number{ff}
}

func add(s run.Scope, args []parse.Node) interface{} {
	var sum float64
	r := &run.Runner{}
	for _, arg := range args {
		switch v := r.Run(s, arg).(type) {
		case Number:
			sum += v.F
		case error:
			return v
		default:
			e := &run.ErrorStack{Message: "not a number"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return Number{sum}
}

func isMissingTerm(arg parse.Node) bool {
	return arg.Call != nil && arg.Call.IsError() && arg.Call.Nodes[1].Token.S == "missing term"
}

func sub(s run.Scope, args []parse.Node) interface{} {
	var sum float64
	r := &run.Runner{}
	for kk, arg := range args {
		if kk == 0 && isMissingTerm(arg) {
			continue
		}

		switch v := r.Run(s, arg).(type) {
		case Number:
			if kk == 0 {
				sum = v.F
			} else {
				sum -= v.F
			}
		case error:
			return v
		default:
			e := &run.ErrorStack{Message: "not a number"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return Number{sum}
}

func mult(s run.Scope, args []parse.Node) interface{} {
	result := Number{1.0}
	r := &run.Runner{}
	for _, arg := range args {
		switch v := r.Run(s, arg).(type) {
		case Number:
			result.F = result.F * v.F
		case error:
			return v
		default:
			e := &run.ErrorStack{Message: "not a number"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return result
}

func div(s run.Scope, args []parse.Node) interface{} {
	result := Number{1.0}
	r := &run.Runner{}
	for kk, arg := range args {
		switch v := r.Run(s, arg).(type) {
		case Number:
			if kk == 0 {
				result = v
			} else {
				result.F = result.F / v.F
			}
		case error:
			return v
		default:
			e := &run.ErrorStack{Message: "not a number"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return result
}
