package from

import (
	"github.com/underbek/datamapper/_test_data/mapper/with_dash/simple/from/from_meta"
)

type UserData struct {
	User User           `map:"-"`
	Meta from_meta.Meta `map:"-"`
}

type User struct {
	ID   int    `map:"id"`
	Name string `map:"name"`
}
