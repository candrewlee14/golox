// package ast implements the Abstract Syntax Tree representation
package ast

import (
	"bytes"
	"golox/token"
)

// AST Node
type Node interface {
	String() string
}

// Stmt node in the AST
type Stmt interface {
	Node
	statementNode()
}

// Program is a list of statements
type Program struct {
	Statements []Stmt
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Expression node in the AST
type Expr interface {
	Node
	expressionNode()
}

// Expression Statement
type ExprStmt struct {
	Token token.Token // first token of expression
	Expr  Expr
}

func (es ExprStmt) statementNode() {}
func (es ExprStmt) String() string {
	es.statementNode()
	return es.Expr.String()
}

// Expression Statement
type FuncDeclStmt struct {
	Token  token.Token // first token of expression
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStmt
}

func (fs FuncDeclStmt) statementNode() {}
func (fs FuncDeclStmt) String() string {
	fs.statementNode()
	var out bytes.Buffer
	out.WriteString("fun ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	for i, p := range fs.Params {
		out.WriteString(p.String())
		if i < len(fs.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(fs.Body.String())

	return out.String()
}

// Identifier is a variable or function name
type Identifier struct {
	Token token.Token // token.IDENT
}

func (i Identifier) expressionNode() {}
func (i Identifier) String() string {
	i.expressionNode() // in order to cover this empty private function in tests
	return i.Token.Lexeme
}

type NumExpr struct {
	Token token.Token // NUMBER token
}

func (n NumExpr) expressionNode() {}
func (n NumExpr) String() string {
	n.expressionNode()
	return n.Token.Lexeme
}

type NilExpr struct {
	Token token.Token // NIL token
}

func (n NilExpr) expressionNode() {}
func (n NilExpr) String() string {
	n.expressionNode()
	return n.Token.Lexeme
}

type StrExpr struct {
	Token token.Token // STRING token
}

func (s StrExpr) expressionNode() {}
func (s StrExpr) String() string {
	s.expressionNode()
	return s.Token.Lexeme
}

type BoolExpr struct {
	Token token.Token // TRUE or FALSE token
}

func (b BoolExpr) expressionNode() {}
func (b BoolExpr) String() string {
	b.expressionNode()
	return b.Token.Lexeme
}

type PrefixExpr struct {
	Token token.Token // unary operator token
	Right Expr
}

func (p PrefixExpr) expressionNode() {}
func (p PrefixExpr) String() string {
	p.expressionNode()
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Token.Lexeme) // operator
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type CallExpr struct {
	Token    token.Token // return token
	Function *Identifier
	Args     []Expr
}

func (ce CallExpr) expressionNode() {}
func (ce CallExpr) String() string {
	ce.expressionNode()
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, a := range ce.Args {
		out.WriteString(a.String())
		if i < len(ce.Args)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")

	return out.String()
}

type InfixExpr struct {
	Left  Expr
	Token token.Token // binary operator token
	Right Expr
}

func (ie InfixExpr) expressionNode() {}
func (ie InfixExpr) String() string {
	ie.expressionNode()
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Token.Lexeme) // operator
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// Var Statement in the form of 'var IDENT = EXPR'
type VarStmt struct {
	Token token.Token // VAR token
	Name  *Identifier
	Value Expr
}

func (vs VarStmt) statementNode() {}
func (vs VarStmt) String() string {
	vs.statementNode()
	var out bytes.Buffer

	out.WriteString(vs.Token.Lexeme + " ")
	out.WriteString(vs.Name.String())

	// TODO: remove nil check
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type BlockStmt struct {
	Token      token.Token // { token
	Statements []Stmt
}

func (bs *BlockStmt) statementNode() {}
func (bs BlockStmt) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

// If Statement in the form of 'if (COND) {OnTrue} else {OnFalse}'
type IfStmt struct {
	Token   token.Token // IF token
	Cond    Expr
	OnTrue  *BlockStmt
	OnFalse *BlockStmt
}

func (ifs IfStmt) statementNode() {}
func (ifs IfStmt) String() string {
	ifs.statementNode()
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ifs.Cond.String())
	out.WriteString(" ")
	out.WriteString(ifs.OnTrue.String())

	if ifs.OnFalse != nil {
		out.WriteString(" else ")
		out.WriteString(ifs.OnFalse.String())
	}
	return out.String()
}

// If Statement in the form of 'if (COND) {OnTrue} else {OnFalse}'
type WhileStmt struct {
	Token token.Token // IF token
	Cond  Expr
	Body  *BlockStmt
}

func (ws WhileStmt) statementNode() {}
func (ws WhileStmt) String() string {
	ws.statementNode()
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Cond.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}

// Return statement in the form 'return EXPR'
type ReturnStmt struct {
	Token       token.Token // return token
	ReturnValue Expr
}

func (rs ReturnStmt) statementNode() {}
func (rs ReturnStmt) String() string {
	rs.statementNode()
	var out bytes.Buffer
	out.WriteString(rs.Token.Lexeme)

	// TODO: remove nil check
	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}
