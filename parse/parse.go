package parse

import (
	"gojs/ast"
	"gojs/tkn"
	"strconv"
)

type Parser struct {
	tokens []tkn.Token
	offset int
}

func NewParser(tokens []tkn.Token) *Parser {
	return &Parser{tokens: tokens, offset: 0}
}

func (p *Parser) Parse() ast.Program {
	nodes := make([]ast.Node, 0, 10)

	n := p.parseStatement()
	for n != nil {
		nodes = append(nodes, n)
		n = p.parseStatement()
	}

	return ast.Program{
		Body: nodes,
	}
}

func (p *Parser) parseStatement() ast.Statement {
	if p.match(tkn.TokenKindEOF) {
		return nil
	}

	if p.matchesExpression() {
		return &ast.ExpressionStatement{
			Expression: p.parseExpression(),
		}
	}

	if p.match(tkn.TokenKindFunction) {
		return p.parseFunction()
	} else if p.match(tkn.TokenKindVar) {
		return p.parseVariableDeclaration()
	} else if p.match(tkn.TokenKindReturn) {
		return p.parseReturn()
	} else if p.match(tkn.TokenKindIf) {
		return p.parseIf()
	} else if p.match(tkn.TokenKindLeftBrace) {
		return p.parseBlockStatement()
	} else if p.match(tkn.TokenKindFor) {
		return p.parseForStatement()
	}

	panic("unknown kind: " + p.kind().String())
}

func (p *Parser) kind() tkn.TokenKind {
	return p.tokens[p.offset].Kind
}

func (p *Parser) value() string {
	return p.tokens[p.offset].Value
}

func (p *Parser) consume(kind tkn.TokenKind) tkn.Token {
	if p.kind() != kind {
		panic("expected kind: " + kind.String() + " but got kind: " + p.kind().String())
	}
	p.offset++
	return p.tokens[p.offset-1]
}

func (p *Parser) match(kind tkn.TokenKind) bool {
	return p.kind() == kind
}

func (p *Parser) parseFunction() *ast.FunctionDeclaration {
	p.consume(tkn.TokenKindFunction)
	name := p.consume(tkn.TokenKindIdentifier).Value
	p.consume(tkn.TokenKindLeftParen)

	args := make([]ast.Identifier, 0)
	for p.match(tkn.TokenKindIdentifier) {
		args = append(args, ast.Identifier{Name: p.consume(tkn.TokenKindIdentifier).Value})
		if p.match(tkn.TokenKindComma) {
			p.consume(tkn.TokenKindComma)
		}
	}
	p.consume(tkn.TokenKindRightParen)

	return &ast.FunctionDeclaration{
		Id:         ast.Identifier{Name: name},
		Parameters: args,
		Body:       p.parseBlockStatement(),
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.consume(tkn.TokenKindLeftBrace)

	statements := make([]ast.Statement, 0, 10)
	for !p.match(tkn.TokenKindRightBrace) {
		node := p.parseStatement()
		statement := node.(ast.Statement)
		statements = append(statements, statement)
	}
	p.consume(tkn.TokenKindRightBrace)

	return &ast.BlockStatement{
		Body: statements,
	}
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	p.consume(tkn.TokenKindFor)
	p.consume(tkn.TokenKindLeftParen)
	init := p.parseStatement()
	p.consume(tkn.TokenKindSemicolon)
	test := p.parseExpression()
	p.consume(tkn.TokenKindSemicolon)
	update := p.parseExpression()
	p.consume(tkn.TokenKindRightParen)
	body := p.parseStatement()

	return &ast.ForStatement{
		Init:   init,
		Test:   test,
		Update: update,
		Body:   body,
	}
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	p.consume(tkn.TokenKindVar)
	name := p.consume(tkn.TokenKindIdentifier).Value
	p.consume(tkn.TokenKindEqual)
	init := p.parseExpression()
	return &ast.VariableDeclaration{
		Declarations: []*ast.VariableDeclarator{
			{
				Id:   &ast.Identifier{Name: name},
				Init: init,
			},
		},
		Kind: "var",
	}

}

func (p *Parser) parseReturn() *ast.ReturnStatement {
	p.consume(tkn.TokenKindReturn)

	var expr ast.Expression
	if p.matchesExpression() {
		expr = p.parseExpression()
	}

	if p.match(tkn.TokenKindSemicolon) {
		p.consume(tkn.TokenKindSemicolon)
	}
	return &ast.ReturnStatement{
		Argument: expr,
	}
}

func (p *Parser) parseIf() *ast.IfStatement {
	p.consume(tkn.TokenKindIf)
	p.consume(tkn.TokenKindLeftParen)
	test := p.parseExpression()
	p.consume(tkn.TokenKindRightParen)
	consequent := p.parseStatement()

	return &ast.IfStatement{
		Test:       test,
		Consequent: consequent,
	}
}

func (p *Parser) parseCallExpression(identifier *ast.Identifier) *ast.CallExpression {
	p.consume(tkn.TokenKindLeftParen)
	args := make([]ast.Expression, 0)
	for !p.match(tkn.TokenKindRightParen) {
		args = append(args, p.parseExpression())
		if p.match(tkn.TokenKindComma) {
			p.consume(tkn.TokenKindComma)
		}
	}
	p.consume(tkn.TokenKindRightParen)

	return &ast.CallExpression{
		Callee:    identifier,
		Arguments: args,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	expr := p.parsePrimaryExpression()
	for p.matchesSecondaryExpression() {
		expr = p.parseSecondaryExpression(expr)
	}
	return expr
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	if p.match(tkn.TokenKindLeftParen) {
		p.consume(tkn.TokenKindLeftParen)
		expr := p.parseExpression()
		p.consume(tkn.TokenKindRightParen)
		return expr
	} else if p.match(tkn.TokenKindIdentifier) {
		return &ast.Identifier{Name: p.consume(tkn.TokenKindIdentifier).Value}
	} else if p.match(tkn.TokenKindIntLiteral) {
		value, _ := strconv.Atoi(p.consume(tkn.TokenKindIntLiteral).Value)
		return &ast.IntLiteral{Value: value}
	} else if p.match(tkn.TokenKindLeftSquareBracket) {
		var elements []ast.Expression
		p.consume(tkn.TokenKindLeftSquareBracket)
		for !p.match(tkn.TokenKindRightSquareBracket) {
			elements = append(elements, p.parseExpression())
			if p.match(tkn.TokenKindComma) {
				p.consume(tkn.TokenKindComma)
			}
		}
		p.consume(tkn.TokenKindRightSquareBracket)
		return &ast.ArrayExpression{Elements: elements}
	} else if p.match(tkn.TokenKindLeftBrace) {
		var properties []ast.Property
		p.consume(tkn.TokenKindLeftBrace)
		for !p.match(tkn.TokenKindRightBrace) {
			key := p.parseExpression()
			p.consume(tkn.TokenKindColon)
			value := p.parseExpression()
			properties = append(properties, ast.Property{Key: key, Value: value})
		}
		p.consume(tkn.TokenKindRightBrace)
		return &ast.ObjectExpression{Properties: properties}
	} else {
		panic("yoo yoo")
	}
}

func (p *Parser) parseSecondaryExpression(lhs ast.Node) ast.Expression {
	if p.match(tkn.TokenKindPlus) {
		p.consume(tkn.TokenKindPlus)
		return &ast.BinaryExpression{Left: lhs.(ast.Expression), Right: p.parseExpression(), Operator: "+"}
	} else if p.match(tkn.TokenKindPlusPlus) {
		p.consume(tkn.TokenKindPlusPlus)
		return &ast.UpdateExpression{Argument: lhs.(ast.Expression), Operator: "++"}
	} else if p.match(tkn.TokenKindMinusMinus) {
		p.consume(tkn.TokenKindMinusMinus)
		return &ast.UpdateExpression{Argument: lhs.(ast.Expression), Operator: "--"}
	} else if p.match(tkn.TokenKindGreaterThan) {
		p.consume(tkn.TokenKindGreaterThan)
		return &ast.BinaryExpression{Left: lhs.(ast.Expression), Right: p.parseExpression(), Operator: ">"}
	} else if p.match(tkn.TokenKindLessThan) {
		p.consume(tkn.TokenKindLessThan)
		return &ast.BinaryExpression{Left: lhs.(ast.Expression), Right: p.parseExpression(), Operator: "<"}
	} else if p.match(tkn.TokenKindEqual) {
		p.consume(tkn.TokenKindEqual)
		return &ast.AssignmentExpression{Left: lhs.(ast.Expression), Right: p.parseExpression(), Operator: "="}
	} else if p.match(tkn.TokenKindLeftParen) {
		return p.parseCallExpression(lhs.(*ast.Identifier))
	} else if p.match(tkn.TokenKindLeftSquareBracket) {
		p.consume(tkn.TokenKindLeftSquareBracket)
		property := p.parseExpression()
		p.consume(tkn.TokenKindRightSquareBracket)
		return &ast.MemberExpression{Object: lhs.(*ast.Identifier), Property: property}
	} else {
		panic("yoo yoo 2")
	}
}

func (p *Parser) matchesStatement() bool {
	k := p.kind()
	return p.matchesExpression() ||
		k == tkn.TokenKindFunction ||
		k == tkn.TokenKindReturn ||
		k == tkn.TokenKindVar ||
		k == tkn.TokenKindLeftBrace ||
		k == tkn.TokenKindFor
}

func (p *Parser) matchesExpression() bool {
	k := p.kind()
	return k == tkn.TokenKindIntLiteral ||
		k == tkn.TokenKindIdentifier ||
		k == tkn.TokenKindLeftParen ||
		k == tkn.TokenKindLeftSquareBracket ||
		k == tkn.TokenKindLeftBrace
}

func (p *Parser) matchesSecondaryExpression() bool {
	k := p.kind()
	return k == tkn.TokenKindPlus ||
		k == tkn.TokenKindPlusPlus ||
		k == tkn.TokenKindMinusMinus ||
		k == tkn.TokenKindGreaterThan ||
		k == tkn.TokenKindLeftParen ||
		k == tkn.TokenKindLessThan ||
		k == tkn.TokenKindEqual ||
		k == tkn.TokenKindLeftSquareBracket
}
