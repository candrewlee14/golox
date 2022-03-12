package obj

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"golox/ast"
)

type Env struct {
	Bindings map[string]*Box
}

type Box struct {
	Ref *Obj
}

func (e *Env) PrintColored() {
	for key, elem := range e.Bindings {
		fmt.Println(color.CyanString("%s", key), "=", color.YellowString("boxed"), *elem.Ref)
	}
}

func NewEnv() Env {
	bindings := make(map[string]*Box)
	return Env{Bindings: bindings}
}

func (e *Env) Bind(name string, val Obj) {
	_, bound := e.Bindings[name]
	if !bound {
		e.Bindings[name] = &Box{&val}
	} else {
		panic(fmt.Sprintf("Variable %q already exists in this scope. Use \"%s = ...;\" to assign instead.", name, name))
	}
}

type ObjType uint8

type Obj interface {
	Type() ObjType
	String() string
}

const (
	NIL_OBJ ObjType = iota
	NUM_OBJ
	BOOL_OBJ
	STR_OBJ
	CLOSURE_OBJ
	RET_VAL_OBJ
)

type Nil struct{}

func (n *Nil) Type() ObjType  { return NIL_OBJ }
func (n *Nil) String() string { return "nil" }

type Closure struct {
	EnvStack []Env
	Params   []*ast.Identifier
	Body     *ast.BlockStmt
}

func (c *Closure) Type() ObjType { return CLOSURE_OBJ }
func (c *Closure) String() string {
	var out bytes.Buffer
	out.WriteString("fun(")
	for i, p := range c.Params {
		out.WriteString(p.String())
		if i < len(c.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(c.Body.String())
	return out.String()
}

type RetVal struct {
	Val Obj
}

func (rv *RetVal) Type() ObjType  { return RET_VAL_OBJ }
func (rv *RetVal) String() string { return "ret " + fmt.Sprint(rv.Val) }

type Num struct {
	Value float64
}

func (n *Num) Type() ObjType  { return NUM_OBJ }
func (n *Num) String() string { return fmt.Sprintf("%f", n.Value) }

type Bool struct {
	Value bool
}

func (b *Bool) Type() ObjType  { return BOOL_OBJ }
func (b *Bool) String() string { return fmt.Sprint(b.Value) }

type Str struct {
	Value string
}

func (s *Str) Type() ObjType  { return STR_OBJ }
func (s *Str) String() string { return s.Value }
