// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

// Fun converts its last argument into a function. All but the last
// argument is considered a parameter to the arguemnt.
func Fun(s run.Scope, args []parse.Node) interface{} {
	if len(args) == 0 {
		return func(s run.Scope, args []parse.Node) interface{} {
			return nil
		}
	}

	params := []string{}
	for kk := range args[:len(args)-1] {
		if args[kk].Token == nil {
			return errors.New("fun: invalid param name")
		}
		params = append(params, args[kk].Token.S)
	}
	n := args[len(args)-1]
	return func(inner run.Scope, args []parse.Node) interface{} {
		r := &run.Runner{}
		values := make([]interface{}, len(args))
		for kk, arg := range args {
			values[kk] = r.LazyRun(inner, arg)
		}
		inv := invocation{params, values, s}
		return r.LazyRun(inv, n)
	}
}

type invocation struct {
	params []string
	args   []interface{}
	s      run.Scope
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
