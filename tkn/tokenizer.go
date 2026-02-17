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
}

func (t *Tokenizer) Tokenize(text string) []Token {
	tokens := make([]Token, 0, 128)

	line := 0
	column := 0

	buffer := ""
	for _, ch := range text {
		if isWhitespace(ch) {
			if token, ok := t.resolveBuffer(buffer, line, column); ok {
				tokens = append(tokens, token)
			}
			buffer = ""
		} else if ch == '(' || ch == ',' || ch == ')' || ch == ';' {
			if token, ok := t.resolveBuffer(buffer, line, column); ok {
				tokens = append(tokens, token)
			}

			buffer = string(ch)
			if token, ok := t.resolveBuffer(buffer, line, column); ok {
				tokens = append(tokens, token)
			}
			buffer = ""
		} else {
			buffer += string(ch)
		}

		column += 1
		if ch == '\n' {
			line++
			column = 0
		}
	}

	if token, ok := t.resolveBuffer(buffer, line, column); ok {
		tokens = append(tokens, token)
	}

	tokens = append(tokens, Token{Kind: TokenKindEOF, Location: Location{Line: line, Column: column}})
	return tokens
}

func (t *Tokenizer) resolveBuffer(buffer string, line, column int) (Token, bool) {
	if len(buffer) == 0 {
		return Token{}, false
	}

	if buffer == "function" {
		return NewToken(TokenKindFunction, line, column), true
	}

	if buffer == "var" {
		return NewToken(TokenKindVar, line, column), true
	}

	if buffer == "return" {
		return NewToken(TokenKindReturn, line, column), true
	}

	if buffer == "<" {
		return NewToken(TokenKindLessThan, line, column), true
	}

	if buffer == "if" {
		return NewToken(TokenKindIf, line, column), true
	}

	if buffer == "for" {
		return NewToken(TokenKindFor, line, column), true
	}

	if buffer == "," {
		return NewToken(TokenKindComma, line, column), true
	}

	if buffer == "(" {
		return NewToken(TokenKindLeftParen, line, column), true
	}

	if buffer == ")" {
		return NewToken(TokenKindRightParen, line, column), true
	}

	if buffer == "{" {
		return NewToken(TokenKindLeftBrace, line, column), true
	}

	if buffer == "}" {
		return NewToken(TokenKindRightBrace, line, column), true
	}

	if buffer == "=" {
		return NewToken(TokenKindEqual, line, column), true
	}

	if buffer == "+" {
		return NewToken(TokenKindPlus, line, column), true
	}

	if buffer == "++" {
		return NewToken(TokenKindPlusPlus, line, column), true
	}

	if buffer == ";" {
		return NewToken(TokenKindSemicolon, line, column), true
	}

	if buffer == ">" {
		return NewToken(TokenKindGreaterThan, line, column), true
	}

	if _, err := strconv.Atoi(buffer); err == nil {
		return NewTokenWithValue(TokenKindIntLiteral, line, column, buffer), true
	}

	return NewTokenWithValue(TokenKindIdentifier, line, column, buffer), true
}
