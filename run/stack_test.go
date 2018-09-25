// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run_test

import (
	"errors"
	"github.com/funnelorg/funnel/run"
	"testing"
)

func TestStack(t *testing.T) {
	inner := &run.ErrorStack{Message: "hello", File: "inner", Offset: 0}
	outer := &run.ErrorStack{Message: "outer", File: "outer", Offset: 10, Inner: inner}
	s := outer.Stack()
	expected := "hello at inner:0\nouter at outer:10\n"
	if s != expected {
		t.Error("unexpected stack", s)
	}

	err := errors.New("inner_error")
	outer = &run.ErrorStack{Message: "outer", File: "outer", Offset: 10, Inner: err}
	s = outer.Stack()
	expected = "inner_error\nouter at outer:10\n"
	if s != expected {
		t.Error("unexpected stack", s)
	}
}
