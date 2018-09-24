// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

// String converts its first token into a string
func String(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return errors.New("string: incorrect number of args")
	}

	if args[0].Token != nil {
		return args[0].Token.S
	}

	result := (&run.Runner{}).Run(s, args[0])
	switch result.(type) {
	case string, error:
		return result
	}
	return errors.New("string: invalid arg type")
}
