# funnel

[![Status](https://travis-ci.com/funnelorg/funnel.svg?branch=master)](https://travis-ci.com/funnelorg/funnel?branch=master)
[![GoDoc](https://godoc.org/github.com/funnelorg/funnel?status.svg)](https://godoc.org/github.com/funnelorg/funnel)
[![codecov](https://codecov.io/gh/funnelorg/funnel/branch/master/graph/badge.svg)](https://codecov.io/gh/funnelorg/funnel)
[![GoReportCard](https://goreportcard.com/badge/github.com/funnelorg/funnel)](https://goreportcard.com/report/github.com/funnelorg/funnel)

Funnel is a an experimental simple functional language.

The language has been designed to have very little by way of syntax:
Everything is just an infix expression with a small set of builtin
binary operators: `+ - * / . = ,` and two grouping operators `()` and
`{}`.

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

There are no built-in types of any sort -- all types and value
composition are done via functions.  For instance, `num(5)` converts
the parameter by interpreting its argument as an number.

## Vision

The main goal of the project is to play around with different type
systems and reactive approaches.  In particular, the short term goal
is to figure out how to integrate this with the
[https://github.com/dotchain/dot] project (which provides data
synchronization support but also the general mutability mechanism)

