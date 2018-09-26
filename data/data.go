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
	"data:map":  run.ArgsResolver(mapf),
}

// List is the native version of data:list
func List(items []interface{}) interface{} {
	return list(items)
}

// Wrap converts any native type into a Map or List
func Wrap(v interface{}) interface{} {
	switch v := v.(type) {
	case mscope:
		return v
	case scopeWithKeys:
		return mscope{v}
	case []interface{}:
		return list(v)
	case map[interface{}]interface{}:
		mapped := run.NewScope([]map[interface{}]interface{}{v}, nil)
		return mscope{mapped.(scopeWithKeys)}
	case map[string]interface{}:
		result := map[interface{}]interface{}{}
		for key, value := range v {
			result[key] = value
		}
		mapped := run.NewScope([]map[interface{}]interface{}{result}, nil)
		return mscope{mapped.(scopeWithKeys)}
	default:
		return v
	}
}
