package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/models"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
)

const internalConvertsPackagePath = "github.com/underbek/datamapper/converts"

var (
	ErrNotFoundStruct = errors.New("not found struct error")
	ErrNotFoundTag    = errors.New("not found tag error")
)

func MapModels(lg logger.Logger, opts options.Options) error {
	//TODO: parse or copy embed sources
	funcs, err := parser.ParseConversionFunctionsByPackage(lg, internalConvertsPackagePath)
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	cfAliases := map[string]string{}
	userFuncs := make(models.Functions)

	if len(opts.ConversionFunctions) != 0 {
		for _, cf := range opts.ConversionFunctions {
			res, err := parser.ParseConversionFunctionsByPackage(lg, cf.Source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			for key, function := range res {
				cfAliases[function.Package.Path] = cf.Alias
				userFuncs[key] = function
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

		for cfPath, alias := range cfAliases {
			aliases[cfPath] = alias
		}

		funcs, userFuncs, err = mapModel(
			lg,
			from,
			to,
			opt.From.Tag,
			opt.To.Tag,
			opt.Destination,
			opt.Inverse,
			opt.Recursive,
			opt.WithPointers,
			aliases,
			funcs,
			userFuncs,
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
	aliases map[string]string,
	funcs, userFuncs models.Functions,
	fromStructs, toStructs map[string]models.Struct,
) (models.Functions, models.Functions, error) {

	from.Fields = utils.FilterFields(fromTag, from.Fields)
	if len(from.Fields) == 0 {
		return nil, nil, fmt.Errorf(
			"%w: source model %s does not contain tag %s",
			ErrNotFoundTag,
			from.Type.Name,
			fromTag,
		)
	}

	to.Fields = utils.FilterFields(toTag, to.Fields)
	if len(to.Fields) == 0 {
		return nil, nil, fmt.Errorf(
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
		return nil, nil, fmt.Errorf("create destination dir %s error: %w", path.Dir(destination), err)
	}

	pkg, err := parser.ParseDestinationPackage(lg, destination)
	if err != nil {
		return nil, nil, fmt.Errorf("parse destination package %s error: %w", destination, err)
	}

	var convertors []string
	var pkgs map[models.Package]struct{}
	for {
		userFuncs = setPackageAliasToFunctions(userFuncs, aliases)
		for key, function := range userFuncs {
			funcs[key] = function
		}

		var convertor string
		pkgs, convertor, err = generator.GenerateConvertor(from, to, pkg, funcs)
		if err == nil {
			convertors = append(convertors, convertor)
			break
		}

		if !recursive {
			return nil, nil, err
		}

		var findError *generator.FindFieldsPairError
		if !errors.As(err, &findError) {
			return nil, nil, err
		}

		if findError.From.Package != from.Type.Package || findError.To.Package != to.Type.Package {
			return nil, nil, err
		}

		fromField, fromOk := fromStructs[findError.From.Name]
		toField, toOk := toStructs[findError.To.Name]

		if !fromOk || !toOk {
			return nil, nil, err
		}

		if withPointers {
			fromField.Type.Pointer = findError.From.Pointer
			toField.Type.Pointer = findError.To.Pointer
		}

		funcs, userFuncs, err = mapModel(
			lg,
			fromField,
			toField,
			fromTag,
			toTag,
			generateDestination(fromField.Type.Name, destination),
			inverse,
			recursive,
			withPointers,
			aliases,
			funcs,
			userFuncs,
			fromStructs,
			toStructs,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	if inverse {
		invertPkgs, convertor, err := generator.GenerateConvertor(to, from, pkg, funcs)
		if err != nil {
			return nil, nil, fmt.Errorf("generate convertor error: %w", err)
		}
		convertors = append(convertors, convertor)

		for key := range invertPkgs {
			pkgs[key] = struct{}{}
		}
	}

	err = generator.CreateConvertorSource(pkg, pkgs, convertors, destination)
	if err != nil {
		return nil, nil, fmt.Errorf("create convertor source error: %w", err)
	}
	lg.Infof("generated convertor source: \"%s\"", destination)

	// parse generated functions
	var parsedFunctions models.Functions
	parsedFunctions, err = parser.ParseConversionFunctionsByPackage(lg, destination)
	if err != nil {
		return nil, nil, fmt.Errorf("parse generated conversion functions error: %w", err)
	}
	for key, function := range parsedFunctions {
		userFuncs[key] = function
	}

	return funcs, userFuncs, nil
}

func generateDestination(typeName, dest string) string {
	fileName := strings.ToLower(fmt.Sprintf("%s_converter.go", typeName))
	dir := utils.ClearFileName(dest)
	return fmt.Sprintf("%s/%s", dir, fileName)
}
