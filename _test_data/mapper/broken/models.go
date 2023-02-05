package transport

import "github.com/google/uuid"

type User struct {
	UUID       uuid.UUID `json:"uuid" map:"id"`
	Name       string    `json:"name" map:"name"`
	Age        string    `map:"age"`
	ChildCount *string   `map:"children"`
}

type BrokenUser [T]struct {
	UUID       uuid.UUID `json:"uuid" map:"id"`
	Name       string    `json:"name" map:"name"`
	Age        string    `map:"age"`
	ChildCount *string   `map:"children"`
}

type BrokenUser2 struct {
	UUID `json:"uuid"`
}
