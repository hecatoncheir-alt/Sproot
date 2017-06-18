package main

import (
	"log"

	"github.com/hecatoncheir/Sproot/engine"
)

func main() {
	puffer := engine.New()
	_, _, err := puffer.DatabaseSetUp("root", "192.168.99.100", 26257, "disable", "Items")

	if err != nil {
		log.Fatal(err)
	}
}
