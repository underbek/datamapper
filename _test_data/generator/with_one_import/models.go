package with_one_import

type From struct {
	ID   int    `map:"id"`
	Name string `map:"name"`
}

type To struct {
	UUID string `map:"id"`
	Name string `map:"name"`
}
