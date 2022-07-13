# datamapper

### Usage

### go generate

```go
package models

//go:generate datamapper --from Model --from-tag dto --to DTO --to-tag json -s models.go -d model_dto_converter.go -p models  
type Model struct {
	ID   int    `dto:"id" dao:"id"`
	Name string `dto:"name" dao:"name"`
}

type DTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
```

### Conversion functions

1. By types

```go
package conversion

func ConvertStringToStringPtr(from string) *string {
	return &from
}
```

2. By comments. Could use generic

```go
package conversion

import "fmt"

func ConvertAnyToString[T int | uint | float32](from T) string {
	return fmt.Sprint(from)
}

func ConvertStringToMany[T int | uint | float32](from int) T {
	return T(from)
}

func ConvertAnyToMany[T,V int | uint | float32](from T) V {
	return V(from)
}
```

### Future

* [x] Parse and filter tag
* [x] Generate empty convertor
* [x] Map similar types
* [x] Simple convertor test
* [x] Create conversion functions
* [x] Use conversion functions in convertor
* [x] Parse conversion functions from sources
* [x] Parse generic types from other package (constrains.Float)
* [x] Parse conversion functions with generic from
* [x] Parse conversion functions with generic to
* [x] Parse conversion functions with generic from and to
* [x] Parse conversion functions with generic struct
* [x] Parse conversion functions with struct
* [x] Create base conversation source
* [x] Generate convertors by other package models
* [ ] Generate convertors by other package fields in models
* [ ] Generate convertors by same package conversion functions
* [ ] Use other conversion functions in convertor
* [ ] Fix tests
* [ ] First generator tests
* [ ] Add CI with tests and linters
* [ ] Parse user struct in struct
* [ ] Use generated convertors in convertor like conversion function
* [ ] Parse embed struct
* [ ] Map filed without tag
* [ ] Parse other package
* [ ] Warning or error politics if tags is not equals
* [ ] Fill some conversion functions
* [ ] Copy using conversion functions from datamapper to target service if flag filled
* [ ] Use conversion functions with error -> convertor with error

### With comment ???

```go
package test

// DATAMAPPER convert from DTO:dto:json 
// DATAMAPPER convert to and from DAO:dao:db 
type Model struct {
	ID   int    `dto:"id" dao:"id"`
	Name string `dto:"name" dao:"name"`
}

type DAO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type DTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
```