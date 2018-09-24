// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package url

import (
	"encoding/json"
	"errors"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
	"net/http"
	"time"
)

func (u URL) json(s run.Scope, args []parse.Node) interface{} {
	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Get(string(u))
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		_ = err
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}
	return result
}
