# datamapper

### Install

```shell
go install github.com/underbek/datamapper@v0.0.4
```

### Usage

```text
Usage:
  datamapper [OPTIONS]

Application Options:
  -d, --destination= Destination file path
      --cf=          User conversion functions sources/packages. Can add package alias like {package_path}:{alias)
      --from=        Model from name
      --from-tag=    Model from tag (default: map)
      --from-source= From model source/package. Can add package alias like {package_path}:{alias) (default: .)
      --to=          Model to name
      --to-tag=      Model to tag (default: map)
      --to-source=   To model source/package. Can add package alias like {package_path}:{alias) (default: .)
  -i, --inverse      Create direct and inverse conversions

Help Options:
  -h, --help         Show this help message
```

### go generate

```go
package models

//go:generate datamapper --from Model --from-tag dto --to DTO --to-tag json -d model_dto_converter.go
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

import "fmt"

func ConvertIntToString(from int) string {
	return fmt.Sprint(from)
}
```

2. By generic

```go
package conversion

import "fmt"

func ConvertAnyToString[T int | uint | float32](from T) string {
	return fmt.Sprint(from)
}

func ConvertStringToMany[T int | uint | float32](from int) T {
	return T(from)
}

func ConvertAnyToMany[T, V int | uint | float32](from T) V {
	return V(from)
}
```

3. With error

```go
package converts

import "github.com/shopspring/decimal"

func ConvertStringToDecimal(from string) (decimal.Decimal, error) {
	return decimal.NewFromString(from)
}
```

### Features

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
* [x] Generate convertors by other package fields in models
* [x] Generate convertors by same package conversion functions
* [x] Fix tests
* [x] Delete package flag
* [x] First mapper tests
* [x] Use other conversion functions in convertor
* [x] Use conversion functions with error -> convertor with error
* [x] Convert with pointer field
* [x] Convert with pointer field with error
* [x] No nil err if from and to fields are pointers
* [x] Add CI with tests and linters
* [x] Parse other package
* [x] First console generate
* [x] Set default options
* [x] Add generation info
* [x] Fill readme
* [x] First release
* [x] Get first value by tag
* [x] Move main to root
* [x] Alias package name for (from/to) models and custom cf
* [x] Use conversion functions with pointers
* [x] Use some conversion functions sources
* [x] Parse type alias
* [x] Parse array, slice, map
* [x] Generate convertors by slice fields
* [x] New alias options (alias to each convertors' path)
* [x] Converts both ways in one source
* [x] Generate convertor with from/to pointer
* [x] Fix error by parsing function as member
* [ ] Update readme
* [ ] Map field without tag
* [ ] Use in my projects
* [ ] Generate convertors with map fields
* [ ] Generate convertors with array fields
* [ ] Recursive convert by option if not found conversions
* [ ] Parse user struct in struct
* [ ] Option for default field value if from field is nil
* [ ] Parse comments
* [ ] Use generated convertors in convertor like conversion function
* [ ] Parse embed struct
* [ ] Parse func aliases
* [ ] Warning or error politics if tags is not equals
* [ ] Fill some conversion functions
* [ ] Copy using conversion functions from datamapper to target service if flag set
* [ ] Parse custom error by conversion functions
* [ ] Fix cyclop linter
* [ ] Tag for default field value if from field is nil

### With comment feature example (not implemented):

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
