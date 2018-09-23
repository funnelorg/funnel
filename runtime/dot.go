// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package runtime

import (
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

// Dot implements the object dot notation accessing the sequence of fields.
func Dot(s run.Scope, args []parse.Node) interface{} {
	r := &run.Runner{}
	result := r.LazyRun(s, args[0])
	for _, next := range args[1:] {
		var key interface{}
		if next.Token != nil {
			key = next.Token.S
		} else {
			key = r.Run(s, next)
		}

		if err, ok := key.(error); ok {
			return err
		}

		switch sx := result.(type) {
		case run.Scope:
			result = sx.Get(key)
		case error:
			return sx
		default:
			return errors.New("cannot use dot with non-map")
		}
	}
	return result
}
