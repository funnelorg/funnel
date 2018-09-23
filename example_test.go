// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package funnel

import (
	"fmt"
	"github.com/funnelorg/funnel/runtime"
)

func ExampleEval_customScope() {
	m := map[interface{}]interface{}{
		"x": 5.0,
		"square": func(n runtime.Number) float64 {
			return n.F * n.F
		},
	}
	myScope := runtime.NewScope(m, runtime.DefaultScope)

	x := Eval(myScope, "myfile.go", "square(x + 2)")
	fmt.Println("square:", x)

	// Output:
	// square: 49
}

func ExampleEval_defaultScope() {
	fmt.Println(Eval(nil, "myfile.go", "5 + 2"))

	// Output:
	// {7}
}
