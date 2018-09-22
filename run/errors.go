// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run

// E simply wraps a string into an error
type E string

// Error returns the underlying string
func (e E) Error() string {
	return string(e)
}

// ErrUnknown is the generic error
var ErrUnknown = E("unknown error")

// ErrNotFunction is returned when a non-function is invoked
var ErrNotFunction = E("not a function")

// ErrInvalidRecursion is returned if expressions are recursively defined
var ErrInvalidRecursion = E("invalid recursion")
