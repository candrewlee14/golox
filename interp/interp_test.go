//go:build unit
// +build unit

package interp

import (
	"golox/ast"
	"golox/obj"
	"golox/token"
	"testing"
)

func testExprBool(t *testing.T, node ast.Node, res bool) {
	intp := New()
	val := intp.Eval(node)
	vb, ok := val.(*obj.Bool)
	if !ok {
		t.Fatalf("Expected result of *obj.Bool, got: %T", val)
	}
	if vb.Value != res {
		t.Fatalf("Expected bool result to be %t, got: %t", res, vb.Value)
	}
}
func testExprNum(t *testing.T, node ast.Node, res float64) {
	intp := New()
	val := intp.Eval(node)
	vb, ok := val.(*obj.Num)
	if !ok {
		t.Fatalf("Expected result of *obj.Num, got: %T", val)
	}
	if vb.Value != res {
		t.Fatalf("Expected num result to be %f, got: %f", res, vb.Value)
	}
}
func testInfixNeqExpr(t *testing.T, node *ast.InfixExpr, res bool) {
	nodeNeq := node
	nodeNeq.Token = token.Token{Type: token.BANG_EQUAL, Lexeme: "!="}
	testExprBool(t, nodeNeq, res)
}

var trueTrueEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestTrueTrueEq(t *testing.T) {
	testExprBool(t, trueTrueEq, true)
	testInfixNeqExpr(t, trueTrueEq, false)
}

var falseTrueEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.BoolExpr{
		Token: token.Token{Type: token.FALSE, Lexeme: "false", Literal: false},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestFalseTrueEq(t *testing.T) {
	testExprBool(t, falseTrueEq, false)
	testInfixNeqExpr(t, falseTrueEq, true)
}

var strTrueEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeet", Literal: "yeet"},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestStrTrueEq(t *testing.T) {
	testExprBool(t, strTrueEq, false)
	testInfixNeqExpr(t, strTrueEq, true)
}

var trueStrEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeet", Literal: "yeet"},
	},
}

func TestTrueStrEq(t *testing.T) {
	testExprBool(t, trueStrEq, false)
	testInfixNeqExpr(t, trueStrEq, true)
}

var numTrueEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestNumTrueEq(t *testing.T) {
	testExprBool(t, numTrueEq, false)
	testInfixNeqExpr(t, numTrueEq, true)
}

var nilTrueEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.NilExpr{
		Token: token.Token{Type: token.NIL, Lexeme: "nil"},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestNilTrueEq(t *testing.T) {
	testExprBool(t, nilTrueEq, false)
	testInfixNeqExpr(t, nilTrueEq, true)
}

var nilNilEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.NilExpr{
		Token: token.Token{Type: token.NIL, Lexeme: "nil"},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.NilExpr{
		Token: token.Token{Type: token.NIL, Lexeme: "nil"},
	},
}

func TestNilNilTrueEq(t *testing.T) {
	testExprBool(t, nilNilEq, true)
	testInfixNeqExpr(t, nilNilEq, false)
}

var strStrMatchEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeet", Literal: "yeet"},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeet", Literal: "yeet"},
	},
}

func TestStrStrMatchEq(t *testing.T) {
	testExprBool(t, strStrMatchEq, true)
	testInfixNeqExpr(t, strStrMatchEq, false)
}

var strStrMismatchEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeetmismatch", Literal: "yeetmismatch"},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.StrExpr{
		Token: token.Token{Type: token.STRING, Lexeme: "yeet", Literal: "yeet"},
	},
}

func TestStrStrMismatchEq(t *testing.T) {
	testExprBool(t, strStrMismatchEq, false)
	testInfixNeqExpr(t, strStrMismatchEq, true)
}

var numNumMatchEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
}

func TestNumNumMatchEq(t *testing.T) {
	testExprBool(t, numNumMatchEq, true)
	testInfixNeqExpr(t, numNumMatchEq, false)
}

var numNumMismatchEq *ast.InfixExpr = &ast.InfixExpr{
	Left: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "13.6", Literal: 13.6},
	},
	Token: token.Token{Type: token.EQUAL_EQUAL, Lexeme: "=="},
	Right: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
}

func TestNumNumMismatchEq(t *testing.T) {
	testExprBool(t, numNumMismatchEq, false)
	testInfixNeqExpr(t, numNumMismatchEq, true)
}

func TestNumNumLess(t *testing.T) {
	numNumLess := numNumMismatchEq
	numNumLess.Token = token.Token{Type: token.LESS, Lexeme: ">"}
	testExprBool(t, numNumLess, true)
}
func TestNumNumLessEq(t *testing.T) {
	numNumLessEq := numNumMismatchEq
	numNumLessEq.Token = token.Token{Type: token.LESS_EQUAL, Lexeme: ">="}
	testExprBool(t, numNumLessEq, true)
}
func TestNumNumGreater(t *testing.T) {
	numNumGreater := numNumMismatchEq
	numNumGreater.Token = token.Token{Type: token.GREATER, Lexeme: "<"}
	testExprBool(t, numNumGreater, false)
}
func TestNumNumGreaterEq(t *testing.T) {
	numNumGreaterEq := numNumMismatchEq
	numNumGreaterEq.Token = token.Token{Type: token.GREATER_EQUAL, Lexeme: "<="}
	testExprBool(t, numNumGreaterEq, false)
}

func TestNumNumStar(t *testing.T) {
	numNumStar := numNumMismatchEq
	numNumStar.Token = token.Token{Type: token.STAR, Lexeme: "*"}
	testExprNum(t, numNumStar, 12.5*13.6)
}
func TestNumNumSlash(t *testing.T) {
	numNumSlash := numNumMismatchEq
	numNumSlash.Token = token.Token{Type: token.SLASH, Lexeme: "/"}
	testExprNum(t, numNumSlash, float64(13.6)/12.5)
}
func TestNumNumPlus(t *testing.T) {
	numNumPlus := numNumMismatchEq
	numNumPlus.Token = token.Token{Type: token.PLUS, Lexeme: "+"}
	testExprNum(t, numNumPlus, float64(12.5)+13.6)
}
func TestNumNumMinus(t *testing.T) {
	numNumMinus := numNumMismatchEq
	numNumMinus.Token = token.Token{Type: token.MINUS, Lexeme: "-"}
	testExprNum(t, numNumMinus, float64(13.6)-12.5)
}

var negNum *ast.PrefixExpr = &ast.PrefixExpr{
	Token: token.Token{Type: token.MINUS, Lexeme: "-"},
	Right: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
}

func TestNegNum(t *testing.T) {
	testExprNum(t, negNum, -float64(12.5))
}

var notNum *ast.PrefixExpr = &ast.PrefixExpr{
	Token: token.Token{Type: token.BANG, Lexeme: "!"},
	Right: ast.NumExpr{
		Token: token.Token{Type: token.NUMBER, Lexeme: "12.5", Literal: 12.5},
	},
}

func TestNotNum(t *testing.T) {
	testExprBool(t, notNum, false)
}

var notBool *ast.PrefixExpr = &ast.PrefixExpr{
	Token: token.Token{Type: token.BANG, Lexeme: "!"},
	Right: ast.BoolExpr{
		Token: token.Token{Type: token.TRUE, Lexeme: "true", Literal: true},
	},
}

func TestNotBool(t *testing.T) {
	testExprBool(t, notBool, false)
}

var notNil *ast.PrefixExpr = &ast.PrefixExpr{
	Token: token.Token{Type: token.BANG, Lexeme: "!"},
	Right: ast.NilExpr{
		Token: token.Token{Type: token.NIL, Lexeme: "nil"},
	},
}

func TestNotNil(t *testing.T) {
	testExprBool(t, notNil, true)
}
