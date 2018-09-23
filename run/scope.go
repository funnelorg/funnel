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
		return ErrInvalidRecursion
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

	ms.values[key] = ms.Runner.Run(ms, ms.parsed.Pairs[idx].Value)
	return ms.values[key]
}

func (ms *mapScope) getKeys() map[interface{}]int {
	if ms.keys != nil {
		return ms.keys
	}

	ms.keys = map[interface{}]int{}
	for idx, pp := range ms.parsed.Pairs {
		ms.keys[ms.Runner.runRaw(ms.base, pp.Key)] = idx
	}

	return ms.keys
}

func unwrapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case Lazy:
		return unwrapValue(v.Value())
	case map[interface{}]interface{}:
		result := map[interface{}]interface{}{}
		for key, value := range v {
			result[key] = unwrapValue(value)
		}
		return result
	}
	return v
}
