//go:build unit
// +build unit

package ast

import (
	"golox/token"
	"testing"
)

var program *Program = &Program{
	Statements: []Stmt{
		&VarStmt{
			Token: token.Token{Type: token.VAR, Lexeme: "var"},
			Name: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Lexeme: "myVar"},
			},
			Value: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Lexeme: "anotherVar"},
			},
		},
	},
}

func TestProgramString(t *testing.T) {
	s := "var myVar = anotherVar;"
	if program.String() != s {
		t.Fatalf("program.String() wrong. expected=%q, got=%q", s, program.String())
	}
}

func TestTokenLexeme(t *testing.T) {
	if program.TokenLexeme() != "var" {
		t.Fatalf("program.tokenLexeme() wrong. expected=%q, got=%q", "var",
			program.TokenLexeme())
	}
	stmt, ok := program.Statements[0].(*VarStmt)
	if !ok {
		t.Fatalf("program statement not *VarStmt. got=%T", program.Statements[0])
	}
	if stmt.TokenLexeme() != "var" {
		t.Fatalf("program.tokenLexeme() wrong. expected=%q, got=%q", "var",
			program.TokenLexeme())
	}
	if stmt.Name.TokenLexeme() != "myVar" {
		t.Fatalf("stmt.Name.tokenLexeme() wrong. expected=%q, got=%q", "myVar",
			stmt.Name.TokenLexeme())
	}
	if stmt.Value.TokenLexeme() != "anotherVar" {
		t.Fatalf("stmt.Name.tokenLexeme() wrong. expected=%q, got=%q", "var",
			stmt.Value.TokenLexeme())
	}
}
