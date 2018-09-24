// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package url_test

import (
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/runtime"
	"github.com/funnelorg/funnel/url"
	"reflect"
	"testing"
)

var jsonURL = `https://jsonplaceholder.typicode.com/todos/1`

func TestUrl(t *testing.T) {
	code := `url("` + jsonURL + `")`
	x := funnel.Eval(url.New(runtime.DefaultScope), "boo", code)
	if !reflect.DeepEqual(x, url.URL(jsonURL)) {
		t.Error("unexpected result", x, reflect.TypeOf(x))
	}

	code = `url("` + jsonURL + `").json()`
	x = funnel.Eval(url.New(runtime.DefaultScope), "boo", code)
	expected := map[string]interface{}{
		"userId":    1.0,
		"id":        1.0,
		"title":     "delectus aut autem",
		"completed": false,
	}

	if !reflect.DeepEqual(x, expected) {
		t.Error("unexpected response", x)
	}
}

func TestUrlErrors(t *testing.T) {
	s := url.New(runtime.DefaultScope)

	x := funnel.Eval(s, "boo", "url()")
	if err, ok := x.(error); !ok || err.Error() != `url: requries exactly 1 arg` {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(builtin:string(a), builtin:string(b))")
	if err, ok := x.(error); !ok || err.Error() != `url: requries exactly 1 arg` {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(url())")
	if err, ok := x.(error); !ok || err.Error() != `url: requries exactly 1 arg` {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(builtin:string(a))")
	if x != url.URL("a") {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(url(builtin:string(a)))")
	if x != url.URL("a") {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(1)")
	if err, ok := x.(error); !ok || err.Error() != `url: not a string` {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", "url(a).boo")
	if err, ok := x.(error); !ok || err.Error() != `no such key` {
		t.Error("unexpected result", x)
	}
}

func TestUrlJSONErrors(t *testing.T) {
	s := url.New(runtime.DefaultScope)

	x := funnel.Eval(s, "boo", `url("http://   /q").json()`)
	if err, ok := x.(error); !ok || err.Error() != `parse http://   /q: invalid character " " in host name` {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", `url("http://google.com/boop").json()`)
	if err, ok := x.(error); !ok || err.Error() != "404 Not Found" {
		t.Error("unexpected result", x)
	}

	x = funnel.Eval(s, "boo", `url("http://google.com").json()`)
	if err, ok := x.(error); !ok || err.Error() != `invalid character '<' looking for beginning of value` {
		t.Error("unexpected result", x)
	}
}
