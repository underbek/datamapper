package main

import (
	"log"
	"os"

	"github.com/underbek/datamapper/mapper"
	"github.com/underbek/datamapper/options"
)

func main() {
	opts, err := options.ParseOptions()
	if err != nil {
		os.Exit(1)
	}

	err = mapper.MapModels(opts)
	if err != nil {
		log.Fatal(err)
	}
}
