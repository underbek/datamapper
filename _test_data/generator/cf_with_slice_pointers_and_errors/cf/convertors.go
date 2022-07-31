package cf

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertDecimalToInt(from decimal.Decimal) (int, error) {
	return int(from.IntPart()), nil
}

func ConvertIntegerToDecimal[T constraints.Integer](from T) (decimal.Decimal, error) {
	return decimal.NewFromInt(int64(from)), nil
}

func ConvertFloatPtrToDecimalPtr[T constraints.Float](from *T) (*decimal.Decimal, error) {
	return nil, nil
}
