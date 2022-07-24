package cf_with_pointers_and_erros

import "github.com/shopspring/decimal"

type From struct {
	ID    int  `map:"id"`
	Age   *int `map:"age"`
	Count *int `map:"count"`
	Orig  int  `map:"orig"`
}

type To struct {
	UUID  *decimal.Decimal `map:"id"`
	Age   decimal.Decimal  `map:"age"`
	Count *decimal.Decimal `map:"count"`
	Orig  decimal.Decimal  `map:"orig"`
}
