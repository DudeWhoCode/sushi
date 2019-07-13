package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/dudewhocode/boa/ast"
)

type ObjectType string

const (
	INTEGEROBJ     = "INTEGER"
	BOOLEANOBJ     = "BOOLEAN"
	NULLOBJ        = "NULL"
	RETURNVALUEOBJ = "RETURN_VALUE"
	ERROROBJ       = "ERROR"
	FUNCTIONOBJ    = "FUNCTION"
	STRINGOBJ      = "STRING"
	BUILTINOBJ     = "BUILTIN"
	ARRAYOBJ       = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Null struct{}

type Error struct {
	Message string
}

type ReturnValue struct {
	Value Object
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type String struct {
	Value string
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

type Array struct {
	Elements []Object
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGEROBJ }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEANOBJ }

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULLOBJ }

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURNVALUEOBJ }

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROROBJ }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

func (f *Function) Type() ObjectType { return FUNCTIONOBJ }

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRINGOBJ }

func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Type() ObjectType { return BUILTINOBJ }

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (a *Array) Type() ObjectType { return ARRAYOBJ }

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
