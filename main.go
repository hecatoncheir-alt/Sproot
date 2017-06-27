package main

import (
	"log"

	"github.com/hecatoncheir/Sproot/engine"
)

func main() {
	puffer := engine.New()
	err := puffer.DatabaseSetUp("http", "192.168.99.100", 8080)

	if err != nil {
		log.Fatal(err)
	}
}
