package domain

import "github.com/shopspring/decimal"

//go:generate ../../../bin/datamapper -i --from User --to User --to-source ../transport -d ../../local_test/domain_to_dto_user_converter.go --cf github.com/underbek/datamapper/_test_data/mapper/convertors --cf github.com/underbek/datamapper/_test_data/mapper/other_convertors
type User struct {
	ID         int             `map:"id"`
	Name       string          `map:"name"`
	Age        decimal.Decimal `map:"age"`
	ChildCount *int            `map:"children"`
	Empty      string
}
