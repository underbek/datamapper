package parser

import "github.com/underbek/datamapper/_test_data/parser/other"

type ModelWithDash struct {
	User  DashUser `map:"-"`
	Count int      `map:"count"`
}

type DashUser struct {
	ID   int                `map:"id"`
	Name string             `map:"name"`
	Meta other.DashUserMeta `map:"-"`
}
