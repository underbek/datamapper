package cf

import (
	"github.com/shopspring/decimal"
)

func ConvertStringToDecimal(from string) (decimal.Decimal, error) {
	return decimal.NewFromString(from)
}
