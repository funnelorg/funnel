// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package url implements some helpful url functions
package url

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/runtime"
)

// New returns a new scope with the default url function
func New(base run.Scope) run.Scope {
	url := map[interface{}]interface{}{"url": urlf}
	return runtime.NewScope(url, base)
}

func urlf(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("url: requries exactly 1 arg")
	}

	result := (&run.Runner{}).Run(s, args[0])
	switch f := result.(type) {
	case string:
		return URL(f)
	case URL:
		return f
	case error:
		if args[0].Token != nil {
			return URL(args[0].Token.S)
		}
		return result
	}
	return errors.New("url: not a string")
}

// URL represents a URL instance
type URL string

// Get returns the "methods" of URL
func (u URL) Get(key interface{}) interface{} {
	if key != "json" {
		return errors.New("no such key")
	}
	return u.json
}
