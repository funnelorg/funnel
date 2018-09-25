// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package url implements some helpful url functions
package url

import "github.com/funnelorg/funnel/run"

// Scope returns a new scope with the default url function
func Scope(base run.Scope) run.Scope {
	return run.NewScope([]map[interface{}]interface{}{Map}, base)
}

// Map contains the functions provided by this package
var Map = map[interface{}]interface{}{"url": run.ArgsResolver(urlf)}

func urlf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "url: requires exactly 1 arg"}
	}

	switch f := args[0].(type) {
	case string:
		return URL(f)
	case URL:
		return f
	case error:
		return f
	}
	return &run.ErrorStack{Message: "url: not a string"}
}

// URL represents a URL instance
type URL string

// Get returns the "methods" of URL
func (u URL) Get(key interface{}) interface{} {
	if key != "json" {
		return &run.ErrorStack{Message: "no such key"}
	}
	return u.json
}
