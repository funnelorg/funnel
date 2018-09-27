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
		return m.filterExpr
	case "map":
		return m.mapExpr
	case "filterf":
		return run.ArgsResolver(m.filterf)
	case "mapf":
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

func (m mscope) filterExpr(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "filter: requires condition function"}
	}

	result := map[interface{}]interface{}{}
	r := &run.Runner{}
	var err error
	m.s.ForEachKeys(func(key interface{}) bool {
		v := m.s.Get(key)
		params := []map[interface{}]interface{}{{"key": key, "value": Wrap(v)}}
		switch check := r.Run(run.NewScope(params, s), args[0]).(type) {
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
	return Wrap(result)
}

func (m mscope) filterf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "filter: requires condition function"}
	}

	result := map[interface{}]interface{}{}
	var err error
	m.s.ForEachKeys(func(key interface{}) bool {
		v := m.s.Get(key)
		params := []interface{}{key, Wrap(v)}
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
	return Wrap(result)
}

func (m mscope) mapExpr(s run.Scope, args []parse.Node) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "map: requires 1 arg"}
	}

	result := map[interface{}]interface{}{}
	r := &run.Runner{}
	m.s.ForEachKeys(func(key interface{}) bool {
		v := m.s.Get(key)
		params := []map[interface{}]interface{}{{"key": key, "value": Wrap(v)}}
		result[key] = r.Run(run.NewScope(params, s), args[0])
		return false
	})
	return Wrap(result)
}

func (m mscope) mapf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "map: requires 1 arg"}
	}

	result := map[interface{}]interface{}{}
	m.s.ForEachKeys(func(key interface{}) bool {
		params := []interface{}{key, Wrap(m.s.Get(key))}
		result[key] = invoke("map", args[0], params)
		return false
	})
	return Wrap(result)
}
