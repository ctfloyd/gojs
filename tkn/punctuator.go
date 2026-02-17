// Punctuator
// https://tc39.es/ecma262/#sec-punctuators

package tkn

type Punctuator struct {
	value         rune
	token         TokenKind
	continuations []Punctuator
}

func (p Punctuator) match(r rune) (Punctuator, bool) {
	for _, c := range p.continuations {
		if c.value == r {
			return c, true
		}
	}
	return Punctuator{}, false
}

func isPunctuatorStart(ch rune) bool {
	_, ok := punctuator[ch]
	return ok
}

var punctuator = map[rune]Punctuator{
	//
	// OptionalChainingPunctuator
	'?': {
		value: '?',
		token: TokenKindQuestion,
		continuations: []Punctuator{
			{
				value: '.',
				token: TokenKindQuestionPeriod,
				//
				// TODO (ctfloyd): this should be illegal if the next digit is a decimal digit.
				// 				   see: https://tc39.es/ecma262/#prod-OptionalChainingPunctuator
			},
			{
				value: '?',
				token: TokenKindQuestionQuestion,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindQuestionQuestionEqual,
					},
				},
			},
		},
	},

	// OtherPunctuator
	'{': {
		value: '{',
		token: TokenKindLeftBrace,
	},
	'(': {
		value: '(',
		token: TokenKindLeftParen,
	},
	')': {
		value: ')',
		token: TokenKindRightParen,
	},
	'[': {
		value: '[',
		token: TokenKindLeftSquareBracket,
	},
	']': {
		value: ']',
		token: TokenKindRightSquareBracket,
	},
	'.': {
		value: '.',
		token: TokenKindPeriod,
		continuations: []Punctuator{
			{
				value: '.',
				token: TokenKindIllegal,
				continuations: []Punctuator{
					{
						value: '.',
						token: TokenKindSpread,
					},
				},
			},
		},
	},
	';': {
		value: ';',
		token: TokenKindSemicolon,
	},
	',': {
		value: ',',
		token: TokenKindComma,
	},
	'<': {
		value: '<',
		token: TokenKindLessThan,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindLessThanOrEqual,
			},
			{
				value: '<',
				token: TokenKindLessThanLessThan,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindLessThanLessThanEqual,
					},
				},
			},
		},
	},
	'>': {
		value: '>',
		token: TokenKindGreaterThan,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindGreaterThanOrEqual,
			},
			{
				value: '>',
				token: TokenKindGreaterThanGreaterThan,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindGreaterThanGreaterThanEqual,
					},
					{
						value: '>',
						token: TokenKindGreaterThanGreaterThanGreaterThan,
						continuations: []Punctuator{
							{
								value: '=',
								token: TokenKindGreaterThanGreaterThanGreaterThanEqual,
							},
						},
					},
				},
			},
		},
	},
	'=': {
		value: '=',
		token: TokenKindEqual,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindEqualEqual,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindEqualEqualEqual,
					},
				},
			},
			{
				value: '>',
				token: TokenKindEqualGreatherThan,
			},
		},
	},
	'!': {
		value: '!',
		token: TokenKindExclamation,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindNotEqual,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindNotEqualEqual,
					},
				},
			},
		},
	},
	'+': {
		value: '+',
		token: TokenKindPlus,
		continuations: []Punctuator{
			{
				value: '+',
				token: TokenKindPlusPlus,
			},
			{
				value: '=',
				token: TokenKindPlusEqual,
			},
		},
	},
	'-': {
		value: '-',
		token: TokenKindMinus,
		continuations: []Punctuator{
			{
				value: '-',
				token: TokenKindMinusMinus,
			},
			{
				value: '=',
				token: TokenKindMinusEqual,
			},
		},
	},
	'*': {
		value: '*',
		token: TokenKindAsterisk,
		continuations: []Punctuator{
			{
				value: '*',
				token: TokenKindAsteriskAsterisk,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindAsteriskAsteriskEqual,
					},
				},
			},
			{
				value: '=',
				token: TokenKindAsteriskEqual,
			},
		},
	},
	'%': {
		value: '%',
		token: TokenKindPercent,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindPercentEqual,
			},
		},
	},
	'&': {
		value: '&',
		token: TokenKindAmpersand,
		continuations: []Punctuator{
			{
				value: '&',
				token: TokenKindAmpersandAmpersand,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindAmperandAmpersandEqual,
					},
				},
			},
			{
				value: '=',
				token: TokenKindAmpersandEqual,
			},
		},
	},
	'|': {
		value: '|',
		token: TokenKindPipe,
		continuations: []Punctuator{
			{
				value: '|',
				token: TokenKindPipePipe,
				continuations: []Punctuator{
					{
						value: '=',
						token: TokenKindPipePipeEqual,
					},
				},
			},
			{
				value: '=',
				token: TokenKindPipeEqual,
			},
		},
	},
	'^': {
		value: '^',
		token: TokenKindCaret,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindCaretEqual,
			},
		},
	},
	'~': {
		value: '~',
		token: TokenKindTilde,
	},
	':': {
		value: ':',
		token: TokenKindColon,
	},

	//
	// DivPunctuator
	'/': {
		value: '/',
		token: TokenKindSlash,
		continuations: []Punctuator{
			{
				value: '=',
				token: TokenKindSlashEqual,
			},
		},
	},

	//
	// RightBracePunctuator
	'}': {
		value: '}',
		token: TokenKindRightBrace,
	},
}
