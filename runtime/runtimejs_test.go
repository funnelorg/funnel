// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build gopherjs

package runtime_test

import (
	"github.com/funnelorg/funnel/runtime"
	"testing"
)

var mm = map[interface{}]interface{}{
	"f": func(args ...interface{}) interface{} {
		return args[0].(int) + args[1].(int)
	},
	"nil": func(args ...interface{}) interface{} { return nil },
	"zf": func(args ...interface{}) interface{} {
		return "z"
	},
	"x": 1,
	"y": 10,
	"invalidf": func() (int, int) {
		return 0, 0
	},
}

func TestNewScopeJS(t *testing.T) {
	s := runtime.NewScope(mm, runtime.DefaultScope)
	x := r.Eval(s, "code", "f(x, y)")
	if x != 11 {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "nil()")
	if x != nil {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "zf()")
	if x != "z" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "invalidf()")
	if err, ok := x.(error); !ok || err.Error() != "not a function" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "zzz")
	if err, ok := x.(error); !ok || err.Error() != "no such key" {
		t.Error("unexpected result", x)
	}
}

func init() {
	// update the m map so that it works with the JS version of Function
	m["zf"] = mm["zf"]
	m["float3f"] = func(args ...interface{}) interface{} {
		return 3.0
	}
	delete(m, "invalidf")
}
