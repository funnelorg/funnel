// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

// Error returns an error. If the first argument is a string, it is
// the message of the error.  The arguments are not evaluated.
func Error(s run.Scope, args []parse.Node) interface{} {
	if len(args) > 0 && args[0].Token != nil {
		return errors.New(args[0].Token.S)
	}
	return errors.New("error")
}
