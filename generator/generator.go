package generator

import (
	"fmt"
	"os"

	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ConvertorType = string
type ImportType = string

type FieldsPair struct {
	FromName       string
	FromType       string
	ToName         string
	ToType         string
	Assignment     string
	Conversions    []string
	WithError      bool
	PointerToValue bool
}

type result struct {
	fields        []FieldsPair
	imports       []string
	conversations []string
}

func CreateConvertor(from, to models.Struct, dest string, functions models.Functions) error {
	content, err := generateConvertor(from, to, dest, functions)
	if err != nil {
		return err
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func generateConvertor(from, to models.Struct, dest string, functions models.Functions) ([]byte, error) {
	pkg, err := utils.LoadPackage(dest)
	if err != nil {
		return nil, err
	}

	res, err := createModelsPair(from, to, pkg.PkgPath, functions)
	if err != nil {
		return nil, err
	}

	res.imports = append(res.imports, from.PackagePath, to.PackagePath)

	convertorName := "Convert"
	if from.PackagePath != pkg.PkgPath {
		convertorName += cases.Title(language.Und, cases.NoLower).String(from.PackageName)
	}
	convertorName += from.Name
	convertorName += "To"
	if to.PackagePath != pkg.PkgPath {
		convertorName += cases.Title(language.Und, cases.NoLower).String(to.PackageName)
	}
	convertorName += to.Name

	pkgName := pkg.Name
	if pkgName == "" {
		if pkg.PkgPath == "" {
			return nil, fmt.Errorf("incorrect parsed package path from destination %s", dest)
		}
		pkgName = getPackageNameByPath(pkg.PkgPath)
	}

	fromName := getFullStructName(from, pkg.PkgPath)
	toName := getFullStructName(to, pkg.PkgPath)

	return createConvertor(pkgName, fromName, toName, convertorName, pkg.PkgPath, res)
}
