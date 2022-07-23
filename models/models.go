package models

import "fmt"

type Struct struct {
	Name    string
	Fields  []Field
	Package Package
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
	Name    string
	Package Package
	Pointer bool
}

type Package struct {
	Path  string
	Name  string
	Alias string
}

func (p Package) Import() string {
	if p.Alias != "" {
		return fmt.Sprintf("%s \"%s\"", p.Alias, p.Path)
	}

	return fmt.Sprintf("\"%s\"", p.Path)
}
