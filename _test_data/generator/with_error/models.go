package with_error

import "github.com/shopspring/decimal"

type From struct {
	UUID string `map:"id"`
	Name string `map:"name"`
	Age  int64  `map:"age"`
}

type To struct {
	ID   decimal.Decimal `map:"id"`
	Name string          `map:"name"`
	Age  decimal.Decimal `map:"age"`
}
