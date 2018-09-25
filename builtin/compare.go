// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func compare(s run.Scope, args []parse.Node, cmp func(n1, n2 Number) bool) interface{} {
	var last Number
	r := &run.Runner{}
	for kk, arg := range args {
		switch v := r.Run(s, arg).(type) {
		case Number:
			if kk > 0 && !cmp(last, v) {
				return false
			}
			last = v
		case error:
			return v
		default:
			e := &run.ErrorStack{Message: "not a number"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return true
}

func more(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F > n2.F })
}

func less(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F < n2.F })
}

func gte(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F >= n2.F })
}

func lte(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F <= n2.F })
}

func equals(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F == n2.F })
}

func notEquals(s run.Scope, args []parse.Node) interface{} {
	return compare(s, args, func(n1, n2 Number) bool { return n1.F != n2.F })
}

func and(s run.Scope, args []parse.Node) interface{} {
	r := &run.Runner{}
	for _, arg := range args {
		switch t := r.Run(s, arg).(type) {
		case bool:
			if !t {
				return false
			}
		case error:
			return t
		default:
			e := &run.ErrorStack{Message: "not a condition"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return true
}

func or(s run.Scope, args []parse.Node) interface{} {
	r := &run.Runner{}
	for _, arg := range args {
		switch t := r.Run(s, arg).(type) {
		case bool:
			if t {
				return true
			}
		case error:
			return t
		default:
			e := &run.ErrorStack{Message: "not a condition"}
			return run.WrapError(e, arg.Loc())
		}
	}
	return false
}
