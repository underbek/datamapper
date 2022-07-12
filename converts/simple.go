package converts

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func ConvertNumericToString[T constraints.Integer | constraints.Float](from T) string {
	return fmt.Sprint(from)
}

func ConvertComplexToString[T constraints.Complex](from T) string {
	return fmt.Sprint(from)
}

func ConvertOrderedToOrdered[T, V constraints.Integer | constraints.Float](from T) V {
	return V(from)
}
