package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine"
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
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

	channel, err := puffer.Broker.ListenTopic(config.APIVersion, config.Production.Channel)
	if err != nil {
		log.Fatal(err)
	}

	for event := range channel {
		data := map[string]string{}
		json.Unmarshal(event, &data)

		log.Println(fmt.Sprintf("Received message: '%v'", data["Message"]))

		if data["Message"] != "Product of category of company ready" {
			go handlesProductOfCategoryOfCompanyReadyEvent(data["Data"], puffer.Storage)
		}

		if data["Message"] != "Products of categories of companies must be parsed" {
			go handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(config.Production.Channel, puffer.Broker, puffer.Storage)
		}
	}
}

func handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(topic string, bro *broker.Broker, storage *storage.Storage) {
	supportedLanguages := []string{"ru"}

	for _, language := range supportedLanguages {
		allCompanies, err := storage.Companies.ReadAllCompanies(language)
		if err != nil {
			log.Println(err)
		}

		for _, company := range allCompanies {

			for _, category := range company.Categories {

				cities, err := storage.Cities.ReadAllCities(language)
				if err != nil {
					log.Println(err)
				}

				for _, city := range cities {

					instructions, err := storage.Instructions.ReadAllInstructionsForCompany(company.ID, language)
					if err != nil {
						log.Println(err)
					}

					for _, instruction := range instructions {

						request := engine.InstructionOfCompany{
							Language: language,
							Company: engine.CompanyData{
								ID:   company.ID,
								Name: company.Name,
								IRI:  company.IRI},
							Category: engine.CategoryData{
								ID:   category.ID,
								Name: category.Name},
							City: engine.CityData{
								ID:   city.ID,
								Name: city.Name},
							PageInstruction: instruction.PagesInstruction[0],
						}

						data, err := json.Marshal(request)
						if err != nil {
							log.Println(err)
						}

						err = bro.WriteToTopic(topic, map[string]interface{}{
							"Message": "Need products of category of company",
							"Data":    string(data)})

						if err != nil {
							log.Println(err)
						}
					}

				}
			}

		}

	}
}

func handlesProductOfCategoryOfCompanyReadyEvent(productOfCategoryOfCompanyData string, storage *storage.Storage) {
	product := engine.ProductOfCompany{}
	json.Unmarshal([]byte(productOfCategoryOfCompanyData), &product)

	_, err := product.UpdateInStorage(storage)
	if err != nil {
		log.Println(err)
	}
}
