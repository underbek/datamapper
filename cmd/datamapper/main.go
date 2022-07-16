package main

import (
	"log"

	"github.com/underbek/datamapper/mapper"
	"github.com/underbek/datamapper/options"
)

func main() {
	opts, err := options.ParseOptions()
	if err != nil {
		log.Fatal(err)
	}

	err = mapper.MapModels(opts)
	if err != nil {
		log.Fatal(err)
	}
}
