// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build js,!jsreflect

package url

import (
	"github.com/funnelorg/funnel/data"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"honnef.co/go/js/xhr"
)

func (u URL) json(s run.Scope, args []parse.Node) interface{} {
	req := xhr.NewRequest("GET", string(u))
	req.Timeout = 1000 // one second, in milliseconds
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		return err
	}
	return data.Wrap(req.Response.Interface())
}
