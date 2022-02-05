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

func TestString(t *testing.T) {
	expectedStr := "var myVar = anotherVar;"
	if program.String() != expectedStr {
		t.Fatalf("program.String() wrong. expected=%q, got=%q", expectedStr,
			program.String())
	}
	stmt, ok := program.Statements[0].(*VarStmt)
	if !ok {
		t.Fatalf("program statement not *VarStmt. got=%T", program.Statements[0])
	}
	if stmt.String() != expectedStr {
		t.Fatalf("stmt.String() wrong. expected=%q, got=%q", expectedStr,
			program.String())
	}
	if stmt.Name.String() != "myVar" {
		t.Fatalf("stmt.Name.String() wrong. expected=%q, got=%q", "myVar",
			stmt.Name.String())
	}
	if stmt.Value.String() != "anotherVar" {
		t.Fatalf("stmt.Value.String() wrong. expected=%q, got=%q", "var",
			stmt.Value.String())
	}
}

var returnProgram *Program = &Program{
	Statements: []Stmt{
		&ReturnStmt{
			Token: token.Token{Type: token.RETURN, Lexeme: "return"},
			ReturnValue: &BoolExpr{
				Token: token.Token{Type: token.TRUE, Lexeme: "true"},
			},
		},
		&ReturnStmt{
			Token: token.Token{Type: token.RETURN, Lexeme: "return"},
			ReturnValue: &StrExpr{
				Token: token.Token{Type: token.STRING, Lexeme: `"hey there"`},
			},
		},
		&ReturnStmt{
			Token: token.Token{Type: token.RETURN, Lexeme: "return"},
			ReturnValue: &NumExpr{
				Token: token.Token{Type: token.NUMBER, Lexeme: "1.345"},
			},
		},
	},
}

func TestReturnString(t *testing.T) {
	returnTests := []struct {
		statementStr string
		valueStr     string
	}{
		{"return true;", "true"},
		{`return "hey there";`, `"hey there"`},
		{`return 1.345;`, `1.345`},
	}
	if len(returnTests) != len(returnProgram.Statements) {
		t.Fatalf("Wrong number of statements. expected=%d, got=%d",
			len(returnTests),
			len(returnProgram.Statements))
	}
	for i, tt := range returnTests {
		stmt, ok := returnProgram.Statements[i].(*ReturnStmt)
		if !ok {
			t.Fatalf("program statement not *VarStmt. got=%T", program.Statements[i])
		}
		if stmt.String() != tt.statementStr {
			t.Fatalf("stmt.String() wrong. expected=%q, got=%q", tt.statementStr,
				stmt.String())
		}
		if stmt.ReturnValue.String() != tt.valueStr {
			t.Fatalf("stmt.ReturnValue.String() wrong. expected=%q, got=%q", tt.valueStr,
				stmt.ReturnValue.String())
		}
	}
}
