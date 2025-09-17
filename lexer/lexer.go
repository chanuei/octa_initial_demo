package lexer

import "strings"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenBlock
	TokenVar
	TokenAssign
	TokenPrint
	TokenNumber
	TokenIdent
	TokenLParen
	TokenRParen
	TokenNewline
)

type Token struct {
	Type  TokenType
	Value string
}

// 简单分词
func Lex(input string) []Token {
	var tokens []Token
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		for _, w := range words {
			switch w {
			case "block":
				tokens = append(tokens, Token{Type: TokenBlock, Value: w})
			case "var":
				tokens = append(tokens, Token{Type: TokenVar, Value: w})
			case "print":
				tokens = append(tokens, Token{Type: TokenPrint, Value: w})
			default:
				tokens = append(tokens, Token{Type: TokenIdent, Value: w})
			}
		}
	}
	return tokens
}
