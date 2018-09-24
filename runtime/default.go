// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package runtime implements the runtime methods and helpers for use
// with the funnel language
package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/run"
)

type rtscope struct {
	m    map[interface{}]interface{}
	base run.Scope
}

func (r *rtscope) Get(key interface{}) interface{} {
	if v, ok := r.m[key]; ok {
		return v
	}
	return r.base.Get(key)
}

var def = map[interface{}]interface{}{
	"!":              Error,
	".":              Dot,
	"+":              Sum,
	"-":              Diff,
	"*":              Multiply,
	"/":              Divide,
	"builtin:number": Num,
	"builtin:string": String,
	"fun":            Fun,
}

// DefaultScope defines the default runtime methods
var DefaultScope = NewScope(def, empty{})

type empty struct{}

func (e empty) Get(key interface{}) interface{} {
	return errors.New("no such key")
}
