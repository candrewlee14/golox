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
var super_cool_bool = true;
var _foo23 = false;
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
		{"super_cool_bool", 5, 4, token.TRUE, true},
		{"_foo23", 6, 4, token.FALSE, false},
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
	numExpr, okNum := varStmt.Value.(ast.NumExpr)
	strExpr, okStr := varStmt.Value.(ast.StrExpr)
	boolExpr, okBool := varStmt.Value.(ast.BoolExpr)
	if okNum {
		if numExpr.Token.Literal != varTest.expectedExprLiteral {
			t.Fatalf("numExpr literal not '%s'. got=%s",
				fmt.Sprint(varTest.expectedExprLiteral),
				fmt.Sprint(varStmt.Name.Token.Literal))
		}
	} else if okStr {
		if strExpr.Token.Literal != varTest.expectedExprLiteral {
			t.Fatalf("strExpr literal not '%s'. got=%s",
				fmt.Sprint(varTest.expectedExprLiteral),
				fmt.Sprint(varStmt.Name.Token.Literal))
		}
	} else if okBool {
		if boolExpr.Token.Literal != varTest.expectedExprLiteral {
			t.Fatalf("boolExpr literal not '%s'. got=%s",
				fmt.Sprint(varTest.expectedExprLiteral),
				fmt.Sprint(varStmt.Name.Token.Literal))
		}
	} else {
		t.Fatalf("varStmt.Value was not a numExpr or strExpr, got=%T", varStmt.Value)
	}
}

type ReturnTest struct {
	expectedLine        int
	expectedLineOffset  int
	expectedExprType    token.TokenType
	expectedExprString  string
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
		{1, 7, token.NUMBER, "1.34", 1.34},
		{2, 7, token.NUMBER, "2", 2.0},
		{3, 7, token.NUMBER, "3814", 3814.0},
		{4, 7, token.STRING, `"hey there"`, "hey there"},
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
		t.Fatalf("stmt not *ast.returnStmt. got=%T", stmt)
	}
	if returnStmt.ReturnValue.String() != test.expectedExprString {
		t.Fatalf("invalid return value string. expected=%q, got=%q",
			test.expectedExprString,
			returnStmt.ReturnValue.String())
	}
}

func assertNoParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}

func TestIdentExpr(t *testing.T) {
	input := "return foobar;"
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
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
}

func TestNumExpr(t *testing.T) {
	input := "return 1.513;"
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. expected=1, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStmt. got=%T",
			program.Statements[0])
	}
	num, ok := stmt.ReturnValue.(ast.NumExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.NumExpr. got=%T",
			stmt.ReturnValue)
	}
	if num.Token.Lexeme != "1.513" {
		t.Errorf("num lexeme is not %s. got=%s", "1.513", num.Token.Lexeme)
	}
	if num.Token.Literal != 1.513 {
		t.Errorf("num literal is not %f. got=%f", 1.513, num.Token.Literal)
	}
}

func TestStrExpr(t *testing.T) {
	valStr := "abc defg hey 12345"
	quotedStr := `"abc defg hey 12345"`
	input := `return "abc defg hey 12345";`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. expected=1, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStmt. got=%T",
			program.Statements[0])
	}
	str, ok := stmt.ReturnValue.(ast.StrExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.StrExpr. got=%T",
			stmt.ReturnValue)
	}
	if str.Token.Lexeme != quotedStr {
		t.Errorf("str lexeme is not %q. got=%q", quotedStr, str.Token.Lexeme)
	}
	if str.Token.Literal != valStr {
		t.Errorf("str literal is not %q. got=%q", valStr, str.Token.Literal)
	}
}

func TestBoolExpr(t *testing.T) {
	input := "return true;"
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. expected=1, got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStmt. got=%T",
			program.Statements[0])
	}
	b, ok := stmt.ReturnValue.(ast.BoolExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.BoolExpr. got=%T",
			stmt.ReturnValue)
	}
	if b.Token.Lexeme != "true" {
		t.Errorf("bool lexeme is not %s. got=%s", "true", b.Token.Lexeme)
	}
	if b.Token.Literal != true {
		t.Errorf("bool literal is not %t. got=%t", true, b.Token.Literal)
	}
}

func TestParsingPrefixExprs(t *testing.T) {
	prefixTests := []struct {
		input string
		op    string
		val   interface{}
	}{
		{"return !false;", "!", false},
		{"return -15;", "-", 15.0},
	}
	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := New(&l)
		program := p.ParseProgram()
		assertNoParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ReturnStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ReturnStmt. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.ReturnValue.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpr. got=%T", stmt.ReturnValue)
		}
		if exp.Token.Lexeme != tt.op {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.op, exp.Token.Lexeme)
		}
		switch exp.Right.(type) {
		case ast.NumExpr:
			numExp, _ := exp.Right.(ast.NumExpr)
			if numExp.Token.Literal != tt.val {
				t.Fatalf("numExp is not '%s'. got=%s",
					tt.val, numExp.Token.Literal)
			}
		case ast.BoolExpr:
			boolExp, _ := exp.Right.(ast.BoolExpr)
			if boolExp.Token.Literal != tt.val {
				t.Fatalf("boolExp is not '%s'. got=%s",
					tt.val, boolExp.Token.Literal)
			}
		default:
			t.Fatalf("exp.Right was not a num or bool")
		}
	}
}

type InfixTest struct {
	input    string
	leftVal  interface{}
	op       string
	rightVal interface{}
}

func testReturnInfixExpr(t *testing.T, tt InfixTest) {
	l := lexer.NewLexer(tt.input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ReturnStmt. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.ReturnValue.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("stmt is not ast.InfixExpr. got=%T", stmt.ReturnValue)
	}
	if exp.Token.Lexeme != tt.op {
		t.Fatalf("exp.Operator is not '%s'. got=%s",
			tt.op, exp.Token.Lexeme)
	}
	assertExprEq(t, tt.leftVal, exp.Left)
	assertExprEq(t, tt.rightVal, exp.Right)
}

func TestParsingInfixExprs(t *testing.T) {
	infixTests := []InfixTest{
		{"return 5 + 6;", 5.0, "+", 6.0},
		{"return 5 - 6;", 5.0, "-", 6.0},
		{"return 5 * 6;", 5.0, "*", 6.0},
		{"return 5 / 6;", 5.0, "/", 6.0},
		{"return 5 > 6;", 5.0, ">", 6.0},
		{"return 5 < 6;", 5.0, "<", 6.0},
		{"return 5 == 6;", 5.0, "==", 6.0},
		{"return 5 != 6;", 5.0, "!=", 6.0},
	}
	for _, tt := range infixTests {
		testReturnInfixExpr(t, tt)
	}
}

func assertExprEq(t *testing.T, expected interface{}, gotVal ast.Expr) {
	switch gotVal.(type) {
	case ast.NumExpr:
		numExp, _ := gotVal.(ast.NumExpr)
		if numExp.Token.Literal != expected {
			t.Fatalf("numExp is not '%s'. got=%s",
				expected, numExp.Token.Literal)
		}
	case ast.BoolExpr:
		boolExp, _ := gotVal.(ast.BoolExpr)
		if boolExp.Token.Literal != expected {
			t.Fatalf("boolExp is not '%s'. got=%s",
				expected, boolExp.Token.Literal)
		}
	default:
		t.Fatalf("Expression was not a num or bool")
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := New(&l)
		program := p.ParseProgram()
		assertNoParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, ident ast.Identifier, expectedStr string) {
	if ident.Token.Lexeme != expectedStr {
		t.Errorf("Identifier didn't match. Expected: %q, got: %q", expectedStr, ident.Token.Lexeme)
	}
}

func TestOperatorGroupParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4;",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5);",
			"(-(5 + 5))",
		},
		{
			"!(true == true);",
			"(!(true == true))",
		},
	}
	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := New(&l)
		program := p.ParseProgram()
		assertNoParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfStmt(t *testing.T) {
	input := `if (x < y) { z; } else { x; }`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStmt. got=%T",
			program.Statements[0])
	}

	if stmt.Cond.String() != "(x < y)" {
		t.Errorf("Condition mismatch. Expected=%q, got=%q\n",
			"(x < y)", stmt.Cond)
	}

	if len(stmt.OnTrue.Statements) != 1 {
		t.Errorf("OnTrue is not 1 statement. got=%d\n",
			len(stmt.OnTrue.Statements))
	}
	consequence, ok := stmt.OnTrue.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExprStmt. got=%T",
			stmt.OnTrue.Statements[0])
	}
	if consequence.Expr.String() != "z" {
		t.Fatalf("expected=x, got=%s", consequence.Expr)
	}

	alt, ok := stmt.OnFalse.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExprStmt. got=%T",
			stmt.OnFalse.Statements[0])
	}
	if alt.Expr.String() != "x" {
		t.Fatalf("expected=x, got=%s", alt.Expr)
	}
}

func TestWhileStmt(t *testing.T) {
	input := `while (x < y) { z; }`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStmt. got=%T",
			program.Statements[0])
	}
	if len(stmt.Body.Statements) != 1 {
		t.Errorf("Body is not 1 statement. got=%d\n",
			len(stmt.Body.Statements))
	}
	consequence, ok := stmt.Body.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExprStmt. got=%T",
			stmt.Body.Statements[0])
	}
	if consequence.Expr.String() != "z" {
		t.Fatalf("expected=x, got=%s", consequence.Expr)
	}
}

func TestFunDeclStmt(t *testing.T) {
	input := `fun FunctionName(x,y,z) {return x * y;}`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FuncDeclStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FuncDeclStmt. got=%T",
			program.Statements[0])
	}
	testIdentifier(t, *stmt.Name, "FunctionName")
	if len(stmt.Params) != 3 {
		t.Errorf("Expected 3 params. got=%d\n",
			len(stmt.Params))
	}
	testIdentifier(t, *stmt.Params[0], "x")
	testIdentifier(t, *stmt.Params[1], "y")
	testIdentifier(t, *stmt.Params[2], "z")
	// TODO: test body content
}

func TestFunCallStmt(t *testing.T) {
	input := `return FunctionName(x,10,z);`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ReturnStmt. got=%T",
			program.Statements[0])
	}
	funCall, ok := stmt.ReturnValue.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Return value is not *ast.CallExpr. got=%T",
			program.Statements[0])
	}
	if len(funCall.Args) != 3 {
		t.Errorf("Expected 3 args. got=%d\n",
			len(funCall.Args))
	}
	if funCall.Args[0].String() != "x" {
		t.Errorf("First param mismatch. Expected=%s, got=%s\n",
			"x", funCall.Args[0])
	}
	if funCall.Args[1].String() != "10" {
		t.Errorf("Second param mismatch. Expected=%s, got=%s\n",
			"10", funCall.Args[1])
	}
	if funCall.Args[2].String() != "z" {
		t.Errorf("Third param mismatch. Expected=%s, got=%s\n",
			"z", funCall.Args[2])
	}
	// TODO: test arg types
}

func TestAssignStmt(t *testing.T) {
	input := `x = 100;`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignStmt. got=%T",
			program.Statements[0])
	}
	if stmt.Name.String() != "x" {
		t.Fatalf("Assignment name mismatch. Expected=%s, got=%T",
			"x", stmt.Name)
	}
	fmt.Printf("%T", stmt.Expr)
	num, ok := stmt.Expr.(ast.NumExpr)
	if !ok {
		t.Fatalf("Assignment value is not *ast.NumExpr. got=%T",
			stmt.Expr)
	}
	if num.String() != "100" {
		t.Fatalf("Assignment value mismatch. Expected=%s, got=%T",
			"100", num.String())
	}
}

func TestAssignNilStmt(t *testing.T) {
	input := `x = nil;`
	l := lexer.NewLexer(input)
	p := New(&l)
	program := p.ParseProgram()
	assertNoParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement, got=%d\n",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignStmt. got=%T",
			program.Statements[0])
	}
	if stmt.Name.String() != "x" {
		t.Fatalf("Assignment name mismatch. Expected=%s, got=%T",
			"x", stmt.Name)
	}
	fmt.Printf("%T", stmt.Expr)
	_, isNil := stmt.Expr.(ast.NilExpr)
	if !isNil {
		t.Fatalf("Assignment value is not ast.NilExpr. got=%T",
			stmt.Expr)
	}
}
