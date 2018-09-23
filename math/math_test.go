// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package math_test

import (
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/math"
	"github.com/funnelorg/funnel/runtime"
	"testing"
)

func TestSquare(t *testing.T) {
	x := funnel.Eval(math.New(), "boo", "math.square(num(5))")
	if x != (runtime.Number{25}) {
		t.Error("unexpected result", x)
	}

	s := runtime.NewScope(map[interface{}]interface{}{"x": 5.0}, math.New())
	x = funnel.Eval(s, "boo", "math.square(x)")
	if x != (runtime.Number{25}) {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.square(5)")
	if x != (runtime.Number{25}) {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.square()")
	if err, ok := x.(error); !ok || err.Error() != "square: must have exactly 1 arg" {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.square(num())")
	if err, ok := x.(error); !ok || err.Error() != "num: incorrect number of args" {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.square(math)")
	if err, ok := x.(error); !ok || err.Error() != "square: not a number" {
		t.Error("unexpected result", x)
	}
}

func TestRoot(t *testing.T) {
	x := funnel.Eval(math.New(), "boo", "math.root(num(25))")
	if x != (runtime.Number{5}) {
		t.Error("unexpected result", x)
	}

	s := runtime.NewScope(map[interface{}]interface{}{"x": 25.0}, math.New())
	x = funnel.Eval(s, "boo", "math.root(x)")
	if x != (runtime.Number{5}) {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.root(25)")
	if x != (runtime.Number{5}) {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.root()")
	if err, ok := x.(error); !ok || err.Error() != "root: must have exactly 1 arg" {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.root(num())")
	if err, ok := x.(error); !ok || err.Error() != "num: incorrect number of args" {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(math.New(), "boo", "math.root(math.root)")
	if err, ok := x.(error); !ok || err.Error() != "root: not a number" {
		t.Error("unexpected result", x)
	}
}
