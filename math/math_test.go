// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package math_test

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/math"
	"github.com/funnelorg/funnel/run"
	"testing"
)

var cases = map[string]interface{}{
	// square
	"math:square(5)":     builtin.Number{25},
	"math:square()":      "math:square: must have exactly 1 arg at file:11",
	"math:square(2++2)":  "missing term at file:14",
	"math:square({x=2})": "math:square: not a number at file:14",

	// root
	"math:root(25)":    builtin.Number{5},
	"math:root()":      "math:root: must have exactly 1 arg at file:9",
	"math:root(2++2)":  "missing term at file:12",
	"math:root({x=2})": "math:root: not a number at file:12",
}

func TestMath(t *testing.T) {
	r := &run.Runner{}
	s := math.Scope(builtin.Scope)

	for code, expected := range cases {
		t.Run(code, func(t *testing.T) {
			value := r.Eval(s, "file", code)
			if err, ok := value.(error); ok {
				if err.Error() != expected {
					t.Error("Failed", value)
				}
			} else if value != expected {
				t.Error("Failed", value)
			}
		})
	}
}
