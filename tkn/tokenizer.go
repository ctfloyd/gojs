package tkn

import (
	"strconv"
	"unicode"
)

type TokenKind int

const (
	TokenKindIllegal TokenKind = iota
	TokenKindComma
	TokenKindEqual
	TokenKindEOF
	TokenKindFor
	TokenKindFunction
	TokenKindGreaterThan
	TokenKindIdentifier
	TokenKindIf
	TokenKindIntLiteral
	TokenKindLeftBrace
	TokenKindLeftParen
	TokenKindLessThan
	TokenKindLessThanOrEqual
	TokenKindPlus
	TokenKindPlusPlus
	TokenKindReturn
	TokenKindRightBrace
	TokenKindRightParen
	TokenKindSemicolon
	TokenKindVar
)

func (tk TokenKind) String() string {
	switch tk {
	case TokenKindIllegal:
		return "Illegal"
	case TokenKindComma:
		return "Comma"
	case TokenKindEqual:
		return "Equal"
	case TokenKindFor:
		return "For"
	case TokenKindFunction:
		return "Function"
	case TokenKindGreaterThan:
		return "GreaterThan"
	case TokenKindIdentifier:
		return "Identifier"
	case TokenKindLeftBrace:
		return "LeftBrace"
	case TokenKindLeftParen:
		return "LeftParen"
	case TokenKindLessThan:
		return "LessThan"
	case TokenKindLessThanOrEqual:
		return "LessThanOrEqual"
	case TokenKindPlus:
		return "Plus"
	case TokenKindPlusPlus:
		return "PlusPlus"
	case TokenKindReturn:
		return "Return"
	case TokenKindRightBrace:
		return "RightBrace"
	case TokenKindRightParen:
		return "RightParen"
	case TokenKindSemicolon:
		return "Semicolon"
	case TokenKindVar:
		return "Var"
	case TokenKindEOF:
		return "EOF"
	default:
		return "Unknown"

	}
}

// https://tc39.es/ecma262/#sec-white-space
var whitespace = map[rune]interface{}{
	'\u0009': nil, // Character Tabulation <TAB>
	'\u000B': nil, // Line Tabulation <VT>
	'\u000C': nil, // Form Feed (FF) <FF>
	'\uFEFF': nil, // Zero Width No-Break Space <ZWNBSP>
}

// https://tc39.es/ecma262/#sec-white-space
func isWhitespace(r rune) bool {
	if _, ok := whitespace[r]; ok {
		return true
	}

	return unicode.IsSpace(r)
}

type Punctuator struct {
	value        rune
	token        TokenKind
	continuation *Punctuator
}

var punctuator = map[rune]Punctuator{
	'<': {
		value: '<',
		token: TokenKindLessThan,
		continuation: &Punctuator{
			value: '=',
			token: TokenKindLessThanOrEqual,
		},
	},
	'>': {
		value: '>',
		token: TokenKindGreaterThan,
		// TODO (c.floyd): continuation
	},
	'=': {
		value: '=',
		token: TokenKindEqual,
	},
	'{': {
		value: '{',
		token: TokenKindLeftBrace,
	},
	'}': {
		value: '}',
		token: TokenKindRightBrace,
	},
	'(': {
		value: '(',
		token: TokenKindLeftParen,
	},
	')': {
		value: ')',
		token: TokenKindRightParen,
	},
	'+': {
		value: '+',
		token: TokenKindPlus,
		continuation: &Punctuator{
			value: '+',
			token: TokenKindPlusPlus,
		},
	},
	',': {
		value: ';',
		token: TokenKindComma,
	},
	';': {
		value: ';',
		token: TokenKindSemicolon,
	},
}

func isPunctuatorStart(ch rune) bool {
	if _, ok := punctuator[ch]; ok {
		return true
	}

	return false
}

type Location struct {
	Line, Column int
}
type Token struct {
	Location Location
	Kind     TokenKind
	Value    string
}

func NewToken(kind TokenKind, line, column int) Token {
	return NewTokenWithValue(kind, line, column, "")
}

func NewTokenWithValue(kind TokenKind, line, column int, value string) Token {
	return Token{
		Kind:     kind,
		Location: Location{Line: line, Column: column},
		Value:    value,
	}
}

type Tokenizer struct {
	text    string
	current int

	line   int
	column int
}

func (t *Tokenizer) Tokenize(text string) []Token {
	t.text = text

	tokens := make([]Token, 0, 128)

	buffer := ""
	for t.current < len(text) {
		ch := t.consume()

		if isWhitespace(ch) {
			if token, ok := t.resolveBuffer(buffer); ok {
				tokens = append(tokens, token)
			}
			buffer = ""
		} else if isPunctuatorStart(ch) {
			if token, ok := t.resolveBuffer(buffer); ok {
				tokens = append(tokens, token)
			}
			buffer = ""

			token := t.resolvePunctuator(ch)
			tokens = append(tokens, token)
		} else {
			buffer += string(ch)
		}
	}

	if token, ok := t.resolveBuffer(buffer); ok {
		tokens = append(tokens, token)
	}

	tokens = append(tokens, Token{Kind: TokenKindEOF, Location: Location{Line: t.line, Column: t.column}})
	return tokens
}

func (t *Tokenizer) peek() rune {
	if t.current >= len(t.text) {
		return -1
	}

	return rune(t.text[t.current])
}

func (t *Tokenizer) consume() rune {
	if t.current >= len(t.text) {
		return -1
	}

	ch := rune(t.text[t.current])
	t.current += 1

	// TODO (c.floyd): This should probably use line separators?
	if ch == '\n' {
		t.column = 0
		t.line += 1
	} else {
		t.column += 1
	}

	return ch
}

func (t *Tokenizer) resolvePunctuator(start rune) Token {
	p := punctuator[start]

	for p.continuation != nil && p.continuation.value == t.peek() {
		t.consume()

		p = *p.continuation
	}
	return NewToken(p.token, t.line, t.column)
}

func (t *Tokenizer) resolveBuffer(buffer string) (Token, bool) {
	if len(buffer) == 0 {
		return Token{}, false
	}

	line, column := t.line, t.column

	if buffer == "function" {
		return NewToken(TokenKindFunction, line, column), true
	}

	if buffer == "var" {
		return NewToken(TokenKindVar, line, column), true
	}

	if buffer == "return" {
		return NewToken(TokenKindReturn, line, column), true
	}

	if buffer == "if" {
		return NewToken(TokenKindIf, line, column), true
	}

	if buffer == "for" {
		return NewToken(TokenKindFor, line, column), true
	}

	if _, err := strconv.Atoi(buffer); err == nil {
		return NewTokenWithValue(TokenKindIntLiteral, line, column, buffer), true
	}

	return NewTokenWithValue(TokenKindIdentifier, line, column, buffer), true
}
