// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package funnel implements components the parser and interpreter for
// a very simple functional language
package funnel

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/runtime"
)

// Scope is the interface that defines the "scope" of a particular
// evaluation. This is a key-value store of the global identifiers
type Scope interface {
	Get(key interface{}) interface{}
	Value() interface{}
}

// Eval evaluates code using the provided scope. The filename is used
// for reporting errors.
func Eval(s Scope, filename, code string) interface{} {
	r := run.Runner{}
	if s == nil {
		s = runtime.DefaultScope
	}
	return r.Run(s, parse.Parse(filename, code))
}
