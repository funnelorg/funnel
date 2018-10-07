// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/funnelorg/funnel"
	"github.com/funnelorg/funnel/builtin"
	"github.com/funnelorg/funnel/code"
	"github.com/funnelorg/funnel/data"
	"github.com/funnelorg/funnel/math"
	"github.com/funnelorg/funnel/run"
	"github.com/funnelorg/funnel/url"
	"github.com/funnelorg/funnel/wiki"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	expr := flag.String("eval", "", "evaluate expression")
	file := flag.String("run", "", "evaluate file")
	help := flag.Bool("help", false, "help")

	flag.Parse()

	if *help {
		log.Fatalf("%s [-run <file> | -eval <expr>]", os.Args[0])
	}

	scope := wiki.Scope(code.Scope(data.Scope(url.Scope(math.Scope(builtin.Scope)))))
	if *expr != "" {
		fmt.Println(funnel.Eval(scope, "funnnel:eval", *expr))
		return
	}

	if *file != "" {
		code, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(funnel.Eval(scope, *file, string(code)))
		return
	}

	values := []map[interface{}]interface{}{{}}
	idx := 1
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("$%d > ", idx)
		code, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		result := funnel.Eval(run.NewScope(values, scope), "repl", code)
		fmt.Println(result)
		values[0][fmt.Sprintf("$%d", idx)] = result
	}
}
