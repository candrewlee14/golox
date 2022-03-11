package interp

import (
	"fmt"
	"github.com/fatih/color"
	"golox/ast"
	"golox/obj"
	"golox/token"
)

type Env struct {
	Bindings map[string]obj.Obj
}

func (e *Env) PrintColored() {
	for key, elem := range e.Bindings {
		fmt.Println(color.CyanString("%s", key), "=", elem)
	}
}

func newEnv() Env {
	bindings := make(map[string]obj.Obj)
	return Env{Bindings: bindings}
}

type Interpreter struct {
	EnvStack []Env
}

func New() Interpreter {
	return Interpreter{EnvStack: []Env{newEnv()}}
}
func (intp *Interpreter) PrintEnv() {
	i := len(intp.EnvStack) - 1
	for i >= 0 {
		fmt.Println("-----")
		intp.EnvStack[i].PrintColored()
		i--
	}
}

func (intp *Interpreter) bind(name string, val obj.Obj) {
	intp.EnvStack[len(intp.EnvStack)-1].bind(name, val)
}
func (intp *Interpreter) assign(name string, val obj.Obj) {
	i := len(intp.EnvStack) - 1
	for i >= 0 {
		_, ok := intp.EnvStack[i].Bindings[name]
		if ok {
			intp.EnvStack[i].Bindings[name] = val
			return
		}
		i--
	}
	panic(fmt.Sprintf("Variable %q does not exist in this scope. Use \"var %s = ...;\" to declare instead.", name, name))
}

func (intp *Interpreter) resolve(name string) (obj.Obj, bool) {
	i := len(intp.EnvStack) - 1
	for i >= 0 {
		val, ok := intp.EnvStack[i].Bindings[name]
		if ok {
			return val, true
		}
		i--
	}
	return nil, false
}

func (e *Env) bind(name string, val obj.Obj) {
	_, bound := e.Bindings[name]
	if !bound {
		e.Bindings[name] = val
	} else {
		panic(fmt.Sprintf("Variable %q already exists in this scope. Use \"%s = ...;\" to assign instead.", name, name))
	}
}

func (intp *Interpreter) Eval(node ast.Node) obj.Obj {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return intp.evalStmts(node.Statements)
	case *ast.ExprStmt:
		return intp.Eval(node.Expr)
	case *ast.ReturnStmt:
		return intp.Eval(node.ReturnValue)
	case *ast.BlockStmt:
		intp.EnvStack = append(intp.EnvStack, newEnv())
		defer func() { intp.EnvStack = intp.EnvStack[:len(intp.EnvStack)-1] }()
		val := intp.evalStmts(node.Statements)
		return val
	case *ast.AssignStmt:
		intp.assign(node.Name.String(), intp.Eval(node.Expr))
		return nil
	case *ast.WhileStmt:
		// TODO: work out what should be truthy here
		for isTruthy(intp.Eval(node.Cond)) {
			_ = intp.Eval(node.Body)
		}
		return nil
	case *ast.VarStmt:
		val := intp.Eval(node.Value)
		intp.bind(node.Name.String(), val)
		return &obj.Nil{}
	case *ast.IfStmt:
		cond := intp.Eval(node.Cond)
		if isTruthy(cond) {
			return intp.Eval(node.OnTrue)
		} else if node.OnFalse != nil {
			return intp.Eval(node.OnFalse)
		} else {
			return nil
		}
	// Expressions
	case ast.NumExpr:
		fl, _ := node.Token.Literal.(float64)
		return &obj.Num{Value: fl}
	case ast.Identifier:
		val, found := intp.resolve(node.Token.Lexeme)
		if !found {
			panic(fmt.Sprintf("Variable %q is not declared in this scope.", node.Token.Lexeme))
		}
		return val
	case ast.NilExpr:
		return &obj.Nil{}
	case ast.StrExpr:
		str, _ := node.Token.Literal.(string)
		return &obj.Str{Value: str}
	case ast.BoolExpr:
		b, _ := node.Token.Literal.(bool)
		return &obj.Bool{Value: b}
	case *ast.PrefixExpr:
		return intp.evalPrefix(node)
	case *ast.InfixExpr:
		return intp.evalInfix(node)
	}
	panic(fmt.Sprintf("Unable to evaluate unexpected expression, got: %T", node))
}

func (intp *Interpreter) evalStmts(stmts []ast.Stmt) obj.Obj {
	var result obj.Obj
	for _, stmt := range stmts {
		result = intp.Eval(stmt)
	}
	return result
}

func (intp *Interpreter) evalPrefix(pe *ast.PrefixExpr) obj.Obj {
	r := intp.Eval(pe.Right)
	switch pe.Token.Type {
	case token.BANG:
		return &obj.Bool{Value: !isTruthy(r)}
	case token.MINUS:
		n := resolveNum(r)
		return &obj.Num{Value: -n.Value}
	}
	panic(fmt.Sprintf("Expected prefix operator, got: %s", pe.Token.Type))
}

func (intp *Interpreter) evalInfix(ie *ast.InfixExpr) obj.Obj {
	l := intp.Eval(ie.Left)
	r := intp.Eval(ie.Right)
	switch ie.Token.Type {
	case token.PLUS:
		return &obj.Num{Value: resolveNum(l).Value + resolveNum(r).Value}
	case token.MINUS:
		return &obj.Num{Value: resolveNum(l).Value - resolveNum(r).Value}
	case token.STAR:
		return &obj.Num{Value: resolveNum(l).Value * resolveNum(r).Value}
	case token.SLASH:
		return &obj.Num{Value: resolveNum(l).Value / resolveNum(r).Value}
	case token.LESS:
		return &obj.Bool{Value: resolveNum(l).Value > resolveNum(r).Value}
	case token.LESS_EQUAL:
		return &obj.Bool{Value: resolveNum(l).Value >= resolveNum(r).Value}
	case token.GREATER:
		return &obj.Bool{Value: resolveNum(l).Value < resolveNum(r).Value}
	case token.GREATER_EQUAL:
		return &obj.Bool{Value: resolveNum(l).Value <= resolveNum(r).Value}
	case token.EQUAL_EQUAL:
		return &obj.Bool{Value: isEq(l, r)}
	case token.BANG_EQUAL:
		return &obj.Bool{Value: !isEq(l, r)}
	}
	panic(fmt.Sprintf("Expected infix operator, got: %s", ie.Token.Type))
}

func resolveNum(o obj.Obj) *obj.Num {
	// TODO: resolve variables to nums
	switch o := o.(type) {
	case *obj.Num:
		return o
	}
	panic(fmt.Sprintf("Unable to resolve object to number. Expected: *obj.Num, got: %T", o))
}

func isEq(a obj.Obj, b obj.Obj) bool {
	// TODO: resolve variables
	switch a := a.(type) {
	case *obj.Bool:
		bb, bIsBool := b.(*obj.Bool)
		if !bIsBool {
			return false
		}
		return a.Value == bb.Value
	case *obj.Num:
		bn, bIsNum := b.(*obj.Num)
		if !bIsNum {
			return false
		}
		return a.Value == bn.Value
	case *obj.Str:
		bs, bIsStr := b.(*obj.Str)
		if !bIsStr {
			return false
		}
		return a.Value == bs.Value
	case *obj.Nil:
		_, bIsNil := b.(*obj.Nil)
		if !bIsNil {
			return false
		}
		return true
	}
	panic(fmt.Sprintf("Unable to compare objects. Got: %T and %T", a, b))
}

func isTruthy(o obj.Obj) bool {
	// TODO: resolve variables to value
	switch o := o.(type) {
	case *obj.Bool:
		return o.Value
	case *obj.Nil:
		return false
	default:
		return true
	}
}
