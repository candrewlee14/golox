// package ast implements the Abstract Syntax Tree representation
package ast

import (
	"bytes"
	//"fmt"
	"golox/token"
)

// AST Node
type Node interface {
	TokenLexeme() string
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

func (p Program) TokenLexeme() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLexeme()
	} else {
		return ""
	}
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
func (i Identifier) TokenLexeme() string {
	return i.Token.Lexeme
}
func (i Identifier) String() string {
	return i.TokenLexeme()
}

// type ExprStatement struct {
// 	Token token.Token
// 	Expr  Expr
// }

// func (es *ExprStatement) statementNode() {}
// func (es *ExprStatement) TokenLexeme() string {
// 	return es.Token.Lexeme
// }
// func (es *ExprStatement) String() string {
// 	// TODO: remove nil check
// 	if es.Expr != nil {
// 		return es.Expr.String()
// 	}
// 	return ""
// }

type NumExpr struct {
	Token token.Token // NUMBER token
}

func (n NumExpr) expressionNode() {}
func (n NumExpr) TokenLexeme() string {
	return n.Token.Lexeme
}

type StrExpr struct {
	Token token.Token // STRING token
}

func (s StrExpr) expressionNode() {}
func (s StrExpr) TokenLexeme() string {
	return s.Token.Lexeme
}

// Var Statement in the form of 'var IDENT = EXPR'
type VarStmt struct {
	Token token.Token // VAR token
	Name  *Identifier
	Value Expr
}

func (vs VarStmt) statementNode() {}
func (vs VarStmt) TokenLexeme() string {
	return vs.Token.Lexeme
}
func (vs VarStmt) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLexeme() + " ")
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
func (rs ReturnStmt) TokenLexeme() string {
	return rs.Token.Lexeme
}
func (rs ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLexeme())

	// TODO: remove nil check
	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}
