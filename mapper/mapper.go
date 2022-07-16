package mapper

import (
	"fmt"

	"github.com/underbek/datamapper/generator"
	"github.com/underbek/datamapper/options"
	"github.com/underbek/datamapper/parser"
	"github.com/underbek/datamapper/utils"
	"golang.org/x/exp/maps"
)

const (
	internalConvertsPackagePath = "../converts"
)

func MapModels(opts options.Options) error {
	if opts.FromSource == "" {
		opts.FromSource = "."
	}

	if opts.ToSource == "" {
		opts.ToSource = "."
	}

	structs, err := parser.ParseModels(opts.FromSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	from, ok := structs[opts.FromName]
	if !ok {
		return fmt.Errorf("source model %s not found from %s", opts.FromName, opts.FromSource)
	}

	structs, err = parser.ParseModels(opts.ToSource)
	if err != nil {
		return fmt.Errorf("parse models error: %w", err)
	}

	to, ok := structs[opts.ToName]
	if !ok {
		return fmt.Errorf("to model %s not found from %s", opts.ToName, opts.ToSource)
	}

	from.Fields = utils.FilterFields(opts.FromTag, from.Fields)
	if len(from.Fields) == 0 {
		return fmt.Errorf("soure model %s does not contain tag %s", opts.FromName, opts.FromTag)
	}

	to.Fields = utils.FilterFields(opts.ToTag, to.Fields)
	if len(to.Fields) == 0 {
		return fmt.Errorf("to model %s does not contain tag %s", opts.ToName, opts.ToTag)
	}

	//TODO: parse or copy embed sources
	funcs, err := parser.ParseConversionFunctions(internalConvertsPackagePath)
	if err != nil {
		return fmt.Errorf("parse internal conversion functions error: %w", err)
	}

	if opts.UserCFSource != "" {
		userFuncs, err := parser.ParseConversionFunctions(opts.UserCFSource)
		if err != nil {
			return fmt.Errorf("parse user conversion functions error: %w", err)
		}

		maps.Copy(funcs, userFuncs)
	}

	err = generator.CreateConvertor(from, to, opts.Destination, funcs)
	if err != nil {
		return fmt.Errorf("generate convertor error: %w", err)
	}

	return nil
}
