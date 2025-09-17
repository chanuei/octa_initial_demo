package lexer

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenVar
	TokenPrint
	TokenBlock
	TokenAssign // =
	TokenLParen // (
	TokenRParen // )
	TokenEOL
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input []rune
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) Next() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	ch := l.input[l.pos]

	// 标识符/关键字
	if unicode.IsLetter(ch) {
		start := l.pos
		for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos])) {
			l.pos++
		}
		word := string(l.input[start:l.pos])
		switch word {
		case "var":
			return Token{Type: TokenVar, Value: word}
		case "print":
			return Token{Type: TokenPrint, Value: word}
		case "block":
			return Token{Type: TokenBlock, Value: word}
		default:
			return Token{Type: TokenIdentifier, Value: word}
		}
	}

	// 数字
	if unicode.IsDigit(ch) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: TokenNumber, Value: string(l.input[start:l.pos])}
	}

	// 符号
	l.pos++
	switch ch {
	case '=':
		return Token{Type: TokenAssign, Value: "="}
	case '(':
		return Token{Type: TokenLParen, Value: "("}
	case ')':
		return Token{Type: TokenRParen, Value: ")"}
	case '\n':
		return Token{Type: TokenEOL, Value: "\n"}
	}

	panic(fmt.Sprintf("unknown char: %c", ch))
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t' || l.input[l.pos] == '\r') {
		l.pos++
	}
}
