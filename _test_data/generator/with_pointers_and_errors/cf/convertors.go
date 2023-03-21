package cf

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertIntPtrToDecimalPtr(from *int) (*decimal.Decimal, error) {
	if from == nil {
		return nil, nil
	}
	res := decimal.NewFromInt(int64(*from))
	return &res, nil
}

func ConvertIntPtrToDecimal[T constraints.Signed](from *T) (decimal.Decimal, error) {
	if from == nil {
		return decimal.Zero, nil
	}
	return decimal.NewFromInt(int64(*from)), nil
}

func ConvertIntToDecimalPtr(from int) (*decimal.Decimal, error) {
	res := decimal.NewFromInt(int64(from))
	return &res, nil
}

func ConvertIntToDecimal(from int) (decimal.Decimal, error) {
	return decimal.NewFromInt(int64(from)), nil
}
