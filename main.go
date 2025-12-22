package main

import (
	"fmt"
	"log"

	"github.com/westleaf/gator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = conf.SetUser("westy")
	if err != nil {
		log.Fatal(err)
	}

	conf, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", conf)
}
