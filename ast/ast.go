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
