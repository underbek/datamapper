package models

import (
	"fmt"
)

type KindOfType int

const (
	BaseType = KindOfType(iota)
	StructType
	InterfaceType
	RedefinedType
	SliceType
	ArrayType
	MapType
)

type Package struct {
	Path  string
	Name  string
	Alias string
}

// slice -> embed_type
// array -> size, embed_type
// map -> key/value

// struct -> fields

// 1- global map
// 2- write info into string
// 3- difference types from cf key and type
// 4- use kustom map comparable

type Type struct {
	Name       string
	Package    Package
	Pointer    bool
	Kind       KindOfType
	Additional any
}

type ArrayAdditional struct {
	InType Type
	Len    int64
}

type SliceAdditional struct {
	InType Type
}

type MapAdditional struct {
	KeyType   Type
	ValueType Type
}

type Tag struct {
	Name  string
	Value string
}

type Field struct {
	Name string
	Type Type
	Tags []Tag
}

type Struct struct {
	Type   Type
	Fields []Field
}

func (t Type) FullName(basePackage string) string {
	ptr := ""
	if t.Pointer {
		ptr = "*"
	}

	if t.Package.Path == basePackage {
		return ptr + t.Name
	}

	if t.Package.Name == "" {
		return ptr + t.Name
	}

	if t.Package.Alias == "" {
		return fmt.Sprintf("%s%s.%s", ptr, t.Package.Name, t.Name)
	}

	return fmt.Sprintf("%s%s.%s", ptr, t.Package.Alias, t.Name)
}

func (p Package) Import() string {
	if p.Alias != "" {
		return fmt.Sprintf("%s \"%s\"", p.Alias, p.Path)
	}

	return fmt.Sprintf("\"%s\"", p.Path)
}
