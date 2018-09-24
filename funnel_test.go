// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package funnel_test

import (
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/runtime"
	"testing"
)

func TestSimple(t *testing.T) {
	result := funnel.Eval(nil, "boo", "4+2")
	expected := runtime.Number{6.0}
	if result != expected {
		t.Error("unexpected result", result)
	}
}
