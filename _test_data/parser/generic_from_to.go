package parser

import (
	"golang.org/x/exp/constraints"
)

func ConvertXFloatToIntegers[T constraints.Float, V int | uint](from T) V {
	return V(from)
}
