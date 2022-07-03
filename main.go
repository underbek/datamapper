package main

import (
	"log"
)

func main() {
	opts, err := parseOptions()
	if err != nil {
		log.Fatal(err)
	}

	err = mapModels(opts)
	if err != nil {
		log.Fatal(err)
	}
}
