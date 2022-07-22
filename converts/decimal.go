package converts

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertDecimalToString(from decimal.Decimal) string {
	return from.String()
}

func ConvertDecimalToNumeric[T constraints.Integer | constraints.Float](from decimal.Decimal) T {
	return T(from.InexactFloat64())
}

func ConvertStringToDecimal(from string) (decimal.Decimal, error) {
	return decimal.NewFromString(from)
}

func ConvertIntegerToDecimal[T constraints.Integer](from T) decimal.Decimal {
	return decimal.NewFromInt(int64(from))
}

func ConvertFloatToDecimal[T constraints.Float](from T) decimal.Decimal {
	return decimal.NewFromFloat(float64(from))
}
