package parser

type TestModel struct {
	ID    int    `json:"id,omitempty" map:"id"`
	Name  string `json:"name,omitempty" map:"name"`
	Empty string
}

type TestModelTo struct {
	UUID string `db:"uuid" map:"id"`
	Name string `db:"name" map:"name"`
}

type Model struct {
	ID string
}

func (t *TestModelTo) foo() {
	panic("test")
}
