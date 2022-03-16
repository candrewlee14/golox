//go:build integration
// +build integration

package interp

import (
	"golox/lexer"
	"golox/parser"
	"testing"
)

// Integration Test
func TestCallExpr(t *testing.T) {
	input := `fun FunctionName(x,y,z) {
            var i = x * y * z;
            while i < 100 {
                i = i * 2;
            }
            return i;
        }
        return FunctionName(10 - 3, 4 - 2, 5 - 3);`
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
        return fib(10);`
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
        return fib(0);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 0)
}

func TestFuncScope(t *testing.T) {
	input := `
        fun testFun() {
            return x;
        }
        var x = 100 + 3;
        return testFun();`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	intp := New()
	val := intp.Eval(program) // This should error for this function
	intp.PrintEnv()
	t.Fatalf("Program should not have x in scope, should've been runtime error. Instead returned %q", val)
}

func TestFuncScopeModification(t *testing.T) {
	input := `
        var x = 100 + 3;
        fun testFun() {
            return x;
        }
        x = 10;
        return testFun();`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 10.0)
}

func TestFuncScopeParam(t *testing.T) {
	input := `
        var x = 100 + 3;
        fun testFun(x) {
            return x;
        }
        x = 10;
        return testFun(5.0);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 5.0)
}

func TestNestedIfReturn(t *testing.T) {
	input := `
        var x = 100;
        fun clamp(min, x, max) {
            if x < min {
                return min;
            }
            if x > max {
                return max;
            }
            return x;
        }
        x = 5;
        return clamp(-1, -135, 100) + clamp(-50, 50, 100) + clamp(0, 560, 126);`
	// -1 + 50 + 126 = 175
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 175.0)
}

func TestRecursion(t *testing.T) {
	input := `
        fun testFib(n) {
            if n < 1 { return 0; }
            if n == 1 { return 1; }
            return testFib(n - 1) + testFib(n - 2);
        }
        return testFib(19) + testFib(0);`
	l := lexer.NewLexer(input)
	p := parser.New(&l)
	program := p.ParseProgram()
	testExprNum(t, program, 4181)
}
