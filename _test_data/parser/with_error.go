package parser

import (
	"strconv"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertStringToSigned[T constraints.Signed](from string) (T, error) {
	res, err := strconv.ParseInt(from, 10, 0)
	if err != nil {
		return 0, err
	}
	return T(res), nil
}

func ConvertStringToDecimal(from string) (decimal.Decimal, error) {
	return decimal.NewFromString(from)
}
