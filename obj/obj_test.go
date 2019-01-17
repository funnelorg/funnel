// Copyright (C) 2019 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package obj_test

import (
	"fmt"
	"github.com/funnelorg/funnel/obj"
)

func ExampleNil() {
	caller := obj.S("caller")
	o := obj.Nil()

	fmt.Println(o.Invoke(caller).Get(caller, "stack"))
	fmt.Println(o.Get(caller, "boo").Get(caller, "stack"))

	// Output:
	// <nil> cannot be invoked
	// cannot access field boo of <nil>
}

func ExampleError() {
	caller := obj.S("caller")
	o := obj.Error(caller, "some error yo")

	fmt.Println("1:", o.Get(caller, "stack"))
	fmt.Println("2:", o.Get(caller, "boo").Get(caller, "stack"))
	fmt.Println("3:", o.Invoke(caller).Get(caller, "stack"))

	// Output:
	// 1: some error yo
	// 2: some error yo
	// cannot access field boo of <error>
	// 3: some error yo
	// some error yo
}

func ExampleS() {
	caller := obj.S("caller")
	o := obj.S("some string")

	fmt.Println("1:", o.Get(caller, "boo").Get(caller, "stack"))
	fmt.Println("2:", o.Invoke(caller).Get(caller, "stack"))

	// Output:
	// 1: cannot access field boo of <string>
	// 2: string is not callable
}

func ExampleTuple() {
	caller := obj.S("caller")
	o := obj.Tuple(map[string]obj.O{"hello": obj.S("world")})
	o2 := obj.Tuple(map[string]obj.O{"()": callable{}})

	fmt.Println("1:", o.Get(caller, "hello"))
	fmt.Println("2:", o2.Invoke(caller, obj.S("args")))
	fmt.Println("3:", o.Invoke(caller).Get(caller, "stack"))
	fmt.Println("4:", o.Get(caller, "boo").Get(caller, "stack"))

	// Output:
	// 1: world
	// 2: called with args
	// 3: tuple is not callable
	// 4: tuple does not have field boo
}

type callable struct{}

func (c callable) Invoke(caller obj.O, args ...obj.O) obj.O {
	return obj.S("called with " + string(args[0].(obj.S)))
}

func (c callable) Get(caller obj.O, name string) obj.O {
	panic("not expected")
}
