# datamapper

### Usage

### go generate
```go
package test

//go:generate datamapper --from Model --from-tag dto --to DTO --to-tag json -s models.go -d model_dto_converter.go -p test  
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

### Future

* [x] Parse and filter tag
* [x] Generate empty convertor
* [x] Map similar types
* [ ] Simple convertor test
* [ ] Create conversion functions
* [ ] Use conversion functions in convertor
* [ ] Parse user struct in struct
* [ ] Use generated conversion functions in convertor
* [ ] Parse embed struct
* [ ] Map filed without tag
* [ ] Parse other package
* [ ] Warning or error politics if tags is not equals


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