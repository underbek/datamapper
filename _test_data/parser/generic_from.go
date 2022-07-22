package parser

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func ConvertAnyToString[A any](from A) string {
	return fmt.Sprint(from)
}

func ConvertIntUintToString[B int | uint](from B) string {
	return fmt.Sprint(from)
}

func ConvertIntegersToString[C int8 | int16 | int32](from C) string {
	return fmt.Sprint(from)
}

func ConvertXFloatToString[D constraints.Float](from D) string {
	return fmt.Sprint(from)
}
