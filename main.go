package main

import (
	"fmt"
	"gojs/ast"
	"gojs/intp"
	"gojs/lang"
	"gojs/parse"
	"gojs/tkn"
)

func jsPrint(values ...lang.Value) {
	str := "[JS] "
	for i, value := range values {
		f := "%s, "
		if i == len(values)-1 {
			f = "%s"
		}

		str += fmt.Sprintf(f, value)
	}
	fmt.Println(str)
}

func main() {
	program := `function foo(a, b) {
                  var c = a + b
				  return c
                }

				for (var i = 0; i < 3; i = i + 1) {
					print(i, foo(1, 2))
				}
`

	t := tkn.Tokenizer{}
	tokens := t.Tokenize(program)

	p := parse.NewParser(tokens)
	pp := p.Parse()

	d := &ast.Dumper{Indent: 4}
	d.DumpNode(&pp, 0)
	fmt.Println(d.Dump)

	//if n.Callee.Name == "print" {
	//	fmt.Printf("[JS] %s\n", args[0].String())
	//	return lang.NewUndefined()
	//}

	i := intp.NewInterpreter()
	i.BindNativeFunction("print", jsPrint)

	v := i.Do(&pp)
	fmt.Println(v.String())
}
