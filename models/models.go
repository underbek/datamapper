package models

type Struct struct {
	Name        string
	Fields      []Field
	PackagePath string
	PackageName string
}

type Field struct {
	Name string
	Type Type
	Tags []Tag
}

type Tag struct {
	Name  string
	Value string
}

type Type struct {
	Name        string
	PackagePath string
	Pointer     bool
}
