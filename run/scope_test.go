// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"testing"
)

func TestForEachKeys(t *testing.T) {
	r := &run.Runner{}
	s := &defaultScope{r}

	mscope := r.LazyRun(s, parse.Parse("code", "{x = 4, y = 2}"))
	found := false
	mscope.(keys).ForEachKeys(func(key interface{}) bool {
		if key != "x" && key != "y" {
			t.Fatal("Unexpected key", key)
		}
		found = true
		return true
	})
	if !found {
		t.Fatal("No keys found")
	}

	m1 := map[interface{}]interface{}{
		"hello": "world",
		"one":   "one",
		11:      "ok",
	}

	m2 := map[interface{}]interface{}{
		"hello": "world2",
		"two":   "two",
	}

	found = false
	sx := run.NewScope([]map[interface{}]interface{}{m1, m2}, nil)
	sx.(keys).ForEachKeys(func(key interface{}) bool {
		if m1[key] == nil && m2[key] == nil {
			t.Fatal("Unexpected key", key)
		}
		found = true
		return true
	})
	if !found {
		t.Fatal("No keys found")
	}
}

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

type keys interface {
	ForEachKeys(fn func(interface{}) bool)
}
