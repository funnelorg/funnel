// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"strconv"
)

type defaultScope struct {
	*run.Runner
}

func (d defaultScope) Value() interface{} {
	panic("internal error")
}

func (d defaultScope) Get(key interface{}) interface{} {
	switch key {
	case "!":
		return d.errorf
	case ".":
		return d.dotf
	case "builtin:number":
		return d.intf
	case "+":
		return d.sumf
	case "builtin:string":
		return d.stringf
	}
	if s, ok := key.(string); ok {
		return run.E("unknown identifier: " + s)
	}
	return run.E("unknown identifier")
}

func (d defaultScope) runArgs(s run.Scope, args []parse.Node) []interface{} {
	result := make([]interface{}, len(args))
	for kk := range args {
		result[kk] = d.LazyRun(s, args[kk])
	}
	return result
}

func (d defaultScope) runRawArgs(s run.Scope, args []parse.Node) []interface{} {
	result := make([]interface{}, len(args))
	for kk := range args {
		if args[kk].Token != nil {
			result[kk] = args[kk].Token.S
		} else {
			result[kk] = d.LazyRun(s, args[kk])
		}
	}
	return result
}

func (d defaultScope) errorf(s run.Scope, args []parse.Node) interface{} {
	params := d.runRawArgs(s, args)
	if len(params) > 0 {
		switch s := params[0].(type) {
		case string:
			return run.E(s)
		case error:
			return s
		}
	}
	return run.ErrUnknown
}

func (d defaultScope) dotf(s run.Scope, args []parse.Node) interface{} {
	result := d.LazyRun(s, args[0])
	for _, next := range d.runRawArgs(s, args[1:]) {
		switch sx := result.(type) {
		case run.Scope:
			result = sx.Get(next)
		case error:
			return sx
		default:
			return run.E("cannot use dot with non-map")
		}
	}
	return result
}

func (d defaultScope) intf(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return run.E("int expects a single arg")
	}
	var result interface{}
	if args[0].Token != nil {
		result = args[0].Token.S
	} else {
		result = d.Run(s, args[0])
	}
	switch result := result.(type) {
	case string:
		v, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		return v
	case int, error:
		return result
	}
	return run.E("int: unknown argument type")
}

func (d defaultScope) sumf(s run.Scope, args []parse.Node) interface{} {
	sum := 0
	for _, v := range d.runArgs(s, args) {
		switch v := v.(type) {
		case int:
			sum += v
		case error:
			return v
		default:
			return run.E("sum: only works with ints")
		}
	}
	return sum
}

func (d defaultScope) stringf(s run.Scope, args []parse.Node) interface{} {
	return args[0].Token.S
}
