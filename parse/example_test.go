// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package parse_test

import (
	"fmt"
	"github.com/funnelorg/funnel/parse"
)

func Example_success() {
	fname := "code.fun"
	expression := "x + y {a = b}"
	parsed := parse.Parse(fname, expression)
	// print the infix formatted canonical version
	fmt.Println("Parsed:", parsed.String())

	// Output:
	// Parsed: x + y({a = b})
}

func Example_error() {
	fname := "code.fun"
	expression := "x ++ y"
	parsed := parse.Parse(fname, expression)
	// print the infix formatted canonical version
	fmt.Println("Parsed:", parsed.String())

	// Output:
	// Parsed: x + !('missing term', 'code.fun:3') + y
}
