// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package run

import "github.com/funnelorg/funnel/parse"

// Scope defines the interface for map-like objects.
type Scope interface {
	Get(key interface{}) interface{}
}

// Lazy defines a value that is "deferred"
type Lazy interface {
	Value() interface{}
}

type mapScope struct {
	parsed parse.Map
	base   Scope
	*Runner
	values  map[interface{}]interface{}
	keys    map[interface{}]int
	running map[interface{}]bool
}

func newMapScope(parsed parse.Map, base Scope, r *Runner) Scope {
	return &mapScope{parsed, base, r, nil, nil, map[interface{}]bool{}}
}

func (ms *mapScope) Value() interface{} {
	result := map[interface{}]interface{}{}
	for key := range ms.getKeys() {
		result[key] = ms.Get(key)
	}
	return result
}

func (ms *mapScope) Get(key interface{}) interface{} {
	v, ok := ms.values[key]
	if ok {
		return v
	}

	if ms.running[key] {
		loc := ms.parsed.Pairs[ms.keys[key]].Loc
		f, o := loc.File, loc.Offset
		return &ErrorStack{Message: "invalid recursion", File: f, Offset: o}
	}
	ms.running[key] = true
	defer func() {
		delete(ms.running, key)
	}()

	keys := ms.getKeys()
	idx, ok := keys[key]
	if !ok {
		return ms.base.Get(key)
	}
	if ms.values == nil {
		ms.values = map[interface{}]interface{}{}
	}

	val := ms.Runner.Run(ms, ms.parsed.Pairs[idx].Value)
	ms.values[key] = WrapError(val, ms.parsed.Pairs[idx].Loc)
	return ms.values[key]
}

func (ms *mapScope) getKeys() map[interface{}]int {
	if ms.keys != nil {
		return ms.keys
	}

	ms.keys = map[interface{}]int{}
	for idx, pp := range ms.parsed.Pairs {
		ms.keys[ms.Runner.run(ms.base, pp.Key)] = idx
	}

	return ms.keys
}

// NewScope creates a scope out of a set of maps, defaulting ot the
// base if no keys are found
func NewScope(maps []map[interface{}]interface{}, base Scope) Scope {
	combined := map[interface{}]interface{}{}
	for _, mm := range maps {
		if len(maps) == 1 {
			return &scope{mm, base}
		}
		for k, v := range mm {
			combined[k] = v
		}
	}
	return &scope{combined, base}
}

type scope struct {
	mm   map[interface{}]interface{}
	base Scope
}

func (s *scope) Get(key interface{}) interface{} {
	if v, ok := s.mm[key]; ok {
		return v
	}
	if s.base == nil {
		if s, ok := key.(string); ok {
			return &ErrorStack{Message: s + ": no such key"}
		}
		return &ErrorStack{Message: "No such key"}
	}
	return s.base.Get(key)
}
