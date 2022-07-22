package tests

type TestModel struct {
	ID    int    `json:"id" map:"id"`
	Name  string `json:"name" map:"name"`
	Empty string
}

type TestModelTo struct {
	UUID string `db:"uuid" map:"id"`
	Name string `db:"name" map:"name"`
}

func (t *TestModelTo) foo() {
	panic("test")
}
