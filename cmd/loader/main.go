package main

import (
	"github.com/underbek/datamapper/loader"
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/parser"
)

const internalConvertsPackagePath = "github.com/underbek/datamapper/converts"

func main() {
	lg := logger.New()
	funcs, err := parser.ParseConversionFunctionsByPackage(lg, internalConvertsPackagePath)
	if err != nil {
		lg.Fatalf("parse internal conversion functions error: %w", err)
	}

	err = loader.Save(funcs)
	if err != nil {
		lg.Fatalf("parse internal conversion functions error: %w", err)
	}
}
