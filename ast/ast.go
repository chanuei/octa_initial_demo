package ast

type Expr interface{}

type NumberExpr struct {
	Value int
}

type VarExpr struct {
	Name string
}

type Stmt interface{}

type VarDeclStmt struct {
	Name string
	Expr Expr
}

type AssignStmt struct {
	Name string
	Expr Expr
}

type PrintStmt struct {
	Expr Expr
}

type FuncStmt struct {
	Name string
	Body []Stmt
}
