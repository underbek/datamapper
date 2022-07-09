package parser

import (
	"golang.org/x/exp/constraints"
)

func ConvertStringToAny[T any](string) T {
	return *new(T)
}

func ConvertStringToIntUint[T int | uint](string) T {
	return T(25)
}

func ConvertStringToIntegers[T int8 | int16 | int32](string) T {
	return T(25)
}

func ConvertStringToXFloat[T constraints.Float](string) T {
	return T(25)
}
