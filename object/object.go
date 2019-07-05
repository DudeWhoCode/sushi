package object

import "fmt"

type ObjectType string

const (
	INTEGEROBJ = "INTEGER"
)

type Object interface {
	Type() string
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() string    { return INTEGEROBJ }
