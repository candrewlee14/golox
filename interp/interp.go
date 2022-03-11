package interp

import (
	"fmt"
	"golox/ast"
	"golox/obj"
	"golox/token"
)

func Eval(node ast.Node) obj.Obj {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStmts(node.Statements)
	case *ast.ExprStmt:
		return Eval(node.Expr)
	case *ast.ReturnStmt:
		return Eval(node.ReturnValue)
	case *ast.VarStmt:
		return Eval(node.Value)
	// Expressions
	case ast.NumExpr:
		fl, _ := node.Token.Literal.(float64)
		return &obj.Num{Value: fl}
	case ast.NilExpr:
		return &obj.Nil{}
	case ast.StrExpr:
		str, _ := node.Token.Literal.(string)
		return &obj.Str{Value: str}
	case ast.BoolExpr:
		b, _ := node.Token.Literal.(bool)
		return &obj.Bool{Value: b}
	case *ast.PrefixExpr:
		return evalPrefix(node)
	case *ast.InfixExpr:
		return evalInfix(node)
	}
	panic(fmt.Sprintf("Unable to evaluate unexpected expression, got: %T", node))
}

func evalStmts(stmts []ast.Stmt) obj.Obj {
	var result obj.Obj
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefix(pe *ast.PrefixExpr) obj.Obj {
	r := Eval(pe.Right)
	switch pe.Token.Type {
	case token.BANG:
		return &obj.Bool{Value: !isTruthy(r)}
	case token.MINUS:
		n := resolveNum(r)
		return &obj.Num{Value: -n.Value}
	}
	panic(fmt.Sprintf("Expected prefix operator, got: %s", pe.Token.Type))
}

func evalInfix(ie *ast.InfixExpr) obj.Obj {
	l := Eval(ie.Left)
	r := Eval(ie.Right)
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
		fmt.Println("HEY")
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
