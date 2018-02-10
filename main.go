package main

import (
	"encoding/json"
	"fmt"
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

	err = puffer.SetUpBroker(config.Production.Broker.Host, config.Production.Broker.Port)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := puffer.Broker.ListenTopic(config.ApiVersion, config.Production.Channel)
	if err != nil {
		log.Fatal(err)
	}

	for event := range channel {
		data := map[string]string{}
		json.Unmarshal(event, &data)

		if data["Message"] != "Product of category of company ready" {
			log.Println(fmt.Sprintf("Received message: '%v'", data["Message"]))
			go handlesProductOfCategoryOfCompanyReadyEvent(data["Data"])
		}
	}
}

func handlesProductOfCategoryOfCompanyReadyEvent(productOfCategoryOfCompanyData string) {
	product := engine.ProductOfCompany{}
	json.Unmarshal([]byte(productOfCategoryOfCompanyData), &product)

	err := product.UpdateInStorage()
	if err != nil {
		log.Println(err)
	}
}
