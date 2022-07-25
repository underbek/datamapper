package models

import "fmt"

type Struct struct {
	Type   Type
	Fields []Field
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

func (t Type) FullName(basePackage string) string {
	if t.Package.Path == basePackage {
		return t.Name
	}

	if t.Package.Name == "" {
		return t.Name
	}

	if t.Package.Alias == "" {
		return fmt.Sprintf("%s.%s", t.Package.Name, t.Name)
	}

	return fmt.Sprintf("%s.%s", t.Package.Alias, t.Name)
}

func (p Package) Import() string {
	if p.Alias != "" {
		return fmt.Sprintf("%s \"%s\"", p.Alias, p.Path)
	}

	return fmt.Sprintf("\"%s\"", p.Path)
}

//func getFullStructName(model models.Struct, pkgPath string) string {
//	if model.Package.Path != pkgPath {
//		if model.Package.Alias != "" {
//			return fmt.Sprintf("%s.%s", model.Package.Alias, model.Name)
//		}
//		return fmt.Sprintf("%s.%s", model.Package.Name, model.Name)
//	}
//
//	return model.Name
//}
//
//func getFullFieldName(filed models.Field, pkgPath string) string {
//	if filed.Type.Package.Name != "" && filed.Type.Package.Path != pkgPath {
//		if filed.Type.Package.Alias != "" {
//			return fmt.Sprintf("%s.%s", filed.Type.Package.Alias, filed.Type.Name)
//		}
//		return fmt.Sprintf("%s.%s", filed.Type.Package.Name, filed.Type.Name)
//	}
//
//	return filed.Type.Name
//}
