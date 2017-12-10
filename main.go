package main

import (
	"log"

	"github.com/hecatoncheir/Sproot/engine"
)

func main() {
	puffer := engine.New()

	err := puffer.SetUpStorage("192.168.99.100", 9080)
	if err != nil {
		log.Fatal(err)
	}
}
