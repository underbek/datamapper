package with_from_and_to_pointers_without_errors

type From struct {
	Name string `map:"name"`
	Age  int    `map:"age"`
}

type To struct {
	Name string `map:"name"`
	Age  uint   `map:"age"`
}
