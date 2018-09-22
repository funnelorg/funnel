// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package parse implements a simple expression parser.
package parse

import (
	"strconv"
	"unicode"
)

// Node is a node in the parse tree. It is a union of the three basic
// types:  Token (for all constants and variables), Call (for a
// call expressions) and Map (for a curly-braces expressions)
type Node struct {
	*Token
	*Call
	*Map
}

// String formats the expression to the canonical infix form
func (n Node) String() string {
	if n.Token != nil {
		return n.Token.String()
	}
	if n.Call != nil {
		return n.Call.String()
	}
	if n.Map != nil {
		return n.Map.String()
	}
	return ""
}

// IsError returns whether the current node is an "error"
func (n Node) IsError() bool {
	return n.Call != nil && n.Call.IsError()
}

// IsOperator returns if the current node is an operator
func (n Node) IsOperator() bool {
	return n.Call != nil && n.Call.IsOperator()
}

// wrap is almost the same as String except it adds a bracket around
// the result if the node is an operator of lower priority
func (n Node) wrap(pri int, isLeft bool) string {
	if n.IsOperator() && priority[n.Call.Nodes[0].Token.S] > pri {
		return n.String()
	}
	if !n.IsOperator() && pri < priority["."] {
		return n.String()
	}

	if n.Call == nil || (priority["."] == pri && isLeft) { // isLeft {
		return n.String()
	}

	return "(" + n.String() + ")"

}

// Loc is the location in the input where this node starts
func (n Node) Loc() *Loc {
	switch {
	case n.Token != nil:
		return n.Token.Loc
	case n.Call != nil:
		return n.Call.Loc
	}
	return n.Map.Loc
}

// Loc is a location in the input file being parsed. Offset is the
// offset in runes
type Loc struct {
	File   string
	Offset int
}

// String returns a human readable string
func (l Loc) String() string {
	return l.File + ":" + strconv.Itoa(l.Offset)
}

// Token is a token of input
type Token struct {
	*Loc
	S string
}

// String formats the token into a regular string
func (t Token) String() string {
	if t.S == "" {
		return "''"
	}

	quote, dquote := false, false
	for _, r := range t.S {
		_, isop := priority[string(r)]
		quote = quote || r == '\''
		dquote = dquote || (r == '"' || unicode.IsSpace(r) || isop)
	}
	if quote {
		return `"` + t.S + `"`
	}
	if dquote {
		return `'` + t.S + `'`
	}
	return t.S
}

// Call represents a function call or an operator expression
// The first node is the operator token for operator expressions and
// the function token/expression.  The rest of the nodes are the
// operands or arguments
type Call struct {
	*Loc
	Nodes []Node
}

// IsError checks if the call is to "error" which is how errors are
// represented
func (c Call) IsError() bool {
	return c.Nodes[0].Token != nil && c.Nodes[0].Token.S == "!"
}

// IsOperator checks if the Call is for an operator expression
func (c Call) IsOperator() bool {
	if c.Nodes[0].Token == nil {
		return false
	}
	_, ok := priority[c.Nodes[0].Token.S]
	return ok
}

// String formats the call expression properly
func (c Call) String() string {
	if c.IsError() {
		loc := Token{S: c.Loc.String()}
		return "!(" + c.Nodes[1].String() + ", " + loc.String() + ")"
	}

	if c.IsOperator() {
		return c.formatOperation()
	}

	result := c.Nodes[0].String() + "("
	for kk, nn := range c.Nodes[1:] {
		if kk == 0 {
			result += nn.String()
		} else {
			result += ", " + nn.String()
		}
	}
	return result + ")"
}

func (c Call) formatOperation() string {
	pri := priority[c.Nodes[0].Token.S]
	space := " "
	if c.Nodes[0].Token.S == "." {
		space = ""
	}
	result := c.Nodes[1].wrap(pri, true)
	for _, nn := range c.Nodes[2:] {
		result += space + c.Nodes[0].Token.S
		result += space + nn.wrap(pri, false)
	}
	return result
}

// Map represents a map structure.
type Map struct {
	*Loc
	Pairs []Pair
}

// String formats the map
func (m Map) String() string {
	result := "{"
	for kk, pp := range m.Pairs {
		if kk == 0 {
			result += pp.String()
		} else {
			result += ", " + pp.String()
		}
	}
	return result + "}"
}

// Pair represnets the key/value pair in a Map
type Pair struct {
	*Loc
	Key, Value Node
}

// String formats the key value pair
func (p Pair) String() string {
	return p.Key.String() + " = " + p.Value.String()
}
