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

	err = puffer.SetUpBroker(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Расставить логи при получении запроса с клиента.
	// При получении данных от парсера. При отправки ответа на клиент.

	puffer.SubscribeOnEvents(puffer.Configuration.Production.SprootTopic)
}
