package main

import (
	"fmt"

	"github.com/underbek/datamapper/parser"
)

func mapModels(opts Options) error {
	structs, err := parser.ParseStructs(opts.Source)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	from, ok := structs[opts.FromName]
	if !ok {
		return fmt.Errorf("source model %s not found from %s", opts.FromName, opts.Source)
	}

	to, ok := structs[opts.ToName]
	if !ok {
		return fmt.Errorf("to model %s not found from %s", opts.ToName, opts.Source)
	}

	from.Fields = filterFields(opts.FromTag, from.Fields)
	if len(from.Fields) == 0 {
		return fmt.Errorf("soure model %s does not contain tag %s", opts.FromName, opts.FromTag)
	}

	to.Fields = filterFields(opts.ToTag, to.Fields)
	if len(to.Fields) == 0 {
		return fmt.Errorf("to model %s does not contain tag %s", opts.ToName, opts.ToTag)
	}

	return nil
}
