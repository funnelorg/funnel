// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package parse_test

import (
	"github.com/funnelorg/funnel/parse"
	"testing"
)

func TestSuccess(t *testing.T) {
	var cases = map[string]string{
		" ":                      "",
		"":                       "",
		"   x  ":                 "x",
		"x + y":                  "x + y",
		"x y":                    "x(y)",
		"x+(y-z)":                "x + (y - z)",
		"f(a, b, c)":             "f(a, b, c)",
		"x*y() + z()*r":          "x * y() + z() * r",
		"{}":                     "{}",
		"{x = f()+3,(y)=z}":      "{x = f() + 3, y = z}",
		"x+    y":                "x + y",
		"x()":                    "x()",
		"((x()).y)()":            "x().y()",
		"x().(y())":              "x().(y())",
		"'hello' + 'world'":      "hello + world",
		"\"hello\" + \"world\"":  "hello + world",
		"'hel\"lo' + 'world'":    "'hel\"lo' + world",
		"\"he'llo\" + \"world\"": "\"he'llo\" + world",
		"f(\"\")":                "f('')",
		"f {x = 42}":             "f({x = 42})",
		"f({x = 42})":            "f({x = 42})",
		"x {":                    "x({})",
		"x (":                    "x()",
		"x.y.z()":                "x.y.z()",
		"(x + y + z)*2":          "(x + y + z) * 2",
		"x + y + z)*2":           "(x + y + z) * 2",
		"x = z }.q":              "{x = z}.q",
		"{x=y,a=b}":              "{x = y, a = b}",
		"{x=y,a=b,m=n}":          "{x = y, a = b, m = n}",
	}

	for test, expected := range cases {
		t.Run(test, func(t *testing.T) {
			actual := parse.Parse("boo", test).String()
			if actual != expected {
				t.Error("Expected", expected, "got", actual)
			}
			if parse.Parse("boo", test).IsError() {
				t.Error("Unexpected error", test)
			}
		})
	}

}

func TestErrors(t *testing.T) {
	var cases = map[string]string{
		"x + + y":            "x + !('missing term', boo:4) + y",
		"x + ":               "x + !('missing term', boo:2)",
		"x + () + z":         "x + !('missing term', boo:4) + z",
		"x + f({x = y) + z":  "x + f(!('mismatched braces', boo:6)) + z",
		"x + f {x = (y} + z": "x + f({x = !('mismatched brackets', boo:11)}) + z",
		"{x}":                "!('invalid braces', boo:1)",
		"{x = y = z}":        "!('invalid equals use', boo:3)",
		"{x=y,z}":            "!('invalid key=value', boo:5)",
	}

	for test, expected := range cases {
		t.Run(test, func(t *testing.T) {
			actual := parse.Parse("boo", test).String()
			if actual != expected {
				t.Error("Expected", expected, "got", actual)
			}
		})
	}

}
