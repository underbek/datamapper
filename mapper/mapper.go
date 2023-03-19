package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/loader"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
)

var (
	ErrNotFoundStruct = errors.New("not found struct error")
	ErrNotFoundTag    = errors.New("not found tag error")
)

func MapModels(lg logger.Logger, opts options.Options) error {
	funcs, err := loader.Read()
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	cfAliases := map[string]string{}

	if len(opts.ConversionFunctions) != 0 {
		for _, cf := range opts.ConversionFunctions {
			res, err := parser.ParseConversionFunctionsByPackage(lg, cf.Source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			for key, function := range res {
				cfAliases[function.Package.Path] = cf.Alias
				funcs[key] = function
			}
		}
	}

	for _, opt := range opts.Options {
		fromStructs, err := parser.ParseModelsByPackage(lg, opt.From.Source)
		if err != nil {
			return fmt.Errorf("parse models error: %w", err)
		}

		fromName, isFromPointer := parseModelName(opt.From.Name)
		from, ok := fromStructs[fromName]
		if !ok {
			return fmt.Errorf(" %w: source model %s from %s", ErrNotFoundStruct, opt.From.Name, opt.From.Source)
		}
		from.Type.Pointer = isFromPointer

		toStructs, err := parser.ParseModelsByPackage(lg, opt.To.Source)
		if err != nil {
			return fmt.Errorf("parse models error: %w", err)
		}

		toName, isToPointer := parseModelName(opt.To.Name)
		to, ok := toStructs[toName]
		if !ok {
			return fmt.Errorf("%w: to model %s from %s", ErrNotFoundStruct, opt.To.Name, opt.To.Source)
		}
		to.Type.Pointer = isToPointer

		aliases := map[string]string{
			from.Type.Package.Path: opt.From.Alias,
			to.Type.Package.Path:   opt.To.Alias,
		}

		maps.Copy(aliases, cfAliases)

		funcs, err = mapModel(
			lg,
			from,
			to,
			opt.From.Tag,
			opt.To.Tag,
			opt.Destination,
			opt.Inverse,
			opt.Recursive,
			opt.WithPointers,
			opt.WithSlice,
			aliases,
			funcs,
			fromStructs,
			toStructs,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func setPackageAlias(p *models.Package, aliases map[string]string) {
	p.Alias = aliases[p.Path]
}

func setPackageAliasToStruct(m *models.Struct, aliases map[string]string) {
	setPackageAlias(&m.Type.Package, aliases)
	for i := range m.Fields {
		setPackageAlias(&m.Fields[i].Type.Package, aliases)
		switch m.Fields[i].Type.Kind {
		case models.SliceType:
			additional := m.Fields[i].Type.Additional.(models.SliceAdditional)
			setPackageAlias(&additional.InType.Package, aliases)
			m.Fields[i].Type.Additional = additional
		}
	}
}

func setPackageAliasToCfKey(key models.ConversionFunctionKey, aliases map[string]string) models.ConversionFunctionKey {
	setPackageAlias(&key.FromType.Package, aliases)
	setPackageAlias(&key.ToType.Package, aliases)

	return key
}

func setPackageAliasToCf(cf models.ConversionFunction, aliases map[string]string) models.ConversionFunction {
	setPackageAlias(&cf.Package, aliases)
	setPackageAlias(&cf.FromType.Package, aliases)
	setPackageAlias(&cf.ToType.Package, aliases)

	return cf
}

func setPackageAliasToFunctions(funcs models.Functions, aliases map[string]string) models.Functions {
	res := make(models.Functions)
	for key, cf := range funcs {
		res[setPackageAliasToCfKey(key, aliases)] = setPackageAliasToCf(cf, aliases)
	}
	return res
}

func parseModelName(modelName string) (string, bool) {
	if strings.HasPrefix(modelName, "*") {
		return strings.TrimPrefix(modelName, "*"), true
	}

	return modelName, false
}

func mapModel(
	lg logger.Logger,
	from, to models.Struct,
	fromTag, toTag string,
	destination string,
	inverse bool,
	recursive bool,
	withPointers bool,
	withSlice bool,
	aliases map[string]string,
	funcs models.Functions,
	fromStructs, toStructs map[string]models.Struct,
) (models.Functions, error) {

	from.Fields = utils.FilterFields(fromTag, from.Fields)
	if len(from.Fields) == 0 {
		return nil, fmt.Errorf(
			"%w: source model %s does not contain tag %s",
			ErrNotFoundTag,
			from.Type.Name,
			fromTag,
		)
	}

	to.Fields = utils.FilterFields(toTag, to.Fields)
	if len(to.Fields) == 0 {
		return nil, fmt.Errorf(
			"%w: to model %s does not contain tag %s",
			ErrNotFoundTag,
			to.Type.Name,
			toTag,
		)
	}

	// set aliases
	setPackageAliasToStruct(&from, aliases)
	setPackageAliasToStruct(&to, aliases)

	err := os.MkdirAll(path.Dir(destination), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("create destination dir %s error: %w", path.Dir(destination), err)
	}

	pkg, err := parser.ParseDestinationPackage(lg, destination)
	if err != nil {
		return nil, fmt.Errorf("parse destination package %s error: %w", destination, err)
	}

	var convertors []string
	pkgs := make(models.Packages)
	var gcf models.GeneratedConversionFunction
	for {
		funcs = setPackageAliasToFunctions(funcs, aliases)
		gcf, err = generator.GenerateConvertor(from, to, pkg, funcs)
		if err == nil {
			convertors = append(convertors, gcf.Body)
			funcs[models.ConversionFunctionKey{
				FromType: gcf.Function.FromType,
				ToType:   gcf.Function.ToType,
			}] = gcf.Function
			maps.Copy(pkgs, gcf.Packages)
			break
		}

		if !recursive {
			return nil, err
		}

		var findError *generator.FindFieldsPairError
		if !errors.As(err, &findError) {
			return nil, err
		}

		if findError.From.Package != from.Type.Package || findError.To.Package != to.Type.Package {
			return nil, err
		}

		fromField, fromOk := fromStructs[findError.From.Name]
		toField, toOk := toStructs[findError.To.Name]

		if !fromOk || !toOk {
			return nil, err
		}

		if withPointers {
			fromField.Type.Pointer = findError.From.Pointer
			toField.Type.Pointer = findError.To.Pointer
		}

		funcs, err = mapModel(
			lg,
			fromField,
			toField,
			fromTag,
			toTag,
			generateDestination(fromField.Type.Name, destination),
			inverse,
			recursive,
			withPointers,
			withSlice,
			aliases,
			funcs,
			fromStructs,
			toStructs,
		)
		if err != nil {
			return nil, err
		}
	}

	if withSlice {
		gcf, err := generator.GenerateSliceConvertor(from.Type, to.Type, pkg, gcf.Function)
		if err != nil {
			return nil, fmt.Errorf("generate convertor slice error: %w", err)
		}
		convertors = append(convertors, gcf.Body)
		funcs[models.ConversionFunctionKey{
			FromType: gcf.Function.FromType,
			ToType:   gcf.Function.ToType,
		}] = gcf.Function
		maps.Copy(pkgs, gcf.Packages)
	}

	if inverse {
		gcf, err := generator.GenerateConvertor(to, from, pkg, funcs)
		if err != nil {
			return nil, fmt.Errorf("generate convertor error: %w", err)
		}
		convertors = append(convertors, gcf.Body)
		funcs[models.ConversionFunctionKey{
			FromType: gcf.Function.FromType,
			ToType:   gcf.Function.ToType,
		}] = gcf.Function
		maps.Copy(pkgs, gcf.Packages)

		if withSlice {
			gcf, err := generator.GenerateSliceConvertor(to.Type, from.Type, pkg, gcf.Function)
			if err != nil {
				return nil, fmt.Errorf("generate convertor slice error: %w", err)
			}
			convertors = append(convertors, gcf.Body)
			funcs[models.ConversionFunctionKey{
				FromType: gcf.Function.FromType,
				ToType:   gcf.Function.ToType,
			}] = gcf.Function
			maps.Copy(pkgs, gcf.Packages)
		}
	}

	err = generator.CreateConvertorSource(pkg, pkgs, convertors, destination)
	if err != nil {
		return nil, fmt.Errorf("create convertor source error: %w", err)
	}
	lg.Infof("generated convertor source: \"%s\"", destination)

	return funcs, nil
}

func generateDestination(typeName, dest string) string {
	fileName := strings.ToLower(fmt.Sprintf("%s_converter.go", typeName))
	dir := utils.ClearFileName(dest)
	return fmt.Sprintf("%s/%s", dir, fileName)
}
