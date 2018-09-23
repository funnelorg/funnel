// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime_test

import (
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/runtime"
	"testing"
)

var m = map[interface{}]interface{}{
	"f": func(args ...int) int {
		return args[0] + args[1]
	},
	"nil": func() {},
	"zf": func() string {
		return "z"
	},
	"x": 1,
	"y": 10,
	"invalidf": func() (int, int) {
		return 0, 0
	},
	"float3": float64(3.0),
}

var r = &run.Runner{}

func TestNewScope(t *testing.T) {
	s := runtime.NewScope(m, runtime.DefaultScope)
	x := r.Eval(s, "code", "f(x, y)")
	if x != 11 {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "nil()")
	if x != nil {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "nil(num(22))")
	if err, ok := x.(error); !ok || err.Error() != "invalid function call" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "invalidf()")
	if err, ok := x.(error); !ok || err.Error() != "invalid function" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "zzz")
	if err, ok := x.(error); !ok || err.Error() != "no such key" {
		t.Error("unexpected result", x)
	}
}

func TestError(t *testing.T) {
	s := runtime.NewScope(m, runtime.DefaultScope)

	x := r.Eval(s, "code", "! boo")
	if err, ok := x.(error); !ok || err.Error() != "boo" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "!()")
	if err, ok := x.(error); !ok || err.Error() != "error" {
		t.Error("unexpected result", x)
	}
}

func TestDot(t *testing.T) {
	s := runtime.NewScope(m, runtime.DefaultScope)

	x := r.Eval(s, "code", "{z = y}.z")
	if x != 10 {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "{z = y}.(zf())")
	if x != 10 {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "{z = y}.(invalidf())")
	if err, ok := x.(error); !ok || err.Error() != "invalid function" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "{z = y}.goop.boop")
	if err, ok := x.(error); !ok || err.Error() != "no such key" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "{z = y}.z.boop")
	if err, ok := x.(error); !ok || err.Error() != "cannot use dot with non-map" {
		t.Error("unexpected result", x)
	}
}

func TestNum(t *testing.T) {
	s := runtime.NewScope(m, runtime.DefaultScope)

	x := r.Eval(s, "code", "num(2, 3)")
	if err, ok := x.(error); !ok || err.Error() != "num: incorrect number of args" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num()")
	if err, ok := x.(error); !ok || err.Error() != "num: incorrect number of args" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(float3)")
	if x != (runtime.Number{3.0}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(float3)")
	if x != (runtime.Number{3.0}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(num(float3))")
	if x != (runtime.Number{3.0}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(num())")
	if err, ok := x.(error); !ok || err.Error() != "num: incorrect number of args" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(y)")
	if err, ok := x.(error); !ok || err.Error() != "num: invalid arg type" {
		t.Error("unexpected result", x)
	}
}

func TestSum(t *testing.T) {
	s := runtime.NewScope(m, runtime.DefaultScope)

	x := r.Eval(s, "code", "1 + 2")
	if x != (runtime.Number{3}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "1 + float3")
	if x != (runtime.Number{4}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "num(1) + float3")
	if x != (runtime.Number{4}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "1 + num(y)")
	if err, ok := x.(error); !ok || err.Error() != "num: invalid arg type" {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "1 + x")
	if err, ok := x.(error); !ok || err.Error() != "sum: not a number" {
		t.Error("unexpected result", x)
	}
}

func TestFun(t *testing.T) {
	s := runtime.DefaultScope
	x := r.Eval(s, "code", "fun()()")
	if x != nil {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "{x = fun(x, y, x + y), y = 42, z = x(num(2), num(3))}.z")
	if x != (runtime.Number{5}) {
		t.Error("unexpected result", x)
	}

	x = r.Eval(s, "code", "fun(x+y,x+2y)")
	if err, ok := x.(error); !ok || err.Error() != "fun: invalid param name" {
		t.Error("unexpected result", x)
	}
}
