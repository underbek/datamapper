package main

import (
	"github.com/underbek/datamapper/logger"
	"github.com/underbek/datamapper/mapper"
	"github.com/underbek/datamapper/options"
)

func main() {
	lg := logger.New()
	opts, err := options.ParseOptions()
	if err != nil {
		lg.Fatal(err)
	}

	err = mapper.MapModels(lg, opts)
	if err != nil {
		lg.Fatal(err)
	}
}
