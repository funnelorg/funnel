// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package code implements funnel import, eval and related functions
package code

import (
	"github.com/funnelorg/funnel/run"
	funnelurl "github.com/funnelorg/funnel/url"
	"net/url"
)

// Scope returns the definitions for all the builtin functions.
// The provided scope will act like the scope for any imported code.
func Scope(base run.Scope) run.Scope {
	im := &importer{"", base}
	mm := []map[interface{}]interface{}{{
		"code:import": run.ArgsResolver(im.importf),
	}}

	return run.NewScope(mm, base)
}

type importer struct {
	url  string
	base run.Scope
}

func (im *importer) Get(key interface{}) interface{} {
	if key == "code:import" {
		return run.ArgsResolver(im.importf)
	}
	return im.base.Get(key)
}

func (im *importer) importf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "import: requires one arg"}
	}
	s, ok := args[0].(string)
	if !ok {
		return &run.ErrorStack{Message: "import: arg must be a string"}
	}

	resolved := im.resolve(s)
	resp, err := (funnelurl.URL(resolved)).Fetch(funnelurl.Text)
	if err != nil {
		return err
	}
	text := resp.(string)

	r := &run.Runner{}
	return r.Eval(&importer{resolved, im.base}, resolved, text)
}

func (im *importer) resolve(current string) string {
	if im.url != "" {
		u, err1 := url.Parse(current)
		base, err2 := url.Parse(im.url)
		if err1 == nil && err2 == nil {
			return base.ResolveReference(u).String()
		}
	}
	return current
}
