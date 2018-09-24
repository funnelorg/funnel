// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"github.com/funnelorg/funnel/run"
	"reflect"
	"testing"
)

func TestSuccess(t *testing.T) {
	cases := map[string]interface{}{
		"builtin:number(4)":           4,
		"4":                           4,
		"4 + 5":                       9,
		"4 + 5 + 6":                   15,
		"{x = 4}":                     map[interface{}]interface{}{"x": 4},
		"{x = 4, y = x+2}.y":          6,
		"{x = 2, y = x + x, z = 3}.y": 4,
		"{builtin:string x = 5}.x":    5,
	}

	r := &run.Runner{}
	s := &defaultScope{r}
	for name, expected := range cases {
		t.Run(name, func(t *testing.T) {
			actual := r.Eval(s, "code", name)
			if !reflect.DeepEqual(actual, expected) {
				t.Error("Expected", expected, "got", actual)
			}
		})
	}
}

func TestError(t *testing.T) {
	cases := map[string]interface{}{
		"boo(5)": "unknown identifier: boo",
		"2.3":    `strconv.Atoi: parsing "2.3": invalid syntax`,
		"builtin:number('2.2')":   `strconv.Atoi: parsing "2.2": invalid syntax`,
		"5(4)":                    "not a function",
		"{x = 5}.(5)":             "unknown identifier",
		"2 ++ 3":                  "missing term",
		"!(5)":                    "unknown error",
		"!(! boo)":                "boo",
		"(5).x":                   "cannot use dot with non-map",
		"builtin:number(5, 3)":    "int expects a single arg",
		"builtin:number({x = 2})": "int: unknown argument type",
		"5 + {x = 2}":             "sum: only works with ints",
		"{x = y, y = x + 2}.x":    "invalid recursion",
	}

	r := &run.Runner{}
	s := &defaultScope{r}
	for name, expected := range cases {
		t.Run(name, func(t *testing.T) {
			actual := r.Eval(s, "code", name)
			if err, ok := actual.(error); !ok {
				t.Error("Unexpected successful result", actual)
			} else if err.Error() != expected {
				t.Error("Expected", expected, "got", err)
			}
		})
	}
}
