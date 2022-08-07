package cf

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertDecimalToInt(from decimal.Decimal) int {
	return int(from.IntPart())
}

func ConvertIntegerToDecimal[T constraints.Integer](from T) decimal.Decimal {
	return decimal.NewFromInt(int64(from))
}

func ConvertFloatPtrToDecimalPtr[T constraints.Float](from *T) *decimal.Decimal {
	return nil
}
