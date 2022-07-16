package parser

import "github.com/shopspring/decimal"

type PointerModel struct {
	ID   *int             `map:"id"`
	Name *string          `map:"name"`
	Age  *decimal.Decimal `map:"age"`
}
