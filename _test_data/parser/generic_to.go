package parser

import (
	"golang.org/x/exp/constraints"
)

func ConvertStringToAny[A any](string) A {
	return *new(A)
}

func ConvertStringToIntUint[B int | uint](string) B {
	return B(25)
}

func ConvertStringToIntegers[C int8 | int16 | int32](string) C {
	return C(25)
}

func ConvertStringToXFloat[D constraints.Float](string) D {
	return D(25)
}
