package same_cf_path

import "github.com/shopspring/decimal"

type From struct {
	ID   int     `map:"id"`
	Name string  `map:"name"`
	Age  float64 `map:"age"`
}

type To struct {
	UUID string          `map:"id"`
	Name string          `map:"name"`
	Age  decimal.Decimal `map:"age"`
}
