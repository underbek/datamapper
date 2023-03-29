package domain

import (
	"github.com/underbek/datamapper/_test_data/mapper/with_dash_and_pointers/domain/user"
)

type Order struct {
	OrderID     *string    `map:"order_id"`
	OrderUUID   string     `map:"order_uuid"`
	User        *user.User `map:"-"`
	SiteUrl     string     `map:"url"`
	RedirectUrl string     `map:"redirect_url"`
}
