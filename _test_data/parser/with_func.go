package parser

import (
	"context"
)

type TestFunc func(ctx context.Context) error

type StructWithFunc struct {
	Func TestFunc
}
