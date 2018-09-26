// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data_test

import (
	"fmt"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/data"
	"github.com/funnelorg/funnel/run"
	"reflect"
	"testing"
)

var mapCases = map[string]interface{}{
	// map creation
	"data:map()":      "map: needs 1 arg at file:8",
	"data:map(1)":     "map: not a map at file:9",
	"data:map{x = 2}": "map[x:{2}]",

	// map count
	"data:map{x = 2}.count()": builtin.Number{1},

	// map filter
	"data:map{x = 2, y = 3}.filter()":           "filter: requires condition function at file:29",
	"data:map{x=2,y=3}.filter(fun(i,v, v < 3))": "map[x:{2}]",
	"data:map{x=2}.filter(1)":                   "filter: arg is not a function at file:21",
	"data:map{x=2}.filter(fun(i,v,v))":          "filter: function returned non-boolean at file:26",
	"data:map{x=2}.filter(fun(i,v,v ++ v))":     "missing term at file:32",

	// map map
	"data:map{x=2}.map()":               "map: requires 1 arg at file:17",
	"data:map{x=2}.map(fun(i,v,v+1))":   "map[x:{3}]",
	"data:map{x=2}.map(builtin:number)": "map[x:builtin:number: must have 1 arg]",

	// unknown field
	"data:map{x=2}.wat":       "wat: no such key at file:13",
	"data:map{x=2}.({x = 2})": "No such key at file:13",

	// testing Wraps
	"data:map{x={y=2}}.map(fun(i,v,v.count()))":         "map[x:{1}]",
	"data:map{x=data:map{y=2}}.map(fun(i,v,v.count()))": "map[x:{1}]",
	"data:map{x=a}.map(fun(i,v,v.item(0)))":             "map[x:hello]",
	"data:map{x=c}.map(fun(i,v,v.count()))":             "map[x:{1}]",

	// nested access
	"data:map{x=2}.x": builtin.Number{2},
}

func TestMap(t *testing.T) {
	vars := []map[interface{}]interface{}{{
		"a": []interface{}{"hello"},
		"b": map[interface{}]interface{}{"zoo": "boo"},
		"c": map[string]interface{}{"boo": "zoo"},
	}}
	s := run.NewScope(vars, data.Scope(builtin.Scope))
	r := &run.Runner{}

	for code, expected := range mapCases {
		t.Run(code, func(t *testing.T) {
			value := native(r.Eval(s, "file", code))
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

func native(v interface{}) interface{} {
	s, ok := v.(scopeWithKeys)
	if !ok {
		return v
	}
	result := map[interface{}]interface{}{}
	s.ForEachKeys(func(key interface{}) bool {
		result[key] = native(s.Get(key))
		return false
	})
	return result
}

type scopeWithKeys interface {
	Get(key interface{}) interface{}
	ForEachKeys(fn func(interface{}) bool)
}
