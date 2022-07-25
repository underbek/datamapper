package parser

import (
	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertIntPtrToString(from *int) string {
	return fmt.Sprint(*from)
}

func ConvertFloatToStringPtr(float32) (res *string) {
	res2 := "test"
	return &res2
}

func ConvertFloatPtrToStringPtr(*float32) *string {
	return nil
}

func ConvertXFloatPointerToDecimal[T constraints.Float](*T) decimal.Decimal {
	return decimal.Zero
}
