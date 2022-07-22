package same_cf_path

import (
	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertNumericToString[T constraints.Integer | constraints.Float](from T) string {
	return fmt.Sprint(from)
}

func ConvertFloatToDecimal[T constraints.Float](from T) decimal.Decimal {
	return decimal.NewFromFloat(float64(from))
}
