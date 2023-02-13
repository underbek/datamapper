package to

type Order struct {
	ID         int          `recursive:"id"`
	User       User         `recursive:"user"`
	Operations []*Operation `recursive:"operations"`
}

type User struct {
	ID      string  `recursive:"id"`
	Account Account `recursive:"account"`
}

type Account struct {
	ID     int64  `recursive:"id"`
	Amount string `recursive:"amount"`
}

type Operation struct {
	ID     uint64 `recursive:"id"`
	Status string `recursive:"status"`
}
