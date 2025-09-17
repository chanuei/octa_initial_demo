package parser

import (
	"octa/ast"
	"octa/lexer"
	"strconv"
)

func Parse(tokens []lexer.Token) *ast.FuncStmt {
	// 简单示例，只支持 block entrance() ... end
	body := []ast.Stmt{
		&ast.VarDeclStmt{Name: "a", Expr: ast.NumberExpr{Value: 1}},
		&ast.VarDeclStmt{Name: "b", Expr: ast.NumberExpr{Value: 2}},
		&ast.PrintStmt{Expr: ast.VarExpr{Name: "a"}},
		&ast.PrintStmt{Expr: ast.VarExpr{Name: "b"}},
	}
	return &ast.FuncStmt{
		Name: "entrance",
		Body: body,
	}
}

// 将字符串转 uint64
func parseNumber(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
