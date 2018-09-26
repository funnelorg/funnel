// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"strconv"
)

type defaultScope struct {
	*run.Runner
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
		return run.ArgsResolver(d.sumf)
	case "builtin:string":
		return d.stringf
	}
	if s, ok := key.(string); ok {
		return errors.New("unknown identifier: " + s)
	}
	return errors.New("unknown identifier")
}

func (d defaultScope) ForEachKeys(fn func(interface{}) bool) {
	// not implemented
}

func (d defaultScope) errorf(s run.Scope, args []parse.Node) interface{} {
	if len(args) > 0 {
		var param interface{}
		r := &run.Runner{}
		if args[0].Token != nil {
			param = args[0].Token.S
		} else {
			param = r.Run(s, args[0])
		}

		switch s := param.(type) {
		case string:
			return errors.New(s)
		case error:
			return s
		}
	}
	return errors.New("unknown error")
}

func (d defaultScope) dotf(s run.Scope, args []parse.Node) interface{} {
	result := d.LazyRun(s, args[0])
	for _, arg := range args[1:] {
		var next interface{}
		if arg.Token != nil {
			next = arg.Token.S
		} else {
			next = d.Run(s, arg)
		}

		switch sx := result.(type) {
		case run.Scope:
			result = sx.Get(next)
		case error:
			return sx
		default:
			return errors.New("cannot use dot with non-map")
		}
	}
	return result
}

func (d defaultScope) intf(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("int expects a single arg")
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
	return errors.New("int: unknown argument type")
}

func (d defaultScope) sumf(args []interface{}) interface{} {
	sum := 0
	for _, v := range args {
		switch v := v.(type) {
		case int:
			sum += v
		case error:
			return v
		default:
			return errors.New("sum: only works with ints")
		}
	}
	return sum
}

func (d defaultScope) stringf(s run.Scope, args []parse.Node) interface{} {
	return args[0].Token.S
}
