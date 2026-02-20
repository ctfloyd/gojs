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
	case *ast.ArrayExpression:
		return i.arrayExpression(n)
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
	case *ast.MemberExpression:
		return i.memberExpression(n)
	case *ast.Program:
		return i.program(n)
	case *ast.ReturnStatement:
		return i.returnStatement(n)
	case *ast.UpdateExpression:
		return i.updateExpression(n)
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

func (i *Interpreter) arrayExpression(n *ast.ArrayExpression) lang.Value {
	results := make([]lang.Value, len(n.Elements))
	for ix, e := range n.Elements {
		results[ix] = i.Do(e)
	}

	return lang.NewObj(&lang.Array{Store: results})
}

func (i *Interpreter) assignmentExpression(n *ast.AssignmentExpression) lang.Value {
	update := i.Do(n.Right)
	if identifier, ok := n.Left.(*ast.Identifier); ok {
		i.put(identifier.Name, update)
	} else if member, ok := n.Left.(*ast.MemberExpression); ok {
		array, idx, _ := i.resolveArrayMemberExpression(member)
		array.Store[idx] = update
	} else {
		panic("unsupported assignment expression")
	}

	return lang.NewUndefined()
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

func (i *Interpreter) identifierUpdateExpression(n *ast.UpdateExpression, identifier *ast.Identifier) lang.Value {
	arg := i.Do(n.Argument)
	if n.Operator == "++" {
		update := lang.NewInt(arg.Int + 1)
		i.put(identifier.Name, update)
		return update
	} else if n.Operator == "--" {
		update := lang.NewInt(arg.Int - 1)
		i.put(identifier.Name, update)
		return update
	} else {
		panic("unsupported operation")
	}
}

func (i *Interpreter) memberUpdateExpression(n *ast.UpdateExpression, me *ast.MemberExpression) lang.Value {
	array, idx, currentValue := i.resolveArrayMemberExpression(me)
	if n.Operator == "++" {
		update := lang.NewInt(currentValue.Int + 1)
		array.Store[idx] = update
		return update
	} else if n.Operator == "--" {
		update := lang.NewInt(currentValue.Int - 1)
		array.Store[idx] = update
		return update
	} else {
		panic("unsupported operation")
	}

}

func (i *Interpreter) updateExpression(n *ast.UpdateExpression) lang.Value {
	if identifier, ok := n.Argument.(*ast.Identifier); ok {
		return i.identifierUpdateExpression(n, identifier)
	}

	if member, ok := n.Argument.(*ast.MemberExpression); ok {
		return i.memberUpdateExpression(n, member)
	}

	panic("unsupported expression type")

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

func (i *Interpreter) memberExpression(n *ast.MemberExpression) lang.Value {
	_, _, val := i.resolveArrayMemberExpression(n)
	return val
}

func (i *Interpreter) resolveArrayMemberExpression(n *ast.MemberExpression) (*lang.Array, int, lang.Value) {
	o := i.Do(n.Object)
	if o.Type != lang.ValueTypeObj {
		panic("invalid object")
	}

	array, ok := o.Obj.(*lang.Array)
	if !ok {
		panic("object is not an array")
	}

	index := i.Do(n.Property)
	if index.Type != lang.ValueTypeInt {
		panic("invalid index")
	}

	if index.Int >= len(array.Store) {
		panic("index out of range")
	}

	return array, index.Int, array.Store[index.Int]
}
