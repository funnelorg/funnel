// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run

import (
	"github.com/funnelorg/funnel/parse"
	"strconv"
)

type defaultScope struct {
	*Runner
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
	case "int":
		return d.intf
	case "+":
		return d.sumf
	}
	if s, ok := key.(string); ok {
		return E("unknown identifier: " + s)
	}
	return E("unknown identifier")
}

func (d defaultScope) runArgs(s Scope, args []parse.Node) []interface{} {
	result := make([]interface{}, len(args))
	for kk := range args {
		result[kk] = d.run(s, args[kk])
	}
	return result
}

func (d defaultScope) runRawArgs(s Scope, args []parse.Node) []interface{} {
	result := make([]interface{}, len(args))
	for kk := range args {
		result[kk] = d.runRaw(s, args[kk])
	}
	return result
}

func (d defaultScope) errorf(s Scope, args []parse.Node) interface{} {
	params := d.runRawArgs(s, args)
	if len(params) > 0 {
		switch s := params[0].(type) {
		case string:
			return E(s)
		case error:
			return s
		}
	}
	return ErrUnknown
}

func (d defaultScope) dotf(s Scope, args []parse.Node) interface{} {
	result := d.run(s, args[0])
	for _, next := range d.runRawArgs(s, args[1:]) {
		switch sx := result.(type) {
		case Scope:
			result = sx.Get(next)
		case error:
			return sx
		default:
			return E("cannot use dot with non-map")
		}
	}
	return result
}

func (d defaultScope) intf(s Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return E("int expects a single arg")
	}
	result := d.runRaw(s, args[0])
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
	return E("int: unknown argument type")
}

func (d defaultScope) sumf(s Scope, args []parse.Node) interface{} {
	sum := 0
	for _, v := range d.runArgs(s, args) {
		switch v := v.(type) {
		case int:
			sum += v
		case error:
			return v
		default:
			return E("sum: only works with ints")
		}
	}
	return sum
}
