// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin_test

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/run"
	"testing"
)

var cases = map[string]interface{}{
	// builtin:number
	"builtin:number()":    "builtin:number: must have 1 arg at file:14",
	"builtin:number(1+2)": "builtin:number: must be a number at file:16",
	"2a":  parseFloatError("2a"),
	"5.2": builtin.Number{5.2},

	// add
	"3 + 2":     builtin.Number{5},
	"3 + {x=2}": "not a number at file:6",
	"3 ++ 2":    "missing term at file:3",

	// sub
	"3 - 2":     builtin.Number{1},
	"-2":        builtin.Number{-2},
	"3 - 2 -2":  builtin.Number{-1},
	"3 - {x=2}": "not a number at file:6",
	"3 -- 2":    "missing term at file:3",

	// mult
	"3*2*3":     builtin.Number{18},
	"3 * {x=2}": "not a number at file:6",
	"3 ** 2":    "missing term at file:3",

	// div
	"12/2/3":    builtin.Number{2},
	"3 / {x=2}": "not a number at file:6",
	"3 // 2":    "missing term at file:3",

	// more
	"2 > 1":         true,
	"2 > 1 > 0":     true,
	"2 > 500 > > 3": false,
	"2 >> 500":      "missing term at file:3",
	"2 > {x=2}":     "not a number at file:6",

	// other comparisons
	"2 < 1":  false,
	"2 >= 2": true,
	"2 >= 3": false,
	"2 <= 3": true,
	"2 == 2": true,
	"2 != 3": true,

	// logical
	"2 < 1 & 4":     false,
	"2 > 1 & 4":     "not a condition at file:8",
	"2 > 1 & 2a":    parseFloatError("2a"),
	"2 > 1 & 2 > 1": true,
	"2 < 1 | 4":     "not a condition at file:8",
	"2 > 1 | 4":     true,
	"2 < 1 | 2a":    parseFloatError("2a"),
	"2 < 1 | 2 < 1": false,

	// if
	"if()":          "if: requires 2 or 3 args at file:2",
	"if(1,2,3)":     "if: not a valid condition at file:4",
	"if(2a,2,3)":    parseFloatError("2a"),
	"if(1 > 2,2,3)": builtin.Number{3},
	"if(1 > 2,2)":   "if: no else value at file:8",
	"if(1 < 2,2)":   builtin.Number{2},

	// error
	"!(1)":                         "unknown error at file:2",
	"error()":                      "unknown error at file:5",
	"error(builtin:string(hello))": "hello at file:21",
	"error(2++1)":                  "missing term at file:8",
	"builtin:string(1)":            "unknown error at file:15",

	// dot
	"{x = 2}.x":                   builtin.Number{2},
	"{x = 2}.(builtin:string(x))": builtin.Number{2},
	"{x = 2}.(2a)":                parseFloatError("2a"),
	"{x = 2}.y":                   "y: no such key at file:7",
	"{x = 2}.x.y":                 "cannot use dot with non-map at file:9",
	"(1 ++ 2).x":                  "missing term at file:4",

	// fun
	"fun()()":        nil,
	"fun(1+1)()":     builtin.Number{2},
	"fun(x, x+2)(2)": builtin.Number{4},
	"fun(x+2,x)":     "fun: invalid param name at file:5",

	// other
	"3*2 + 2*5": builtin.Number{16},
	"3 - 5*2":   builtin.Number{-7},
	"12/3/2":    builtin.Number{2},
}

func TestBuiltin(t *testing.T) {
	r := &run.Runner{}
	s := builtin.Scope

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

func parseFloatError(s string) string {
	return `strconv.ParseFloat: parsing "` + s + `": invalid syntax`
}
