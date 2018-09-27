// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// +build !js jsreflect

package url

import (
	"encoding/json"
	"github.com/funnelorg/funnel/run"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetch the URL contents
func (u URL) Fetch(fmt string) (interface{}, error) {
	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Get(string(u))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		_ = err
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, &run.ErrorStack{Message: resp.Status}
	}

	if fmt == Text {
		body, err := ioutil.ReadAll(resp.Body)
		return string(body), err
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
