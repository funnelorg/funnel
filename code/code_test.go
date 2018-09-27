// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package code_test

import (
	"fmt"
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/code"
	"github.com/funnelorg/funnel/data"
	"reflect"
	"testing"
)

var codeURL = "https://raw.githubusercontent.com/funnelorg/funnel/master/x/tests/simple.fun"
var reimportURL = "https://raw.githubusercontent.com/funnelorg/funnel/master/x/tests/reimport.fun?q=3"
var importCall = `code:import(builtin:string "` + codeURL + `")`
var reimportCall = `code:import(builtin:string "` + reimportURL + `")`
var s = code.Scope(data.Scope(builtin.Scope))

var cases = map[string]interface{}{
	"code:import()":            "import: requires one arg at file:11",
	"code:import(1)":           "import: arg must be a string at file:12",
	"code:import(string boop)": `Get boop: unsupported protocol scheme ""`,
	importCall:                 "[hello world]",
	reimportCall:               "[hello world]",
}

func TestUrl(t *testing.T) {
	for code, expected := range cases {
		t.Run(code, func(t *testing.T) {
			value := native(funnel.Eval(s, "file", code))
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
