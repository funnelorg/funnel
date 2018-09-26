// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package url_test

import (
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/url"
	"reflect"
	"testing"
)

var jsonURL = `https://api.myjson.com/bins/ngfbw`
var urlCall = `url(builtin:string "` + jsonURL + `")`
var s = url.Scope(builtin.Scope)
var data = map[interface{}]interface{}{"hello": "world"}

var cases = map[string]interface{}{
	"url()":                "url: requires exactly 1 arg at file:3",
	"url(1)":               "url: not a string at file:4",
	urlCall:                url.URL(jsonURL),
	"url(" + urlCall + ")": url.URL(jsonURL),
	urlCall + ".json()":    data,
	urlCall + ".boo":       "no such key at file:55",
	"url(1++2)":            "missing term at file:6",

	"url(builtin:string 'http://google.com').json()":      "invalid character '<' looking for beginning of value",
	"url(builtin:string 'http://google.com/boop').json()": "404 Not Found",
	`url(builtin:string "http:// /q").json()`:             `parse http:// /q: invalid character " " in host name`,
}

func TestUrl(t *testing.T) {
	for code, expected := range cases {
		t.Run(code, func(t *testing.T) {
			value := native(funnel.Eval(s, "file", code))
			if err, ok := value.(error); ok {
				if err.Error() != expected {
					t.Error("Failed", value)
				}
			} else if !reflect.DeepEqual(value, expected) {
				t.Error("Failed", value)
			}
		})
	}
}

func native(v interface{}) interface{} {
	s, ok := v.(scopeWithKeys)
	if !ok {
		return v
	}
	result := map[interface{}]interface{}{}
	s.ForEachKeys(func(key interface{}) bool {
		result[key] = native(s.Get(key))
		return false
	})
	return result
}

type scopeWithKeys interface {
	Get(key interface{}) interface{}
	ForEachKeys(fn func(interface{}) bool)
}
