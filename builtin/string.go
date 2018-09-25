// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func stringf(s run.Scope, args []parse.Node) interface{} {
	if len(args) > 0 && args[0].Token != nil {
		return args[0].Token.S
	}
	return &run.ErrorStack{Message: "unknown error"}
}
