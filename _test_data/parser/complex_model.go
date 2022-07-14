package parser

import "github.com/shopspring/decimal"

type ComplexModel struct {
	ID  Model           `json:"id" map:"id"`
	Age decimal.Decimal `json:"age" map:"age"`
}
