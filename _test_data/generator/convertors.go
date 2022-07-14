package generator

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

func ConvertStringToDecimal(from string) decimal.Decimal {
	// TODO: handle error from conversion function
	to, _ := decimal.NewFromString(from)
	return to
}

func ConvertIntegerToDecimal[T constraints.Integer](from T) decimal.Decimal {
	return decimal.NewFromInt(int64(from))
}

func ConvertFloatToDecimal[T constraints.Float](from T) decimal.Decimal {
	return decimal.NewFromFloat(float64(from))
}
