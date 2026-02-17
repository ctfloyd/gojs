package tkn

import (
	"strconv"
	"unicode"
)

type TokenKind int

const (
	TokenKindIllegal TokenKind = iota
	TokenKindAmperandAmpersandEqual
	TokenKindAmpersand
	TokenKindAmpersandAmpersand
	TokenKindAmpersandEqual
	TokenKindAsterisk
	TokenKindAsteriskAsterisk
	TokenKindAsteriskAsteriskEqual
	TokenKindAsteriskEqual
	TokenKindCaret
	TokenKindCaretEqual
	TokenKindColon
	TokenKindComma
	TokenKindEOF
	TokenKindEqual
	TokenKindEqualEqual
	TokenKindEqualEqualEqual
	TokenKindEqualGreatherThan
	TokenKindExclamation
	TokenKindFor
	TokenKindFunction
	TokenKindGreaterThan
	TokenKindGreaterThanGreaterThan
	TokenKindGreaterThanGreaterThanEqual
	TokenKindGreaterThanGreaterThanGreaterThan
	TokenKindGreaterThanGreaterThanGreaterThanEqual
	TokenKindGreaterThanOrEqual
	TokenKindIdentifier
	TokenKindIf
	TokenKindIntLiteral
	TokenKindLeftBrace
	TokenKindLeftParen
	TokenKindLeftSquareBracket
	TokenKindLessThan
	TokenKindLessThanLessThan
	TokenKindLessThanOrEqual
	TokenKindLessThanLessThanEqual
	TokenKindMinus
	TokenKindMinusEqual
	TokenKindMinusMinus
	TokenKindNotEqual
	TokenKindNotEqualEqual
	TokenKindPercent
	TokenKindPercentEqual
	TokenKindPeriod
	TokenKindPipe
	TokenKindPipeEqual
	TokenKindPipePipe
	TokenKindPipePipeEqual
	TokenKindPlus
	TokenKindPlusEqual
	TokenKindPlusPlus
	TokenKindQuestion
	TokenKindQuestionPeriod
	TokenKindQuestionQuestion
	TokenKindQuestionQuestionEqual
	TokenKindReturn
	TokenKindRightBrace
	TokenKindRightParen
	TokenKindRightSquareBracket
	TokenKindSemicolon
	TokenKindSlash
	TokenKindSlashEqual
	TokenKindSpread
	TokenKindTilde
	TokenKindVar
)

func (tk TokenKind) String() string {
	switch tk {
	case TokenKindIllegal:
		return "Illegal"
	case TokenKindAmperandAmpersandEqual:
		return "AmperandAmpersandEqual"
	case TokenKindAmpersand:
		return "Ampersand"
	case TokenKindAmpersandAmpersand:
		return "AmpersandAmpersand"
	case TokenKindAmpersandEqual:
		return "AmpersandEqual"
	case TokenKindAsterisk:
		return "Asterisk"
	case TokenKindAsteriskAsterisk:
		return "AsteriskAsterisk"
	case TokenKindAsteriskAsteriskEqual:
		return "AsteriskAsteriskEqual"
	case TokenKindAsteriskEqual:
		return "AsteriskEqual"
	case TokenKindCaret:
		return "Caret"
	case TokenKindCaretEqual:
		return "CaretEqual"
	case TokenKindColon:
		return "Colon"
	case TokenKindComma:
		return "Comma"
	case TokenKindEOF:
		return "EOF"
	case TokenKindEqual:
		return "Equal"
	case TokenKindEqualEqual:
		return "EqualEqual"
	case TokenKindEqualEqualEqual:
		return "EqualEqualEqual"
	case TokenKindEqualGreatherThan:
		return "EqualGreatherThan"
	case TokenKindExclamation:
		return "Exclamation"
	case TokenKindFor:
		return "For"
	case TokenKindFunction:
		return "Function"
	case TokenKindGreaterThan:
		return "GreaterThan"
	case TokenKindGreaterThanGreaterThan:
		return "GreaterThanGreaterThan"
	case TokenKindGreaterThanGreaterThanEqual:
		return "GreaterThanGreaterThanEqual"
	case TokenKindGreaterThanGreaterThanGreaterThan:
		return "GreaterThanGreaterThanGreaterThan"
	case TokenKindGreaterThanGreaterThanGreaterThanEqual:
		return "GreaterThanGreaterThanGreaterThanEqual"
	case TokenKindGreaterThanOrEqual:
		return "GreaterThanOrEqual"
	case TokenKindIdentifier:
		return "Identifier"
	case TokenKindIf:
		return "If"
	case TokenKindIntLiteral:
		return "IntLiteral"
	case TokenKindLeftBrace:
		return "LeftBrace"
	case TokenKindLeftParen:
		return "LeftParen"
	case TokenKindLeftSquareBracket:
		return "LeftSquareBracket"
	case TokenKindLessThan:
		return "LessThan"
	case TokenKindLessThanLessThan:
		return "LessThanLessThan"
	case TokenKindLessThanOrEqual:
		return "LessThanOrEqual"
	case TokenKindLessThanLessThanEqual:
		return "LessThanLessThanEqual"
	case TokenKindMinus:
		return "Minus"
	case TokenKindMinusEqual:
		return "MinusEqual"
	case TokenKindMinusMinus:
		return "MinusMinus"
	case TokenKindNotEqual:
		return "NotEqual"
	case TokenKindNotEqualEqual:
		return "NotEqualEqual"
	case TokenKindPercent:
		return "Percent"
	case TokenKindPercentEqual:
		return "PercentEqual"
	case TokenKindPeriod:
		return "Period"
	case TokenKindPipe:
		return "Pipe"
	case TokenKindPipeEqual:
		return "PipeEqual"
	case TokenKindPipePipe:
		return "PipePipe"
	case TokenKindPipePipeEqual:
		return "PipePipeEqual"
	case TokenKindPlus:
		return "Plus"
	case TokenKindPlusEqual:
		return "PlusEqual"
	case TokenKindPlusPlus:
		return "PlusPlus"
	case TokenKindQuestion:
		return "Question"
	case TokenKindQuestionPeriod:
		return "QuestionPeriod"
	case TokenKindQuestionQuestion:
		return "QuestionQuestion"
	case TokenKindQuestionQuestionEqual:
		return "QuestionQuestionEqual"
	case TokenKindReturn:
		return "Return"
	case TokenKindRightBrace:
		return "RightBrace"
	case TokenKindRightParen:
		return "RightParen"
	case TokenKindRightSquareBracket:
		return "RightSquareBracket"
	case TokenKindSemicolon:
		return "Semicolon"
	case TokenKindSlash:
		return "Slash"
	case TokenKindSlashEqual:
		return "SlashEqual"
	case TokenKindSpread:
		return "Spread"
	case TokenKindTilde:
		return "Tilde"
	case TokenKindVar:
		return "Var"
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

// https://tc39.es/ecma262/#sec-punctuators
func (t *Tokenizer) resolvePunctuator(start rune) Token {
	p := punctuator[start]
	for next, ok := p.match(t.peek()); ok; next, ok = p.match(t.peek()) {
		t.consume()
		p = next
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
