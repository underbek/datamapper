package from

import "github.com/shopspring/decimal"

type Order struct {
	ID         int64       `recursive:"id"`
	User       User        `recursive:"user"`
	Operations []Operation `recursive:"operations"`
}

type User struct {
	ID      int64    `recursive:"id"`
	Account *Account `recursive:"account"`
}

type Account struct {
	ID     int64           `recursive:"id"`
	Amount decimal.Decimal `recursive:"amount"`
}

type Operation struct {
	ID     int64  `recursive:"id"`
	Status string `recursive:"status"`
}
