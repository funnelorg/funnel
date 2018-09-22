// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package parse

var priority = map[string]int{
	"{": 0,
	"}": 0,
	"(": 0,
	")": 0,
	",": 1,
	"=": 2,
	"+": 10,
	"-": 10,
	"*": 20,
	"/": 20,
	".": 30,
}

var isAssoc = map[string]bool{
	",": true,
	"=": true,
	"+": true,
	"-": true,
	"*": true,
	"/": true,
}

type shunt struct {
	terms  []Node
	ops    []Token
	last   *string
	isCall map[*Loc]bool
}

func (s *shunt) lastWasTerm() bool {
	if s.last == nil {
		return false
	}

	_, ok := priority[*s.last]
	return !ok || *s.last == ")" || *s.last == "}"
}

func (s *shunt) Parsed() Node {
	if s.last == nil {
		return Node{}
	}

	l := len(s.ops)
	missingTerm := !s.lastWasTerm()
	for kk := l - 1; kk >= 0; kk-- {
		op := s.ops[kk]
		if op.S == "{" {
			s.Push(Token{op.Loc, "}"})
			missingTerm = false
		} else if op.S == "(" {
			s.Push(Token{op.Loc, ")"})
			missingTerm = false
		}
	}

	if len(s.ops) > 0 && missingTerm {
		op, _ := s.topOp()
		s.terms = append(s.terms, s.error("missing term", op.Loc))
	}

	for op, ok := s.popOp(); ok; op, ok = s.popOp() {
		s.mergeLastTwoTerms(op)
	}

	return s.terms[0]
}

func (s *shunt) Push(t Token) {
	switch t.S {
	case "(", "{":
		if s.lastWasTerm() {
			s.isCall[t.Loc] = true
		}

		s.ops = append(s.ops, t)
	case ")":
		s.pushEndBrackets(t)
	case "}":
		s.pushEndBraces(t)
	default:
		if pri, isOp := priority[t.S]; isOp {
			s.pushOp(t, pri)
		} else {
			s.pushTerm(t)
		}
	}
	s.last = &t.S
}

func (s *shunt) error(message string, errorLoc *Loc) Node {
	nodes := []Node{
		{Token: &Token{errorLoc, "!"}},
		{Token: &Token{errorLoc, message}},
	}
	return Node{Call: &Call{Loc: errorLoc, Nodes: nodes}}
}

func (s *shunt) popTerm(errorLoc *Loc) Node {
	l := len(s.terms)
	last := s.terms[l-1]
	s.terms = s.terms[:l-1]
	return last
}

func (s *shunt) popOp() (Token, bool) {
	if l := len(s.ops); l > 0 {
		last := s.ops[l-1]
		s.ops = s.ops[:l-1]
		return last, true
	}
	return Token{}, false
}

func (s *shunt) makeCall(op Token, arg1, arg2 Node) Node {
	nn := []Node{{Token: &op}, arg1, arg2}
	return Node{Call: &Call{Loc: op.Loc, Nodes: nn}}
}

func (s *shunt) isAssoc(n Node, op string) bool {
	if !isAssoc[op] || n.Call == nil || n.Call.Nodes[0].Token == nil {
		return false
	}
	return n.Call.Nodes[0].Token.S == op
}

func (s *shunt) pushOp(t Token, pri int) {
	if !s.lastWasTerm() {
		// two consecutive ops.  add a missing term error
		s.terms = append(s.terms, s.error("missing term", t.Loc))
	}
	for len(s.ops) > 0 && priority[s.ops[len(s.ops)-1].S] >= pri {
		last, _ := s.popOp()
		s.mergeLastTwoTerms(last)
	}
	s.ops = append(s.ops, t)
}

func (s *shunt) pushTerm(t Token) {
	term := Node{Token: &t}
	if s.lastWasTerm() {
		nn := []Node{s.popTerm(t.Loc), term}
		s.terms = append(s.terms, Node{Call: &Call{Loc: t.Loc, Nodes: nn}})
	} else {
		s.terms = append(s.terms, term)
	}
}

func (s *shunt) mergeLastTwoTerms(op Token) {
	arg2 := s.popTerm(op.Loc)
	arg1 := s.popTerm(op.Loc)
	if s.isAssoc(arg1, op.S) {
		arg1.Call.Nodes = append(arg1.Call.Nodes, arg2)
		s.terms = append(s.terms, arg1)
	} else {
		s.terms = append(s.terms, s.makeCall(op, arg1, arg2))
	}
}

func (s *shunt) topOp() (Token, bool) {
	if l := len(s.ops); l > 0 {
		return s.ops[l-1], true
	}
	return Token{}, false
}

func (s *shunt) rollupDots() {
	top, ok := s.topOp()
	for ok && top.S == "." {
		s.popOp()
		s.mergeLastTwoTerms(top)
		top, ok = s.topOp()
	}
}

func (s *shunt) pushEndBrackets(op Token) {
	top, ok := s.topOp()
	if s.last != nil && *s.last == "(" {
		s.popOp()
		if s.isCall[top.Loc] {
			s.rollupDots()
			term := s.popTerm(op.Loc)
			nn := Node{Call: &Call{top.Loc, []Node{term}}}
			s.terms = append(s.terms, nn)
		} else {
			nn := s.error("missing term", top.Loc)
			s.terms = append(s.terms, nn)
		}
		return
	}

	for ok && top.S != "(" {
		if top.S == "{" {
			// fake an "}" and error out
			s.pushEndBraces(Token{op.Loc, "}"})
			s.popTerm(op.Loc)
			s.terms = append(s.terms, s.error("mismatched braces", top.Loc))
		} else {
			s.popOp()
			s.mergeLastTwoTerms(top)
		}
		top, ok = s.topOp()
	}
	if !ok {
		return
	}
	s.popOp()
	if s.isCall[top.Loc] {
		term := s.popTerm(top.Loc)
		s.rollupDots()
		fn := s.popTerm(top.Loc)
		nn := Node{Call: &Call{term.Loc(), []Node{fn, term}}}
		if s.isAssoc(term, ",") {
			nn.Call.Nodes = append([]Node{fn}, term.Call.Nodes[1:]...)
		}
		s.terms = append(s.terms, nn)
	}
}

func (s *shunt) pushEndBraces(op Token) {
	top, ok := s.topOp()
	if s.last != nil && *s.last == "{" {
		s.popOp()
		mm := Node{Map: &Map{Loc: top.Loc}}
		if s.isCall[top.Loc] {
			term := s.popTerm(top.Loc)
			nn := Node{Call: &Call{top.Loc, []Node{term, mm}}}
			s.terms = append(s.terms, nn)
		} else {
			s.terms = append(s.terms, mm)
		}
		return
	}

	for ok && top.S != "{" {
		if top.S == "(" {
			// fake an ")" and error out
			s.pushEndBrackets(Token{op.Loc, ")"})
			s.popTerm(op.Loc)
			s.terms = append(s.terms, s.error("mismatched brackets", top.Loc))
		} else {
			s.popOp()
			s.mergeLastTwoTerms(top)
		}
		top, ok = s.topOp()
	}
	s.popOp()
	mm := s.makeMap(s.popTerm(op.Loc))
	if s.isCall[top.Loc] {
		fn := s.popTerm(op.Loc)
		nn := Node{Call: &Call{top.Loc, []Node{fn, mm}}}
		s.terms = append(s.terms, nn)
	} else {
		s.terms = append(s.terms, mm)
	}
}

func (s *shunt) makeMap(n Node) Node {
	if s.isAssoc(n, "=") {
		if len(n.Call.Nodes) != 3 {
			return s.error("invalid equals use", n.Call.Loc)
		}
		key, value := n.Call.Nodes[1], n.Call.Nodes[2]
		return Node{Map: &Map{n.Loc(), []Pair{{n.Call.Loc, key, value}}}}
	}

	if !s.isAssoc(n, ",") {
		return s.error("invalid braces", n.Loc())
	}

	pairs := []Pair{}
	for _, node := range n.Call.Nodes[1:] {
		if !s.isAssoc(node, "=") || len(node.Call.Nodes) != 3 {
			return s.error("invalid key=value", node.Loc())
		}
		key, value := node.Call.Nodes[1], node.Call.Nodes[2]
		pairs = append(pairs, Pair{node.Call.Loc, key, value})
	}
	return Node{Map: &Map{n.Loc(), pairs}}
}
