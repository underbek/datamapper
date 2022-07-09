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

type ConversionFunctionKey struct {
	FromType, ToType string
}

type ConversionFunction struct {
	Name string
}
