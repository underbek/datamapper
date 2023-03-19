package cf

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertStringSliceToDecimalSlice(from []string) []decimal.Decimal {
	return nil
}

func ConvertIntegerToDecimal[T constraints.Integer](from T) decimal.Decimal {
	return decimal.NewFromInt(int64(from))
}
