// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"github.com/funnelorg/funnel/run"
	"testing"
)

func TestScope(t *testing.T) {
	m1 := map[interface{}]interface{}{
		"hello": "world",
		"one":   "one",
		11:      "ok",
	}

	m2 := map[interface{}]interface{}{
		"hello": "world2",
		"two":   "two",
	}

	s := run.NewScope([]map[interface{}]interface{}{m1, m2}, nil)

	cases := map[interface{}]interface{}{
		"hello": "world2",
		"one":   "one",
		"two":   "two",
		11:      "ok",
		44:      "No such key",
		"boo":   "boo: no such key",
	}

	check := func(v1, v2 interface{}) bool {
		if err, ok := v1.(error); ok {
			return err.Error() == v2
		}
		return v1 == v2
	}

	for key, expected := range cases {
		if x := s.Get(key); !check(x, expected) {
			t.Error("Unexpected result", key, "=>", x)
		}
	}

	m3 := map[interface{}]interface{}{"hello": "world3"}
	s = run.NewScope([]map[interface{}]interface{}{m3}, s)
	cases["hello"] = "world3"

	for key, expected := range cases {
		if x := s.Get(key); !check(x, expected) {
			t.Error("Unexpected result", key, "=>", x)
		}
	}
}
