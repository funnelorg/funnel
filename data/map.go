// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package data

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/parse"
	"github.com/funnelorg/funnel/run"
)

func mapf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "map: needs 1 arg"}
	}
	s, ok := args[0].(scopeWithKeys)
	if !ok {
		return &run.ErrorStack{Message: "map: not a map"}
	}
	return mscope{s}
}

type scopeWithKeys interface {
	Get(key interface{}) interface{}
	ForEachKeys(fn func(interface{}) bool)
}

type mscope struct {
	s scopeWithKeys
}

func (m mscope) Get(key interface{}) interface{} {
	switch key {
	case "filter":
		return run.ArgsResolver(m.filterf)
	case "map":
		return run.ArgsResolver(m.mapf)
	case "count":
		return m.countf
	}
	return m.s.Get(key)
}

func (m mscope) ForEachKeys(fn func(interface{}) bool) {
	m.s.ForEachKeys(fn)
}

func (m mscope) countf(s run.Scope, args []parse.Node) interface{} {
	count := 0
	m.s.ForEachKeys(func(k interface{}) bool {
		count++
		return false
	})
	return builtin.Number{float64(count)}
}

func (m mscope) filterf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "filter: requires condition function"}
	}

	result := map[interface{}]interface{}{}
	var err error
	m.s.ForEachKeys(func(key interface{}) bool {
		v := m.s.Get(key)
		params := []interface{}{key, v}
		switch check := invoke("filter", args[0], params).(type) {
		case bool:
			if check {
				result[key] = v
			}
			return false
		case error:
			err = check
			return true
		}
		err = &run.ErrorStack{Message: "filter: function returned non-boolean"}
		return true
	})
	if err != nil {
		return err
	}

	filtered := run.NewScope([]map[interface{}]interface{}{result}, nil)
	return mscope{filtered.(scopeWithKeys)}
}

func (m mscope) mapf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "map: requires 1 arg"}
	}

	result := map[interface{}]interface{}{}
	m.s.ForEachKeys(func(key interface{}) bool {
		params := []interface{}{key, m.s.Get(key)}
		result[key] = invoke("map", args[0], params)
		return false
	})

	mapped := run.NewScope([]map[interface{}]interface{}{result}, nil)
	return mscope{mapped.(scopeWithKeys)}
}
