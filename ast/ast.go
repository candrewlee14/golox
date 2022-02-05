// package ast implements the Abstract Syntax Tree representation
package ast

import (
	"golox/token"
)

// AST Node
type Node interface {
	TokenLexeme() string
}

// Stmt node in the AST
type Stmt interface {
	Node
	statementNode()
}

// Expression node in the AST
type Expr interface {
	Node
	expressionNode()
}

type NumExpr struct {
	Token token.Token // NUMBER token
}

func (n *NumExpr) expressionNode() {}
func (n *NumExpr) TokenLexeme() string {
	return n.Token.Lexeme
}

type StrExpr struct {
	Token token.Token // STRING token
}

func (s *StrExpr) expressionNode() {}
func (s *StrExpr) TokenLexeme() string {
	return s.Token.Lexeme
}

// Program is a list of statements
type Program struct {
	Statements []Stmt
}

func (p *Program) TokenLexeme() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLexeme()
	} else {
		return ""
	}
}

// Var Statement in the form of 'var IDENT = EXPR'
type VarStmt struct {
	Token token.Token // VAR token
	Name  *Identifier
	Value Expr
}

func (vs *VarStmt) statementNode() {}
func (vs *VarStmt) TokenLexeme() string {
	return vs.Token.Lexeme
}

// Return statement in the form 'return EXPR'
type ReturnStmt struct {
	Token       token.Token // return token
	ReturnValue Expr
}

func (rs *ReturnStmt) statementNode() {}
func (rs *ReturnStmt) TokenLexeme() string {
	return rs.Token.Lexeme
}

// Identifier is a variable or function name
type Identifier struct {
	Token token.Token // token.IDENT
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLexeme() string {
	return i.Token.Lexeme
}
