package main

import (
	"log"

	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Sproot/engine"
)

func main() {
	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Sproot"
	}

	var err error

	puffer := engine.New(config)
	err = puffer.SetUpStorage(config.Production.Database.Host,
		config.Production.Database.Port)
	if err != nil {
		log.Fatal(err)
	}

	err = puffer.SetUpModel()
	if err != nil {
		log.Fatal(err)
	}

	err = puffer.SetUpBroker(config.Production.EventBus.Host, config.Production.EventBus.Port)
	if err != nil {
		log.Fatal(err)
	}

	puffer.SubscribeOnEvents(config.Production.SprootTopic)
}
