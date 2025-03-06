package main

import (
	"fmt"
	"github.com/mattr/gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.SetUser("matt")
	if err != nil {
		log.Fatal(err)
	}

	updated, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", updated)
}
