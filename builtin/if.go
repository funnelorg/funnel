// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func iff(s run.Scope, args []parse.Node) interface{} {
	if len(args) <= 1 {
		return &run.ErrorStack{Message: "if: requires 2 or 3 args"}
	}
	r := &run.Runner{}
	switch cond := r.Run(s, args[0]).(type) {
	case bool:
		result := args[1]
		if !cond {
			if len(args) < 3 {
				return &run.ErrorStack{Message: "if: no else value"}
			}
			result = args[2]
		}
		return r.Run(s, result)
	case error:
		return cond
	}
	return &run.ErrorStack{Message: "if: not a valid condition"}
}
