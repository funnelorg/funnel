// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package builtin

import (
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func errorInternal(s run.Scope, args []parse.Node) interface{} {
	if len(args) > 0 && args[0].Token != nil {
		return &run.ErrorStack{Message: args[0].Token.S}
	}
	return &run.ErrorStack{Message: "unknown error"}
}

func errorf(s run.Scope, args []parse.Node) interface{} {
	r := &run.Runner{}
	if len(args) > 0 {
		switch t := r.Run(s, args[0]).(type) {
		case error:
			return &run.ErrorStack{Inner: t}
		case string:
			return &run.ErrorStack{Message: t}
		}
	}
	return &run.ErrorStack{Message: "unknown error"}
}
