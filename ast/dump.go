package ast

import (
	"fmt"
	"strconv"
)

type Dumper struct {
	Dump   string
	Indent int
}

func (d *Dumper) DumpNode(node Node, level int) {
	switch n := node.(type) {
	case *Program:
		d.printIndent(level)
		d.append("Program[\n")
		for _, p := range n.Body {
			d.DumpNode(p, level+1)
		}
		d.printIndent(level)
		d.append("]\n")
	case *BlockStatement:
		d.printIndent(level)
		d.append("BlockStatement[\n")
		for _, p := range n.Body {
			d.DumpNode(p, level+1)
		}
		d.printIndent(level)
		d.append("]\n")
	case *FunctionDeclaration:
		d.printIndent(level)
		d.append("FunctionDeclaration[\n")
		d.DumpNode(n.Body, level+1)
		d.printIndent(level)
		d.append("]\n")
	case *VariableDeclaration:
		d.printIndent(level)
		d.append("VariableDeclaration[\n")
		for _, decl := range n.Declarations {
			d.DumpNode(decl, level+1)
		}
		d.printIndent(level)
		d.append("]\n")
	case *VariableDeclarator:
		d.printIndent(level)
		d.append("VariableDeclarator[\n")
		d.printIndent(level + 1)
		d.append(n.Id.Name)
		if n.Init != nil {
			d.append("=")
			d.DumpNode(n.Init, level+1)
		}
		d.printIndent(level)
		d.append("]\n")
	case *ReturnStatement:
		d.printIndent(level)
		d.append("ReturnStatement[\n")
		if n.Argument != nil {
			d.DumpNode(n.Argument, level+1)
		}
		d.printIndent(level)
		d.append("]\n")
	case *IfStatement:
		d.printIndent(level)
		d.append("IfStatement[\n")
		d.DumpNode(n.Test, level+1)
		d.DumpNode(n.Consequent, level+1)
		d.printIndent(level)
		d.append("]\n")
	case *ExpressionStatement:
		d.printIndent(level)
		d.append("ExpressionStatement[\n")
		d.DumpNode(n.Expression, level+1)
		d.printIndent(level)
		d.append("]\n")
	case *Identifier:
		d.printIndent(level)
		d.append("Identifier[")
		d.append(n.Name)
		d.append("]\n")
	case *CallExpression:
		d.printIndent(level)
		d.append("CallExpression[\n")
		d.printIndent(level + 1)
		d.append(n.Callee.Name)
		d.append("(\n")
		for i, arg := range n.Arguments {
			d.printIndent(level + 2)
			d.append(fmt.Sprintf("arg%d=", i))
			d.DumpNode(arg, level+2)
		}
		d.printIndent(level + 1)
		d.append(")\n")
		d.printIndent(level)
		d.append("]\n")
	case *BinaryExpression:
		d.printIndent(level)
		d.append("BinaryExpression[\n")
		d.printIndent(level + 1)
		d.append("lhs=")
		d.DumpNode(n.Left, level+1)
		d.printIndent(level + 1)
		d.append("op=(")
		d.append(n.Operator)
		d.append(")\n")
		d.printIndent(level + 1)
		d.append("rhs=")
		d.DumpNode(n.Right, level+1)
		d.printIndent(level)
		d.append("]\n")
	case *IntLiteral:
		d.printIndent(level)
		d.append(strconv.Itoa(n.Value))
		d.append("\n")
	case *AssignmentExpression:
		d.printIndent(level)
		d.append("AssignmentExpression[\n")
		d.printIndent(level + 1)
		d.append("lhs=")
		d.DumpNode(n.Left, level+1)
		d.printIndent(level + 1)
		d.append("op=(" + n.Operator + ")\n")
		d.printIndent(level + 1)
		d.append("rhs=")
		d.DumpNode(n.Right, level+1)
		d.printIndent(level)
		d.append("]\n")
	case *ForStatement:
		d.printIndent(level)
		d.append("ForStatement[\n")
		d.printIndent(level + 1)
		d.append("init=")
		d.DumpNode(n.Init, level+1)
		d.printIndent(level + 1)
		d.append("test=")
		d.DumpNode(n.Test, level+1)
		d.printIndent(level + 1)
		d.append("update=")
		d.DumpNode(n.Update, level+1)
		d.printIndent(level)
		d.append("]\n")
	}
}

func (d *Dumper) printIndent(level int) {
	if len(d.Dump) > 0 && d.Dump[len(d.Dump)-1] == '\n' {
		for i := 0; i < level; i++ {
			d.Dump += "│"
			for j := 0; j < d.Indent; j++ {
				d.Dump += " "
			}
		}
	}
}

func (d *Dumper) append(value string) {
	d.Dump += value
}
