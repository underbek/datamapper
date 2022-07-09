package parser

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func ConvertAnyToString[T any](from T) string {
	return fmt.Sprint(from)
}

func ConvertIntUintToString[T int | uint](from T) string {
	return fmt.Sprint(from)
}

func ConvertIntegersToString[T int8 | int16 | int32](from T) string {
	return fmt.Sprint(from)
}

func ConvertXFloatToString[T constraints.Float](from T) string {
	return fmt.Sprint(from)
}
