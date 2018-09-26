// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"fmt"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/data"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"reflect"
	"testing"
)

var cases = map[string]interface{}{
	// list count
	"data:list(1, 2, 3).count()": builtin.Number{3},

	// list item
	"data:list(1).item()":        "item: requires 1 arg at file:17",
	"data:list(1).item({x = 2})": "item: not a number at file:21",
	"data:list(1).item(-1)":      "item: out of bounds at file:18",
	"data:list(1,3).item(1)":     builtin.Number{3},

	// list slice
	"data:list(1, 2, 3).slice(1,2)": "[{2}]",
	"data:list(1).slice()":          "slice: requires 2 args at file:18",
	"data:list(1).slice({x=2},2)":   "slice: not a number at file:24",
	"data:list(1).slice(2,{x=2})":   "slice: not a number at file:20",
	"data:list(1).slice(-2,2)":      "slice: out of bounds at file:21",
	"data:list(1).slice(2,1)":       "slice: out of bounds at file:20",

	// list splice
	"data:list(1, 2, 3).splice(1,2,data:list(9))": "[{1} {9} {3}]",

	"data:list(1).splice()":          "splice: requires 3 args at file:19",
	"data:list(1).splice({x=2},2,1)": "splice: not a number at file:25",
	"data:list(1).splice(2,{x=2},1)": "splice: not a number at file:21",
	"data:list(1).splice(-2,2,1)":    "splice: out of bounds at file:22",
	"data:list(1).splice(2,1,1)":     "splice: out of bounds at file:21",
	"data:list(1).splice(0,1,1)":     "splice: arg 3 must be a list at file:21",

	// list filter
	"data:list(1).filter()":                "filter: requires condition function at file:19",
	"data:list(1).filter(fun(i,v, v < 5))": "[{1}]",
	"data:list(1).filter(1)":               "filter: arg is not a function at file:20",
	"data:list(1).filter(fun(i,v,v))":      "filter: function returned non boolean at file:25",
	"data:list(1).filter(fun(i,v,v ++ v))": "missing term at file:31",

	// list map
	"data:list(1).map()":               "map: requires one arg at file:16",
	"data:list(1).map(fun(i,v,v+1))":   data.List([]interface{}{builtin.Number{2}}),
	"data:list(1).map(builtin:number)": "[builtin:number: must have 1 arg]",

	// unknown field
	"data:list(1).wat":       "unknown field: wat at file:12",
	"data:list(1).({x = 2})": "unknown field at file:12",

	// special function types
	"data:list(1).map(raw)":      "[raw]",
	"data:list(1).map(callable)": "[callable]",
}

func TestList(t *testing.T) {
	fns := []map[interface{}]interface{}{{
		"raw": func([]interface{}) interface{} {
			return "raw"
		},
		"callable": callable{},
	}}
	s := run.NewScope(fns, data.Scope(builtin.Scope))
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

type callable struct{}

func (c callable) Call(s run.Scope, nodes []parse.Node) interface{} {
	return "callable"
}
