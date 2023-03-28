package user

import "time"

type User struct {
	ID        string `map:"user_id"`
	UserTimes *Times `map:"-"`
}

type Times struct {
	CreatedAt time.Time `map:"user_created_at"`
}
