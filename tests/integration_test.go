//go:build integration
// +build integration

package tests

import (
	"golox/ast"
	"golox/interp"
	"golox/lexer"
	"golox/obj"
	"golox/parser"
	"testing"
)

func testExprNum(t *testing.T, node ast.Node, res float64) {
	intp := interp.New()
	val := intp.Eval(node)
	vb, ok := val.(*obj.Num)
	if !ok {
		t.Fatalf("Expected result of *obj.Num, got: %T", val)
	}
	if vb.Value != res {
		t.Fatalf("Expected num result to be %f, got: %f", res, vb.Value)
	}
}

// Integration Test
func TestCallExpr(t *testing.T) {
	input := `fun FunctionName(x,y,z) {
            var i = x * y * z;
            while i < 100 {
                i = i * 2;
            }
            return i;
        }
        FunctionName(10 - 3, 4 - 2, 5 - 3);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 112.0)
}

func TestFib(t *testing.T) {
	// this doesn't get fib(0) right
	input := `
        fun fib(n) {
            var pf = 0;
            var f = 1;
            var i = 0;
            while i < n - 1 {
                var temp = f;
                f = f + pf;
                pf = temp;
                i = i + 1;
            }
            return f;
        }
        fib(10);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 55.0)
}

func TestFibNestedReturn(t *testing.T) {
	// this should get fib(0) right
	input := `
        fun fib(n) {
            if n == 0 {
                return 0;
            }
            var pf = 0;
            var f = 1;
            var i = 0;
            while i < n - 1 {
                var temp = f;
                f = f + pf;
                pf = temp;
                i = i + 1;
            }
            return f;
        }
        fib(0);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 0)
}
