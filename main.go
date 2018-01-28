package main

import (
	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine"
	"log"
)

func main() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	puffer := engine.New()
	err = puffer.SetUpStorage(config.Production.Database.Host, config.Production.Database.Port)
	if err != nil {
		log.Fatal(err)
	}
}
