// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"strconv"
)

func listf(args []interface{}) interface{} {
	return list(args)
}

type list []interface{}

func (t list) Get(key interface{}) interface{} {
	switch key {
	case "item":
		return run.ArgsResolver(t.itemf)
	case "filter":
		return run.ArgsResolver(t.filterf)
	case "map":
		return run.ArgsResolver(t.mapf)
	case "slice":
		return run.ArgsResolver(t.slicef)
	case "count":
		return t.countf
	case "splice":
		return run.ArgsResolver(t.splicef)
	}
	if s, ok := key.(string); ok {
		return &run.ErrorStack{Message: "unknown field: " + s}
	}
	return &run.ErrorStack{Message: "unknown field"}
}

func (t list) itemf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "item: requires 1 arg"}
	}

	idx := 0
	if n, ok := args[0].(builtin.Number); ok {
		idx = int(n.F)
	} else {
		return &run.ErrorStack{Message: "item: not a number"}
	}

	if idx < 0 || idx >= len(t) {
		return &run.ErrorStack{Message: "item: out of bounds"}
	}
	return t[idx]
}

func (t list) slicef(args []interface{}) interface{} {
	start, end := 0, 0
	if len(args) != 2 {
		return &run.ErrorStack{Message: "slice: requires 2 args"}
	}
	if n, ok := args[0].(builtin.Number); ok {
		start = int(n.F)
	} else {
		return &run.ErrorStack{Message: "slice: not a number"}
	}

	if n, ok := args[1].(builtin.Number); ok {
		end = int(n.F)
	} else {
		return &run.ErrorStack{Message: "slice: not a number"}
	}

	if end > len(t) || start < 0 || end < start {
		return &run.ErrorStack{Message: "slice: out of bounds"}
	}

	return t[start:end]
}

func (t list) countf(s run.Scope, args []parse.Node) interface{} {
	return builtin.Number{float64(len(t))}
}

func (t list) splicef(args []interface{}) interface{} {
	start, end := 0, 0
	if len(args) != 3 {
		return &run.ErrorStack{Message: "splice: requires 3 args"}
	}
	if n, ok := args[0].(builtin.Number); ok {
		start = int(n.F)
	} else {
		return &run.ErrorStack{Message: "splice: not a number"}
	}

	if n, ok := args[1].(builtin.Number); ok {
		end = int(n.F)
	} else {
		return &run.ErrorStack{Message: "splice: not a number"}
	}

	if end > len(t) || start < 0 || start > end {
		return &run.ErrorStack{Message: "splice: out of bounds"}
	}

	if rep, ok := args[2].(list); ok {
		return append(append(t[0:start:start], rep...), t[end:]...)
	}
	return &run.ErrorStack{Message: "splice: arg 3 must be a list"}
}

func (t list) filterf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "filter: requires condition function"}
	}

	result := []interface{}(nil)
	for kk, elt := range t {
		params := []interface{}{builtin.Number{float64(kk)}, elt}
		switch check := t.invoke("filter", args[0], params).(type) {
		case bool:
			if check {
				result = append(result, elt)
			}
		case error:
			return check
		default:
			return &run.ErrorStack{Message: "filter: function returned non boolean"}
		}
	}
	return list(result)
}

func (t list) mapf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "map: requires one arg"}
	}

	result := []interface{}(nil)
	for kk, elt := range t {
		params := []interface{}{builtin.Number{float64(kk)}, elt}
		result = append(result, t.invoke("map", args[0], params))
	}
	return list(result)
}

func (t list) invoke(kind string, fn interface{}, args []interface{}) interface{} {
	var f func(s run.Scope, args []parse.Node) interface{}

	switch fn := fn.(type) {
	case func([]interface{}) interface{}:
		return fn(args)
	case callResolved:
		return fn.CallResolved(args)
	case func(s run.Scope, args []parse.Node) interface{}:
		f = fn
	case callable:
		f = fn.Call
	default:
		return &run.ErrorStack{Message: kind + ": arg is not a function"}
	}

	m := []map[interface{}]interface{}{{}}
	nodes := make([]parse.Node, len(args))
	for kk, arg := range args {
		key := "arg: " + strconv.Itoa(kk)
		m[0][key] = arg
		nodes[kk] = parse.Node{Token: &parse.Token{&parse.Loc{}, key}}
	}
	return f(run.NewScope(m, nil), nodes)
}

type callResolved interface {
	CallResolved(args []interface{}) interface{}
}

type callable interface {
	Call(s run.Scope, args []parse.Node) interface{}
}
