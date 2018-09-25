// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func fun(s run.Scope, args []parse.Node) interface{} {
	if len(args) == 0 {
		return func(s run.Scope, args []parse.Node) interface{} {
			return nil
		}
	}

	params := []string{}
	for _, arg := range args[:len(args)-1] {
		if arg.Token == nil {
			e := &run.ErrorStack{Message: "fun: invalid param name"}
			return run.WrapError(e, arg.Loc())
		}
		params = append(params, arg.Token.S)
	}
	expr := args[len(args)-1]
	return &callable{params, expr, s}
}

type callable struct {
	params []string
	expr   parse.Node
	s      run.Scope
}

func (c *callable) Call(s run.Scope, args []parse.Node) interface{} {
	r := &run.Runner{}
	return c.CallResolved(r.ResolveArgs(s, args))
}

func (c *callable) CallResolved(args []interface{}) interface{} {
	r := &run.Runner{}
	return r.LazyRun(&invocation{c, args}, c.expr)
}

type invocation struct {
	*callable
	args []interface{}
}

func (i invocation) Get(key interface{}) interface{} {
	if str, ok := key.(string); ok {
		for kk, param := range i.params {
			if param == str {
				return i.args[kk]
			}
		}
	}
	return i.s.Get(key)
}
