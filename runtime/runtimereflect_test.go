// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js

package runtime_test

import (
	"github.com/funnelorg/funnel/runtime"
	"testing"
)

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

	x = r.Eval(s, "code", "nil(22)")
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
