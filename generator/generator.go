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
	fromName      string
	toName        string
	fromTag       string
	toTag         string
	fields        []FieldsPair
	packages      models.Packages
	conversions   []string
	withError     bool
}

type sliceResult struct {
	convertorName string
	fromName      string
	toName        string
	packages      models.Packages
	conversion    string
	withError     bool
}

func CreateConvertorSource(pkg models.Package, packages models.Packages, convertors []string, dest string) error {
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
	models.GeneratedConversionFunction, error,
) {

	res, err := createModelsPair(from, to, pkg.Path, functions)
	if err != nil {
		return models.GeneratedConversionFunction{}, err
	}

	res.packages[from.Type.Package] = struct{}{}
	res.packages[to.Type.Package] = struct{}{}

	res.convertorName = generateConvertorName(from.Type, to.Type, pkg.Path, models.StructType)

	res.fromName = from.Type.FullName(pkg.Path)
	res.toName = to.Type.FullName(pkg.Path)

	res.fromTag = from.Fields[0].Tags[0].Name
	res.toTag = to.Fields[0].Tags[0].Name

	convertor, err := fillConvertor(res)
	if err != nil {
		return models.GeneratedConversionFunction{}, err
	}

	return models.GeneratedConversionFunction{
		Function: models.ConversionFunction{
			Name:      res.convertorName,
			Package:   pkg,
			FromType:  from.Type,
			ToType:    to.Type,
			TypeParam: models.NoTypeParam,
			WithError: res.withError,
		},
		Packages: res.packages,
		Body:     convertor,
	}, nil
}

func GenerateSliceConvertor(from, to models.Type, pkg models.Package, cf models.ConversionFunction) (
	models.GeneratedConversionFunction, error,
) {
	res := sliceResult{}

	res.packages = make(models.Packages)

	res.packages[from.Package] = struct{}{}
	res.packages[to.Package] = struct{}{}

	res.convertorName = generateConvertorName(from, to, pkg.Path, models.SliceType)

	res.fromName = from.FullName(pkg.Path)
	res.toName = to.FullName(pkg.Path)

	res.conversion = getConversionFunctionCall(cf, from, to, pkg.Path, "from")

	res.withError = cf.WithError

	convertor, err := fillSliceConvertor(res)
	if err != nil {
		return models.GeneratedConversionFunction{}, err
	}

	return models.GeneratedConversionFunction{
		Function: models.ConversionFunction{
			Name:      res.convertorName,
			Package:   pkg,
			FromType:  from,
			ToType:    to,
			TypeParam: models.NoTypeParam,
			WithError: res.withError,
		},
		Packages: res.packages,
		Body:     convertor,
	}, nil
}
