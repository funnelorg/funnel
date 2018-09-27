// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build js,!jsreflect

package url

import "honnef.co/go/js/xhr"

// Fetch the URL contents
func (u URL) Fetch(fmt string) (interface{}, error) {
	req := xhr.NewRequest("GET", string(u))
	req.Timeout = 10000 // 10s
	req.ResponseType = fmt
	err := req.Send(nil)
	if err != nil {
		return nil, err
	}
	return req.Response.Interface()
}
