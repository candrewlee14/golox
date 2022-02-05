//go:build unit
// +build unit

package parser

import (
	"fmt"
	"golox/ast"
	"golox/lexer"
	"golox/token"
	"testing"
)

type VarTest struct {
	expectedIdentifier  string
	expectedLine        int
	expectedLineOffset  int
	expectedExprType    token.TokenType
	expectedExprLiteral interface{}
}

func TestVarStmts(t *testing.T) {
	input := `
var x = 1.34;
var y = 2;
var foobar = 3814;
var str = "hey there";
`
	l := lexer.NewLexer(input)
	p := New(&l)

	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	tests := []VarTest{
		{"x", 1, 4, token.NUMBER, 1.34},
		{"y", 2, 4, token.NUMBER, 2.0},
		{"foobar", 3, 4, token.NUMBER, 3814.0},
		{"str", 4, 4, token.STRING, "hey there"},
	}
	if len(program.Statements) != len(tests) {
		t.Errorf("program.Statements: %s", program.Statements)
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			len(tests),
			len(program.Statements))
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		testVarStmt(t, stmt, tt)
	}
}

func testVarStmt(t *testing.T, s ast.Stmt, varTest VarTest) {
	if s.TokenLexeme() != "var" {
		t.Fatalf("s.TokenLexeme not 'var'. got=%q", s.TokenLexeme())
	}
	varStmt, ok := s.(*ast.VarStmt)
	if !ok {
		t.Fatalf("s not *ast.VarStatement. got=%q", s)
	}
	if varStmt.Name.Token.Lexeme != varTest.expectedIdentifier {
		t.Fatalf("s Name lexeme not '%s'. got=%s",
			varTest.expectedIdentifier,
			varStmt.Name.Token.Lexeme)
	}
	if varStmt.Name.Token.Line != varTest.expectedLine {
		t.Fatalf("varStmt.Name line not '%d'. got=%d",
			varTest.expectedLine,
			varStmt.Name.Token.Line)
	}
	if varStmt.Name.Token.LineOffset != varTest.expectedLineOffset {
		t.Fatalf("varStmt.Name line offset not '%d'. got=%d",
			varTest.expectedLineOffset,
			varStmt.Name.Token.LineOffset)
	}

	// TODO: check expr value literal equality
	// numExpr, okNum := varStmt.Value.(*ast.NumExpr)
	// strExpr, okStr := varStmt.Value.(*ast.StrExpr)
	// if okNum {
	//     if numExpr.Token.Literal != varTest.expectedExprLiteral {
	//     t.Fatalf("numExpr literal not '%s'. got=%s",
	//         fmt.Sprint(varTest.expectedExprLiteral),
	//         fmt.Sprint(varStmt.Name.Token.Literal))
	//     }
	// } else if okStr {
	//     if strExpr.Token.Literal != varTest.expectedExprLiteral {
	//     t.Fatalf("strExpr literal not '%s'. got=%s",
	//         fmt.Sprint(varTest.expectedExprLiteral),
	//         fmt.Sprint(varStmt.Name.Token.Literal))
	//     }
	// } else {
	//     t.Fatalf("varStmt.Value was not a numExpr or strExpr")
	// }
}

type ReturnTest struct {
	expectedLine        int
	expectedLineOffset  int
	expectedExprType    token.TokenType
	expectedExprLiteral interface{}
}

func TestReturnStmts(t *testing.T) {
	input := `
return 1.34;
return 2;
return 3814;
return "hey there";
`
	l := lexer.NewLexer(input)
	p := New(&l)

	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	tests := []ReturnTest{
		{1, 7, token.NUMBER, 1.34},
		{2, 7, token.NUMBER, 2.0},
		{3, 7, token.NUMBER, 3814.0},
		{4, 7, token.STRING, "hey there"},
	}
	if len(program.Statements) != len(tests) {
		t.Errorf("program.Statements: %s", program.Statements)
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			len(tests),
			len(program.Statements))
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		testReturnStmt(t, stmt, tt)
	}
}

func testReturnStmt(t *testing.T, stmt ast.Stmt, test ReturnTest) {
	returnStmt, ok := stmt.(*ast.ReturnStmt)
	if !ok {
		t.Errorf("stmt not *ast.returnStmt. got=%T", stmt)
	}
	if returnStmt.TokenLexeme() != "return" {
		t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
			returnStmt.TokenLexeme())
	}
	// TODO: check expr value literal equality
}

func assertNoParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentExpr(t *testing.T) {
	input := "return foobar;"
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	fmt.Println(program)
	assertNoParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. expected=1, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.ReturnValue.(ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.Identifier. got=%T",
			stmt.ReturnValue)
	}
	if ident.Token.Lexeme != "foobar" {
		t.Errorf("ident.Value is not %s. got=%s", "foobar", ident.Token.Lexeme)
	}
	if ident.TokenLexeme() != "foobar" {
		t.Errorf("ident.TokenLexeme() is not %s. got=%s", "foobar", ident.TokenLexeme())
	}
}
