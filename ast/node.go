package ast

type Node interface {
	Node()
}

type Statement interface {
	Node()
	_Statement()
}

type Expression interface {
	Node()
	_Expression()
}

type Program struct {
	Start, End int
	Body       []Node
}

func (p *Program) Node() {}

type Identifier struct {
	Start, End int
	Name       string
}

func (i *Identifier) Node()        {}
func (i *Identifier) _Statement()  {}
func (i *Identifier) _Expression() {}

type FunctionDeclaration struct {
	Start, End                   int
	Expression, Generator, Async bool
	Id                           Identifier
	Parameters                   []Identifier
	Body                         Statement
}

func (f *FunctionDeclaration) Node()       {}
func (f *FunctionDeclaration) _Statement() {}

type BlockStatement struct {
	Start, End int
	Body       []Statement
}

func (b *BlockStatement) Node()       {}
func (b *BlockStatement) _Statement() {}

type VariableDeclaration struct {
	Start, End   int
	Declarations []*VariableDeclarator
	Kind         string
}

func (v *VariableDeclaration) Node()       {}
func (v *VariableDeclaration) _Statement() {}

type VariableDeclarator struct {
	Start, End int
	Id         *Identifier
	Init       Expression
}

func (v *VariableDeclarator) Node() {}

type ArrayExpression struct {
	Start, End int
	Elements   []Expression
}

func (a *ArrayExpression) Node()        {}
func (a *ArrayExpression) _Expression() {}

type MemberExpression struct {
	Start, End int
	Object     Expression
	Property   Expression
}

func (m *MemberExpression) Node()        {}
func (m *MemberExpression) _Expression() {}

type BinaryExpression struct {
	Start, End  int
	Left, Right Expression
	Operator    string
}

func (b *BinaryExpression) Node()        {}
func (b *BinaryExpression) _Expression() {}

type UpdateExpression struct {
	Argument Expression
	Operator string
}

func (u *UpdateExpression) Node()        {}
func (u *UpdateExpression) _Expression() {}

type ReturnStatement struct {
	Start, End int
	Argument   Expression
}

func (r *ReturnStatement) Node()       {}
func (r *ReturnStatement) _Statement() {}

type IfStatement struct {
	Start, End int
	Test       Expression
	Consequent Statement
}

func (i *IfStatement) Node()       {}
func (i *IfStatement) _Statement() {}

type ExpressionStatement struct {
	Start, End int
	Expression Expression
}

func (e *ExpressionStatement) Node()       {}
func (e *ExpressionStatement) _Statement() {}

type CallExpression struct {
	Start, End int
	Callee     *Identifier
	Arguments  []Expression
	Optional   bool
}

func (c *CallExpression) Node()        {}
func (c *CallExpression) _Expression() {}

type IntLiteral struct {
	Start, End int
	Value      int
}

func (i *IntLiteral) Node()        {}
func (i *IntLiteral) _Expression() {}

type ForStatement struct {
	Start, End   int
	Init         Statement
	Test, Update Expression
	Body         Statement
}

func (f *ForStatement) Node()       {}
func (f *ForStatement) _Statement() {}

type AssignmentExpression struct {
	Start, End  int
	Operator    string
	Left, Right Expression
}

func (a *AssignmentExpression) Node()        {}
func (a *AssignmentExpression) _Expression() {}
