// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package wiki implements some helpful wiki functions
package wiki

import (
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/data"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/url"
	"github.com/tvastar/htmltables"
	"strconv"
	"strings"
)

// Scope returns a new scope with the default wiki functions
func Scope(base run.Scope) run.Scope {
	return run.NewScope([]map[interface{}]interface{}{Map}, base)
}

// Map contains the functions provided by this package
var Map = map[interface{}]interface{}{"wiki:table": run.ArgsResolver(tablesf)}

func tablesf(args []interface{}) interface{} {
	if len(args) != 1 {
		return &run.ErrorStack{Message: "wiki:tables: requires exactly 1 arg"}
	}

	switch f := args[0].(type) {
	case string:
		return wikiTables(f)
	case error:
		return f
	}
	return &run.ErrorStack{Message: "wiki:tables: not a string"}
}

func wikiTables(name string) interface{} {
	prefix := "https://en.wikipedia.org/w/api.php?action=parse&format=json&page="
	u := url.URL(prefix + escape([]byte(name)) + "&origin=*")

	json, err := u.Fetch("json")
	if err != nil {
		return &run.ErrorStack{Message: "wiki:tables fetch error", Inner: err}
	}
	text, err := jsonPath(json, "parse", "text", "*")
	if err != nil {
		return err
	}

	tables, err := htmltables.Parse(text)
	if err == nil && len(tables) > 0 {
		return format(getTable(tables))
	}

	return &run.ErrorStack{Message: "wiki:tables parse error", Inner: err}
}

func jsonPath(json interface{}, keys ...string) (string, error) {
	if len(keys) == 0 {
		text, ok := json.(string)
		if !ok {
			return "", &run.ErrorStack{Message: "Parse failure: not a string"}
		}
		return text, nil
	}
	mm, ok := json.(map[string]interface{})
	if !ok {
		return "", &run.ErrorStack{Message: "Parse failure: " + keys[0]}
	}
	return jsonPath(mm[keys[0]], keys[1:]...)
}

func escape(b []byte) string {
	result := make([]byte, 0, len(b))
	for _, c := range b {
		if c == ' ' {
			result = append(result, '+')
		} else if shouldNotEncode(c) {
			result = append(result, c)
		} else {
			result = append(result, '%')
			result = append(result, "0123456789ABCDEF"[c>>4])
			result = append(result, "0123456789ABCDEF"[c&15])
		}
	}
	return string(result)
}

func shouldNotEncode(c byte) bool {
	return 'A' <= c && c <= 'Z' ||
		'a' <= c && c <= 'z' ||
		'0' <= c && c <= '9' ||
		c == '-' || c == '_' || c == '.' || c == '~'
}

func format(t *htmltables.Table) interface{} {
	result := []interface{}{}
	for _, row := range t.Rows {
		entry := map[string]interface{}{}
		for kk := range t.Headers {
			entry[t.Headers[kk]] = formatCell(row[kk])
		}
		result = append(result, entry)
	}
	return data.Wrap(result)
}

func formatCell(s string) interface{} {
	if x, err := strconv.ParseFloat(strings.Replace(s, ",", "", 10), 64); err == nil {
		return builtin.Number{x}
	}
	return s
}

func getTable(tables []*htmltables.Table) *htmltables.Table {
	for _, table := range tables {
		if strings.Contains(table.Attributes["class"], "wikitable") {
			return table
		}
	}
	return tables[0]
}
