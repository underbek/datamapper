package mapper

import (
	"errors"
	"fmt"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
)

const internalConvertsPackagePath = "github.com/underbek/datamapper/converts"

var (
	ErrNotFoundStruct = errors.New("not found struct error")
	ErrNotFoundTag    = errors.New("not found tag error")
)

func MapModels(opts options.Options) error {
	structs, err := parser.ParseModelsByPackage(opts.FromSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	from, ok := structs[opts.FromName]
	if !ok {
		return fmt.Errorf(" %w: source model %s from %s", ErrNotFoundStruct, opts.FromName, opts.FromSource)
	}

	if opts.FromPackageAlias != "" {
		from.Package.Alias = opts.FromPackageAlias
	}

	structs, err = parser.ParseModelsByPackage(opts.ToSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	to, ok := structs[opts.ToName]
	if !ok {
		return fmt.Errorf("%w: to model %s from %s", ErrNotFoundStruct, opts.ToName, opts.ToSource)
	}

	if opts.ToPackageAlias != "" {
		to.Package.Alias = opts.ToPackageAlias
	}

	from.Fields = utils.FilterFields(opts.FromTag, from.Fields)
	if len(from.Fields) == 0 {
		return fmt.Errorf(
			"%w: source model %s does not contain tag %s",
			ErrNotFoundTag,
			opts.FromName,
			opts.FromTag,
		)
	}

	to.Fields = utils.FilterFields(opts.ToTag, to.Fields)
	if len(to.Fields) == 0 {
		return fmt.Errorf("%w: to model %s does not contain tag %s", ErrNotFoundTag, opts.ToName, opts.ToTag)
	}

	//TODO: parse or copy embed sources
	funcs, err := parser.ParseConversionFunctionsByPackage(internalConvertsPackagePath)
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	if len(opts.UserCFSources) != 0 {
		for _, source := range opts.UserCFSources {
			userFuncs, err := parser.ParseConversionFunctionsByPackage(source)
			if err != nil {
				return fmt.Errorf("parse user conversion functions error: %w", err)
			}

			if opts.UserCFPackageAlias != "" {
				for key, cf := range userFuncs {
					cf.Package.Alias = opts.UserCFPackageAlias
					userFuncs[key] = cf
				}
			}

			maps.Copy(funcs, userFuncs)
		}
	}

	err = generator.CreateConvertor(from, to, opts.Destination, funcs)
	if err != nil {
		return fmt.Errorf("generate convertor error: %w", err)
	}

	return nil
}
