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

func TestFuncDeclInvalid(t *testing.T) {
	progs := []string{
		`fun testFun(x,x) {return x;}`,
		`fun testFun( {return x;}`,
		`fun testFun) {return x;}`,
		`fun testFun(`,
		`fun testFun {return x;}`,
		`fun testFun(x) return x;}`,
		`fun testFun(x, y) return x;}`,
		`fun testFun(x) {return x}`,
		`fun testFun() {return x;`,
		`fun testFun() {return x;}}`,
		`fun testFun() return return x;}`,
		`fn testFun(x) {return x;}`,
		`fun 1() {return x;}`,
		`fun testFun(x, y, z,) {return x;`,
		`fun testFun(x, y, z,) {return x;`,
		`fun testFun(x, y z, a) {return x;}`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}

func TestFuncDeclValid(t *testing.T) {
	progs := []string{
		`fun testFun(x,y) {return x;}`,
		`fun testFun(x,) {return x;}`,
		`fun testFun(x, y, z,) {return x;}`,
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

func TestBlockValid(t *testing.T) {
	progs := []string{
		`{}`,
		`var x = 0;
        {while x != 10 {x = x + 1;}}`,
		`var x = 0;
        {while x != 10 {x = x + 1;}}
        {while x != 12 {x = x + 1; print x * 3;}}`,
		`{ var x = 0; }
        var x = 10;`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestBlockInvalid(t *testing.T) {
	progs := []string{
		`{}}`,
		`{{}`,
		`{};`,
		`{;};`,
		`var x = 0;
        {while x != 10 {x = x + 1; print x;}};`,
		`{{ var x = 0; }
        var x = 10;`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}

func TestCallExprValid(t *testing.T) {
	progs := []string{
		`testFun();`,
		`testFun(1, 2 + 3, x);`,
		`testFun(1, 2 + 3, x,);`,
		`testFun(testFun2());`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestCallExprInvalid(t *testing.T) {
	progs := []string{
		`10();`,
		`testFun(1,, 2 + 3, x);`,
		`testFun(1, 2 + 3, x,,);`,
		`testFun(10());`,
		`testFun(1 } 2);`,
		`testFun(1 2);`,
		`testFun(`,
		`testFun(;`,
		`testFun((;`,
		`testFun((;`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}

func TestVarValid(t *testing.T) {
	progs := []string{
		`var x = "hey";`,
		`var y = 13.5;`,
		`var x = 13.5 + 8 / 2 - 3 * 4;`,
		`var z = x;`,
		`var a = testFun(testFun2());`,
		`var x = testFun();
         var y = testFun();`,
		`var y = testFun(1, 2 + 3, x);`,
		`var z = testFun(1, 2 + 3, x,);`,
		`var a = testFun(testFun2());`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestVarInvalid(t *testing.T) {
	progs := []string{
		`var 10 = "hey";`,
		`var y 13.5;`,
		`var y { 13.5;`,
		`var y = 13.5`,
		`var z * x`,
		`var a testFun(testFun2());`,
		`var x = return testFun();`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}
func TestAssignValid(t *testing.T) {
	progs := []string{
		`x = "hey";`,
		`y = 13.5;`,
		`x = 13.5 + 8 / 2 - 3 * 4;`,
		`z = x;`,
		`a = testFun(testFun2());`,
		`x = testFun();
         y = testFun();`,
		`y = testFun(1, 2 + 3, x);`,
		`z = 1.5 + testFun(1, 2 + 3, x,);`,
		`a = testFun(testFun2());`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestAssignInvalid(t *testing.T) {
	progs := []string{
		`10 = "hey";`,
		`y = 13.5`,
		`y { 13.5;`,
		`y = 13.5 1;`,
		`y = z * x z`,
		`y = z * x z;`,
		`x = 10; y = z * x z`,
		`a b = testFun(testFun2());`,
		`x = return testFun();`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}
func TestReturnValid(t *testing.T) {
	progs := []string{
		`return "hey";`,
		`return 13.5 + 8 / 2 - 3 * 4;`,
		`return x;`,
		`return;`,
		`return testFun(testFun2());`,
		`return testFun();
         return testFun();`,
		`return testFun(1, 2 + 3, x);`,
		`return 1 + testFun(1, 2 + 3, x,);`,
		`return testFun(testFun2());`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestReturnInvalid(t *testing.T) {
	progs := []string{
		`return "hey"`,
		`return "hey; return"`,
		`return "hey; return;"`,
		`return "hey; return 19"`,
		`return "hey; return return"`,
		`return "hey; return return;"`,
		`return "hey; return return 19"`,
		`return "hey; return return 19;"`,
		`return 13.5 } 1`,
		`return {}`,
		`return fun hey(){}`,
		`return return testFun();`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}
func TestExprValid(t *testing.T) {
	progs := []string{
		`"hey";`,
		`13.5 + 8 / 2 - 3 * 4;`,
		`x;`,
		`x23;`,
		`testFun(testFun2());`,
		`testFun();
         testFun();`,
		`testFun(1, 2 + 3, x);`,
		`1 + testFun(1, 2 + 3, x,);`,
		`testFun(testFun2());`,
	}
	for _, progStr := range progs {
		assertNoErrors(t, progStr)
	}
}
func TestExprInvalid(t *testing.T) {
	progs := []string{
		`"hey"`,
		`x23 , 3`,
		`x23 . 3`,
	}
	for _, progStr := range progs {
		assertInvalid(t, progStr)
	}
}
