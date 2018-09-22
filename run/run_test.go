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
		"int(4)":                           4,
		"int(int(4))":                      4,
		"int(4) + int(5)":                  9,
		"int(4) + int(5) + int(6)":         15,
		"{x = int(4)}":                     map[interface{}]interface{}{"x": 4},
		"{x = int(4), y = x+int(2)}.y":     6,
		"{x = int(2), y = x + x, z = 3}.y": 4,
	}

	r := &run.Runner{}
	for name, expected := range cases {
		t.Run(name, func(t *testing.T) {
			actual := r.Eval(nil, "code", name)
			if !reflect.DeepEqual(actual, expected) {
				t.Error("Expected", expected, "got", actual)
			}
		})
	}
}

func TestError(t *testing.T) {
	cases := map[string]interface{}{
		"boo(5)":               "unknown identifier: boo",
		"2.3":                  "unknown identifier: 2",
		"int('2.2')":           `strconv.Atoi: parsing "2.2": invalid syntax`,
		"int(5)(4)":            "not a function",
		"{x = 5}.(int(5))":     "unknown identifier",
		"int(2) ++ int(3)":     "missing term",
		"!(int(5))":            "unknown error",
		"!(! boo)":             "boo",
		"int(5).x":             "cannot use dot with non-map",
		"int(5, 3)":            "int expects a single arg",
		"int({x = 2})":         "int: unknown argument type",
		"int(5) + {x = 2}":     "sum: only works with ints",
		"{x = y, y = x + 2}.x": "invalid recursion",
	}

	r := &run.Runner{}
	for name, expected := range cases {
		t.Run(name, func(t *testing.T) {
			actual := r.Eval(nil, "code", name)
			if err, ok := actual.(error); !ok {
				t.Error("Unexpected successful result", actual)
			} else if err.Error() != expected {
				t.Error("Expected", expected, "got", err)
			}
		})
	}
}
