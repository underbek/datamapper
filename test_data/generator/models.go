package generator

type Model struct {
	ID    int    `json:"id" map:"id"`
	Name  string `json:"name" map:"name"`
	Empty string
}

type DAO struct {
	UUID string `db:"uuid" map:"id"`
	Name string `db:"name" map:"name"`
}
