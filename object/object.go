package object

import (
	"bytes"
	"fmt"
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
