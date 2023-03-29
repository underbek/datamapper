package dao

import "time"

type Order struct {
	ID   int64  `db:"order_id"`
	UUID string `db:"order_uuid"`
}

type User struct {
	ID        int64     `db:"user_id"`
	CreatedAt time.Time `db:"user_created_at"`
}

type OrderUrls struct {
	OrderID     *int64 `db:"-"`
	SiteUrl     string `db:"url"`
	RedirectUrl string `db:"redirect_url"`
}

type OrderData struct {
	Order    *Order    `db:"-"`
	UserData *User     `db:"-"`
	Urls     OrderUrls `db:"-"`
}
