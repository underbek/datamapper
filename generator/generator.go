package generator

import (
	"errors"
	"fmt"
	"os"

	"github.com/underbek/datamapper/models"
)

var (
	ErrNotFound                = errors.New("not found error")
	ErrNothingToConvert        = errors.New("nothing to convert error")
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
	Types          []TypeWithName
}

type TypeWithName struct {
	FieldName string
	Type      models.Type
}

type ModelWithPairs struct {
	Type   TypeWithName
	models []ModelWithPairs
	fields []FieldsPair
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

func GenerateConvertor(from, to models.Struct, fromTag, toTag string, pkg models.Package, functions models.Functions) (
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

	res.fromTag = fromTag
	res.toTag = toTag

	if len(res.fields) == 0 {
		return models.GeneratedConversionFunction{}, fmt.Errorf(
			"%w %s by tag %s -> %s by tag %s",
			ErrNothingToConvert,
			from.Type.Name,
			fromTag,
			to.Type.Name,
			toTag,
		)
	}

	headType := findHead(res.fields)
	modelWithPairs := createModelWithPairs(res.fields, headType)

	resultStruct, err := createResultConverter(pkg.Path, res.toName, modelWithPairs)
	if err != nil {
		return models.GeneratedConversionFunction{}, err
	}

	convertor, err := fillConvertor(res, resultStruct)
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
			Name:    res.convertorName,
			Package: pkg,
			FromType: models.Type{
				Kind: models.SliceType,
				Additional: models.SliceAdditional{
					InType: from,
				},
			},
			ToType: models.Type{
				Kind: models.SliceType,
				Additional: models.SliceAdditional{
					InType: to,
				},
			},
			TypeParam: models.NoTypeParam,
			WithError: res.withError,
		},
		Packages: res.packages,
		Body:     convertor,
	}, nil
}
