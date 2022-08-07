package generator

import (
	"errors"
	"os"

	"github.com/underbek/datamapper/models"
)

var (
	ErrNotFound                = errors.New("not found error")
	ErrUndefinedConversionRule = errors.New("undefined conversion rule error")
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
	convertorName string
	pkg           models.Package
	fromName      string
	toName        string
	fromTag       string
	toTag         string
	fields        []FieldsPair
	packages      map[models.Package]struct{}
	conversions   []string
	withError     bool
}

func CreateConvertorSource(pkg models.Package, packages map[models.Package]struct{}, convertors []string,
	dest string) error {

	content, err := fillConvertorsSource(pkg, packages, convertors)
	if err != nil {
		return err
	}

	file, err := os.Create(dest) //nolint:gosec
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

func GenerateConvertor(from, to models.Struct, pkg models.Package, functions models.Functions) (
	map[models.Package]struct{}, string, error) {

	res, err := createModelsPair(from, to, pkg.Path, functions)
	if err != nil {
		return nil, "", err
	}

	res.packages[from.Type.Package] = struct{}{}
	res.packages[to.Type.Package] = struct{}{}

	res.convertorName = generateConvertorName(from, to, pkg.Path)

	res.fromName = from.Type.FullName(pkg.Path)
	res.toName = to.Type.FullName(pkg.Path)

	res.fromTag = from.Fields[0].Tags[0].Name
	res.toTag = to.Fields[0].Tags[0].Name

	res.withError = isReturnError(res.fields)

	convertor, err := fillConvertor(res)
	if err != nil {
		return nil, "", err
	}

	return res.packages, convertor, nil
}
