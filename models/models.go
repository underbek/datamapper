package models

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	Tags []Tag
}

type Tag struct {
	Name  string
	Value string
}

type Type struct {
	Name    string
	Package string
}

type ConversionFunctionKey struct {
	FromType, ToType Type
}

type ConversionFunction struct {
	Name string
}
