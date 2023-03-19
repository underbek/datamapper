package cf

import (
	"github.com/shopspring/decimal"
	"golang.org/x/exp/constraints"
)

func ConvertIntPtrToDecimalPtr(from *int) *decimal.Decimal {
	if from == nil {
		return nil
	}
	res := decimal.NewFromInt(int64(*from))
	return &res
}

func ConvertIntPtrToDecimal[T constraints.Signed](from *T) decimal.Decimal {
	if from == nil {
		return decimal.Zero
	}
	return decimal.NewFromInt(int64(*from))
}

func ConvertIntToDecimalPtr(from int) *decimal.Decimal {
	res := decimal.NewFromInt(int64(from))
	return &res
}

func ConvertIntToDecimal(from int) decimal.Decimal {
	return decimal.NewFromInt(int64(from))
}
