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

var funcDeclProgram *Program = &Program{
	Statements: []Stmt{
		&FuncDeclStmt{
			Token: token.Token{Type: token.FUN, Lexeme: "fun"},
			Name: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Lexeme: "FunctionName"},
			},
			Params: []*Identifier{
				&Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Lexeme: "x"},
				},
				&Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Lexeme: "y"},
				},
				&Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Lexeme: "z"},
				},
			},
			Body: &BlockStmt{
				Token: token.Token{Type: token.LEFT_BRACE, Lexeme: "{"},
				Statements: []Stmt{
					&ReturnStmt{
						Token: token.Token{Type: token.RETURN, Lexeme: "return"},
						ReturnValue: &Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Lexeme: "x"},
						},
					},
				},
			},
		},
	},
}

func TestFuncDeclString(t *testing.T) {
	expectedStr := `fun FunctionName(x, y, z) {return x;}`
	if funcDeclProgram.String() != expectedStr {
		t.Fatalf("Function declaration statement String mismatch. Expected: %q, got=%q",
			expectedStr, funcDeclProgram)
	}
}

var exprProgram *Program = &Program{
	Statements: []Stmt{
		&ExprStmt{
			Token: token.Token{Type: token.IDENTIFIER, Lexeme: "testIdent"},
			Expr: &PrefixExpr{
				Token: token.Token{Type: token.MINUS, Lexeme: "-"},
				Right: &InfixExpr{
					Left: &Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Lexeme: "testIdent"},
					},
					Token: token.Token{Type: token.PLUS, Lexeme: "+"},
					Right: &NumExpr{
						Token: token.Token{Type: token.NUMBER, Lexeme: "10"},
					},
				},
			},
		},
	},
}

func TestExprStmtString(t *testing.T) {
	expectedStr := `(-(testIdent + 10))`
	if exprProgram.String() != expectedStr {
		t.Fatalf("Expression statement String mismatch. Expected: %q, got=%q",
			expectedStr, exprProgram)
	}
}

var ifProgram *Program = &Program{
	Statements: []Stmt{
		&IfStmt{
			Token: token.Token{Type: token.IF, Lexeme: "if"},
			Cond: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Lexeme: "condition"},
			},
			OnTrue: &BlockStmt{
				Token: token.Token{Type: token.LEFT_BRACE, Lexeme: "{"},
				Statements: []Stmt{
					&ReturnStmt{
						Token: token.Token{Type: token.RETURN, Lexeme: "return"},
						ReturnValue: &Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Lexeme: "x"},
						},
					},
				},
			},
			OnFalse: &BlockStmt{
				Token: token.Token{Type: token.LEFT_BRACE, Lexeme: "{"},
				Statements: []Stmt{
					&ReturnStmt{
						Token: token.Token{Type: token.RETURN, Lexeme: "return"},
						ReturnValue: &Identifier{
							Token: token.Token{Type: token.IDENTIFIER, Lexeme: "y"},
						},
					},
				},
			},
		},
	},
}

func TestIfStmtString(t *testing.T) {
	expectedStr := `if condition {return x;} else {return y;}`
	if ifProgram.String() != expectedStr {
		t.Fatalf("If statement String mismatch. Expected: %q, got=%q",
			expectedStr, ifProgram)
	}
}
