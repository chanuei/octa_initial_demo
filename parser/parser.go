package parser

import (
	"fmt"
	"strconv"

	"octa/ast"
	"octa/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() lexer.Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) expect(tt lexer.TokenType) lexer.Token {
	tok := p.next()
	if tok.Type != tt {
		panic(fmt.Sprintf("expected %v, got %v", tt, tok))
	}
	return tok
}

func Parse(tokens []lexer.Token) *ast.FuncStmt {
	p := New(tokens)
	tok := p.next()
	if tok.Type != lexer.TokenBlock {
		panic("expected 'block'")
	}
	funcName := p.expect(lexer.TokenIdentifier)
	p.expect(lexer.TokenLParen)
	p.expect(lexer.TokenRParen)

	f := &ast.FuncStmt{Name: funcName.Value}

	for p.pos < len(p.tokens) && p.peek().Type != lexer.TokenEOF {
		stmt := parseStmt(p)
		if stmt != nil {
			f.Body = append(f.Body, stmt)
		}
	}
	return f
}

func parseStmt(p *Parser) ast.Stmt {
	tok := p.peek()
	switch tok.Type {
	case lexer.TokenVar:
		p.next()
		nameTok := p.expect(lexer.TokenIdentifier)
		p.expect(lexer.TokenAssign)
		valTok := p.expect(lexer.TokenNumber)
		val, _ := strconv.Atoi(valTok.Value)
		return &ast.VarDeclStmt{
			Name: nameTok.Value,
			Expr: ast.NumberExpr{Value: val},
		}
	case lexer.TokenIdentifier:
		name := p.next().Value
		p.expect(lexer.TokenAssign)
		valTok := p.expect(lexer.TokenNumber)
		val, _ := strconv.Atoi(valTok.Value)
		return &ast.AssignStmt{
			Name: name,
			Expr: ast.NumberExpr{Value: val},
		}
	case lexer.TokenPrint:
		p.next()
		p.expect(lexer.TokenLParen)
		varName := p.expect(lexer.TokenIdentifier)
		p.expect(lexer.TokenRParen)
		return &ast.PrintStmt{
			Expr: ast.VarExpr{Name: varName.Value},
		}
	default:
		p.next()
		return nil
	}
}
