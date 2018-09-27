# funnel

[![Status](https://travis-ci.com/funnelorg/funnel.svg?branch=master)](https://travis-ci.com/funnelorg/funnel?branch=master)
[![GoDoc](https://godoc.org/github.com/funnelorg/funnel?status.svg)](https://godoc.org/github.com/funnelorg/funnel)
[![codecov](https://codecov.io/gh/funnelorg/funnel/branch/master/graph/badge.svg)](https://codecov.io/gh/funnelorg/funnel)
[![GoReportCard](https://goreportcard.com/badge/github.com/funnelorg/funnel)](https://goreportcard.com/report/github.com/funnelorg/funnel)

Core features of the language:

   - Interpreted but plans to compile to JS or Go
   - Infix expressions
   - Pure functional, no assignments
   - Map expression to support complicated expression: `{x = 23, y = x+2}`
   - Define closures with `fun(arg1, arg2,.. expression)` syntax
   - data:list and data:map provide ability to do filter/map
   - general type system not yet implemented
   - easy to add custom functions
   - ability to import code

## Playground

See the language in action
[here](https://funnelorg.github.io/playground/)

## Introduction

Funnel is a an experimental simple functional language.

The language has been designed to have very little by way of syntax:
Everything is just an infix expression with a small set of builtin
binary operators: `+ - * / . = , < > <= >= != == & |` and two grouping
operators `()` and `{}`.

The curly braces functions as a "let expression":

```
    {
       x  = <some_expression>,
       y = x + x,
    }
```

Everything else is done via runtime functions.  Defining a function is
done via the `fun` helper:   `fun(x, y, x+y)`

The **dot notation** is used access the fields of a the `{}`
expression:

```
    {
       x = <some_expression>,
       y = x + x,
    }.y
```

The `data:list` function provides ability to work with array-like
entites:

```
    data:list(1, 2, 3).item(1)
```

Lists also support `filter` and `map` functions:

```
   data:list(1, 2, 3).filter(index < 2 || value > 2)
   data:list(1, 2, 3).map(value*2)
```

The `data:map` functions provides access to `filter` and `map` on top
of the regular map expression:

```
   data:map{x = 42}.filter(key == string x || value == 42)
```

Other functions include `code:import(url)` to fetch and execute code
at url.

## Vision

The main goal of the project is to play around with different type
systems and reactive approaches.  In particular, the short term goal
is to figure out how to integrate this with the
[project](https://github.com/dotchain/dot) (which provides data
synchronization support but also the general mutability mechanism)

