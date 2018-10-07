// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package wiki_test

import (
	"fmt"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/wiki"
	"reflect"
	"testing"
)

var cases = map[string]interface{}{
	"wiki:table(string x, string y)": "wiki:tables: requires exactly 1 arg at file:19",
	"wiki:table(1 ++ 2)":             "missing term at file:14",
	"wiki:table(1)":                  "wiki:tables: not a string at file:11",

	"wiki:table(string \"<>\")": "Parse failure: text at file:18",

	"wiki:table(string \"List of prime numbers\").count()":                 "{50}",
	`wiki:table(string "Countries of the United Kingdom").map(value.Name)`: `[England Northern Ireland Scotland Wales United Kingdom]`,
}

func TestWiki(t *testing.T) {
	s := wiki.Scope(builtin.Scope)
	r := &run.Runner{}

	for code, expected := range cases {
		t.Run(code, func(t *testing.T) {
			value := r.Eval(s, "file", code)
			if err, ok := value.(error); ok {
				if err.Error() != expected {
					t.Error("Failed", value)
				}
			} else if !reflect.DeepEqual(value, expected) {
				if fmt.Sprintf("%v", value) != expected {
					t.Error("Failed", value)
				}
			}
		})
	}
}
