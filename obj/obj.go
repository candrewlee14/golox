package obj

import "fmt"

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
)

type Nil struct{}

func (n *Nil) Type() ObjType  { return NIL_OBJ }
func (n *Nil) String() string { return "nil" }

type Num struct {
	Value float64
}

func (n *Num) Type() ObjType  { return NUM_OBJ }
func (n *Num) String() string { return fmt.Sprint(n.Value) }

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
