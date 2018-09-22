// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package parse implements the parser for the funnel language.
//
// The language is a simple expression (no statements):
//
// The standard binary operators +, -, *, / can be used in infix
// form.  No unary operators are supported.
//
// Function calls have the traditional f(a, b) syntax but also support
// the shortened syntax of f x (for single parameter case only)
//
// The map expression can be used to create a map of values:
//
//    {
//       x = something,
//       y = some other expressoin
//    }
//
// The fields of a map can be acessed with dot notation
//
// The syntax supports quotes and double quotes but these are not
// native types but instead escape sequences. That is, "hello world"
// is treated as an identifier whose name includes the space.
//
// The parser is quite loose -- automatically accepting unterminated
// brackets and braces (by auto-filling them). Missing terms (such as
// "x + + y") do cause an error. Errors are represented as a function
// call to the "!" function.  The previous expression would be parsed
// as if it were:
//
//      x + !("missing term", filename:3) + y
//
//
// Note that xyz:2 is a valid identifier. In fact, other than operator
// characters, quotes or space, everything is allowed as part of
// identifier names.  Even these restricted characters can appear if
// quotes are used: "+zero+" is a valid identifier.
package parse

import "unicode"

type parser struct {
	sh                       shunt
	fname, s                 string
	digit, id, quote, dquote int
}

func (p *parser) isOperator(r rune) bool {
	_, ok := priority[string([]rune{r})]
	return ok
}

func (p *parser) flush(idx int) {
	// if  there is a pending digit or id, flush it
	if p.digit >= 0 {
		loc := Loc{p.fname, p.digit}
		p.sh.Push(Token{&loc, p.s[p.digit:idx]})
		p.digit = -1
	}
	if p.id >= 0 {
		loc := Loc{p.fname, p.id}
		p.sh.Push(Token{&loc, p.s[p.id:idx]})
		p.id = -1
	}
}

func (p *parser) parse() Node {
	for kk, rr := range p.s {
		p.process(kk, rr)
	}

	p.flush(len(p.s))
	return p.sh.Parsed()
}

func (p *parser) process(kk int, rr rune) {
	loc := Loc{p.fname, kk}
	switch {
	case p.quote >= 0 && rr != '\'':
	case p.dquote >= 0 && rr != '"':
	case rr == '\'', rr == '"':
		p.onQuote(kk, rr)
	case unicode.IsSpace(rr):
		p.flush(kk)
	case p.isOperator(rr):
		p.flush(kk)
		p.sh.Push(Token{&loc, string(rr)})
	case unicode.IsDigit(rr):
		if p.digit == -1 && p.id == -1 {
			p.digit = kk
		}
	default:
		if p.digit == -1 && p.id == -1 {
			p.id = kk
		}
	}
}

func (p *parser) onQuote(kk int, rr rune) {
	if rr == '\'' {
		if p.quote < 0 {
			p.flush(kk)
			p.quote = kk
		} else {
			l := Loc{p.fname, p.quote}
			p.sh.Push(Token{&l, p.s[p.quote+1 : kk]})
			p.quote = -1
		}
		return
	}

	if p.dquote < 0 {
		p.flush(kk)
		p.dquote = kk
	} else {
		l := Loc{p.fname, p.dquote}
		p.sh.Push(Token{&l, p.s[p.dquote+1 : kk]})
		p.dquote = -1
	}
}

// Parse returns the parsed form of the provided string
func Parse(fname, s string) Node {
	p := &parser{shunt{isCall: map[*Loc]bool{}}, fname, s, -1, -1, -1, -1}
	return p.parse()
}
