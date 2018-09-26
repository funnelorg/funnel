// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package data implements list, filter, map and reduce
package data

import "github.com/funnelorg/funnel/run"

// Scope returns the definitions for all the builtin functions.
func Scope(base run.Scope) run.Scope {
	return run.NewScope([]map[interface{}]interface{}{Map}, base)
}

// Map contains the map of builtins
var Map = map[interface{}]interface{}{
	"data:list": run.ArgsResolver(listf),
}

// List is the native version of data:list
func List(items []interface{}) interface{} {
	return list(items)
}
