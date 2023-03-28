package to

import (
	"github.com/shopspring/decimal"
	"github.com/underbek/datamapper/_test_data/mapper/with_dash/simple/to/to_meta"
)

type UserData struct {
	User User            `map:"-"`
	Age  decimal.Decimal `map:"age"`
}

type User struct {
	UUID     string       `map:"id"`
	Name     string       `map:"name"`
	MetaData to_meta.Meta `map:"-"`
}
