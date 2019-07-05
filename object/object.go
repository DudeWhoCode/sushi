package object

import "fmt"

type ObjectType string

const (
	INTEGEROBJ = "INTEGER"
	BOOLEANOBJ = "BOOLEAN"
)

type Object interface {
	Type() string
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() string    { return INTEGEROBJ }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() string    { return BOOLEANOBJ }
