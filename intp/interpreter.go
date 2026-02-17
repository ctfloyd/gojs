package intp

import (
	"gojs/ast"
	"gojs/lang"
)

type scope struct {
	vars map[string]lang.Value
}

func newScope() scope {
	return scope{make(map[string]lang.Value)}
}

func (s scope) Put(name string, value lang.Value) {
	s.vars[name] = value
}

func (s scope) Get(name string) (lang.Value, bool) {
	v, ok := s.vars[name]
	return v, ok
}

type Interpreter struct {
	scope []scope
}

func NewInterpreter() *Interpreter {
	return &Interpreter{scope: []scope{newScope()}}
}

func (i *Interpreter) BindNativeFunction(name string, f func(values ...lang.Value)) {
	binder := lang.NewObj(&lang.NativeFunction{Function: f})
	i.put(name, binder)
}

func (i *Interpreter) get(name string) lang.Value {
	for idx := len(i.scope) - 1; idx >= 0; idx-- {
		s := i.scope[idx]
		if v, ok := s.Get(name); ok {
			return v
		}
	}
	panic("var: " + name + " not in scope")
	return lang.Value{}
}

func (i *Interpreter) put(name string, value lang.Value) {
	i.latestScope().Put(name, value)
}

func (i *Interpreter) latestScope() scope {
	return i.scope[len(i.scope)-1]
}

func (i *Interpreter) enterScope() {
	i.scope = append(i.scope, newScope())
}

func (i *Interpreter) exitScope() {
	i.scope = i.scope[:len(i.scope)-1]
}

func (i *Interpreter) Do(n ast.Node) lang.Value {
	switch n := n.(type) {
	case *ast.AssignmentExpression:
		return i.assignmentExpression(n)
	case *ast.BinaryExpression:
		return i.binaryExpression(n)
	case *ast.BlockStatement:
		return i.blockStatement(n)
	case *ast.CallExpression:
		return i.callExpression(n)
	case *ast.ExpressionStatement:
		return i.expressionStatement(n)
	case *ast.ForStatement:
		return i.forStatement(n)
	case *ast.FunctionDeclaration:
		return i.functionDeclaration(n)
	case *ast.Identifier:
		return i.identifier(n)
	case *ast.IfStatement:
		return i.ifStatement(n)
	case *ast.IntLiteral:
		return i.intLiteral(n)
	case *ast.Program:
		return i.program(n)
	case *ast.ReturnStatement:
		return i.returnStatement(n)
	case *ast.VariableDeclarator:
		return i.variableDeclarator(n)
	case *ast.VariableDeclaration:
		return i.variableDeclaration(n)
	default:
		panic("unsupported node")
	}
	return lang.Value{}
}

func (i *Interpreter) blockStatement(n *ast.BlockStatement) lang.Value {
	lv := lang.Value{}
	for _, n1 := range n.Body {
		lv = i.Do(n1)
	}
	return lv
}

func (i *Interpreter) program(n *ast.Program) lang.Value {
	lv := lang.Value{}
	for _, n1 := range n.Body {
		lv = i.Do(n1)
	}
	return lv
}

func (i *Interpreter) identifier(n *ast.Identifier) lang.Value {
	return i.get(n.Name)
}

func (i *Interpreter) ifStatement(n *ast.IfStatement) lang.Value {
	test := i.Do(n.Test)

	v := lang.NewUndefined()
	if test.Bool {
		i.enterScope()
		v = i.Do(n.Consequent)
		i.exitScope()
	}
	return v
}

func (i *Interpreter) variableDeclarator(n *ast.VariableDeclarator) lang.Value {
	init := i.Do(n.Init)
	i.put(n.Id.Name, init)
	return init
}

func (i *Interpreter) assignmentExpression(n *ast.AssignmentExpression) lang.Value {
	identifier := n.Left.(*ast.Identifier)
	update := i.Do(n.Right)
	i.put(identifier.Name, update)
	return update
}

func (i *Interpreter) binaryExpression(n *ast.BinaryExpression) lang.Value {
	l := i.Do(n.Left)
	r := i.Do(n.Right)
	if n.Operator == "+" {
		return lang.NewInt(l.Int + r.Int)
	} else if n.Operator == ">" {
		return lang.NewBool(l.Int > r.Int)
	} else if n.Operator == "<" {
		return lang.NewBool(l.Int < r.Int)
	} else {
		panic("unsupported operation")
	}
}

func (i *Interpreter) functionDeclaration(n *ast.FunctionDeclaration) lang.Value {
	f := lang.NewObj(&lang.Function{Name: n.Id.Name, Body: n.Body, Parameters: n.Parameters})
	i.put(n.Id.Name, f)
	return f
}

func (i *Interpreter) returnStatement(n *ast.ReturnStatement) lang.Value {
	return i.Do(n.Argument)
}

func (i *Interpreter) variableDeclaration(n *ast.VariableDeclaration) lang.Value {
	lv := lang.Value{}
	for _, d := range n.Declarations {
		lv = i.Do(d)
	}
	return lv
}

func (i *Interpreter) expressionStatement(n *ast.ExpressionStatement) lang.Value {
	return i.Do(n.Expression)
}

func (i *Interpreter) forStatement(n *ast.ForStatement) lang.Value {
	i.enterScope()
	i.Do(n.Init)
	for i.Do(n.Test).Bool {
		i.Do(n.Body)
		i.Do(n.Update)
	}
	i.exitScope()
	return lang.NewUndefined()
}

func (i *Interpreter) callExpression(n *ast.CallExpression) lang.Value {
	args := []lang.Value{}
	for _, a := range n.Arguments {
		args = append(args, i.Do(a))
	}

	v := lang.NewUndefined()

	f := i.get(n.Callee.Name)
	if nf, ok := f.Obj.(*lang.NativeFunction); ok {
		nf.Function(args...)
	} else if lf, ok := f.Obj.(*lang.Function); ok {
		i.enterScope()
		for idx, a := range lf.Parameters {
			i.put(a.Name, args[idx])
		}
		v = i.Do(lf.Body)
		i.exitScope()
	} else {
		panic("unhandled function reference type")
	}

	return v
}

func (i *Interpreter) intLiteral(n *ast.IntLiteral) lang.Value {
	return lang.NewInt(n.Value)
}
