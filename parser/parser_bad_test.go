//go:build unit
// +build unit

package parser

import (
	//"fmt"
	//"golox/ast"
	"golox/lexer"
    //"golox/token"
	"testing"
)

func assertInvalid(t *testing.T, source string) {
	l := lexer.NewLexer(source)
	p := New(&l)
	program := p.ParseProgram()
    if len(p.Errors()) == 0 {
        t.Fatalf("Expected errors in program. The parsed string: %q, resulting program: %q.",
			source,
			program)
    }
}
func assertNoErrors(t *testing.T, source string) {
	l := lexer.NewLexer(source)
	p := New(&l)
	p.ParseProgram()
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

func TestFuncDeclInvalid(t *testing.T){
    progs := []string{
        `fun testFun(x,x) {return x;}`,
        `fun testFun( {return x;}`,
        `fun testFun() return x;}`,
        `fun testFun() {return x}`,
        `fun testFun() {return x;`,
        `fun testFun() {return x;}}`,
        `fun testFun() return return x;}`,
    }
    for _, progStr := range progs {
        assertInvalid(t, progStr)
    }
}

func TestFuncDeclValid(t *testing.T){
    progs := []string{
        `fun testFun(x,y) {return x;}`,
        `fun testFun(x,) {return x;}`,
        `fun testFun(x) {return x;}`,
        `fun testFun() {return x;}`,
        `fun testFun() {
            var i = 0;
            var n = 2;
            while i < 10 {
                n = n * 2;
                i = i + 1;
            }
            return n;
        }`,
        `fun testFun() {
            var i = 0;
            var n = 2;
            while i < 10 {
                n = n * 2;
                i = i + 1;
            }
        }`,
    }
    for _, progStr := range progs {
        assertNoErrors(t, progStr)
    }
}
