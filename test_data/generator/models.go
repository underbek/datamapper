package generator

type Model struct {
	ID    int     `json:"id" map:"id"`
	Name  string  `json:"name" map:"name"`
	Age   float64 `json:"age" map:"age"`
	Empty string
}

type DAO struct {
	UUID string `db:"uuid" map:"id"`
	Name string `db:"name" map:"name"`
	Age  uint8  `json:"age" map:"age"`
}
