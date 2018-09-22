// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package funnel_test

import (
	"github.com/funnelorg/funnel"
	"testing"
)

func TestSimple(t *testing.T) {
	result := funnel.Eval(nil, "boo", "int(4)+int(2)")
	if result != 6 {
		t.Error("unexpected result", result)
	}
}
