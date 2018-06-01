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

	// channel, err := puffer.Broker.ListenTopic(config.Production.SprootTopic, config.APIVersion)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for event := range channel {
	// 	data := broker.EventData{}
	// 	json.Unmarshal(event, &data)

	// 	if data.APIVersion != config.APIVersion {
	// 		continue
	// 	}

	// 	log.Println(fmt.Sprintf("Received message: '%v'", data.Message))

	// 	if data.Message == "Need items by name" {
	// 		details := storage.ProductsByNameForPage{}
	// 		json.Unmarshal([]byte(data.Data), &details)
	// 		go handlesProductsByNameAndPagination(details, data.ClientID, data.APIVersion, config.Production.InitialTopic, puffer.Broker, puffer.Storage)

	// 	}

	// 	if data.Message == "Product of category of company ready" {
	// 		go handlesProductOfCategoryOfCompanyReadyEvent(data.Data, puffer.Storage)
	// 	}

	// 	if data.Message == "Products of categories of companies must be parsed" {
	// 		go handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(config.Production.SprootTopic, puffer.Broker, puffer.Storage)
	// 	}
	// }
}

// func handlesProductsByNameAndPagination(details storage.ProductsByNameForPage, clientID, APIVersion string, topic string, bro *broker.Broker, store *storage.Storage) {
// 	productsForPage, err := store.Products.ReadProductsByNameWithPagination(details.SearchedName, details.Language, details.CurrentPage, details.TotalProductsForOnePage)

// 	if err != nil && err != storage.ErrProductsByNameNotFound {
// 		log.Println(err)
// 	}

// 	if err != nil && err == storage.ErrProductsByNameNotFound {
// 		data, err := json.Marshal(productsForPage)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		event := broker.EventData{
// 			Message:    "Items by name not found",
// 			Data:       string(data),
// 			APIVersion: APIVersion,
// 			ClientID:   clientID}

// 		err = bro.WriteToTopic(topic, event)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// 	if err == nil {
// 		data, err := json.Marshal(productsForPage)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		event := broker.EventData{
// 			Message:    "Items by name ready",
// 			Data:       string(data),
// 			APIVersion: APIVersion,
// 			ClientID:   clientID}

// 		err = bro.WriteToTopic(topic, event)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}

// }

// func handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(topic string, bro *broker.Broker, storage *storage.Storage) {
// 	supportedLanguages := []string{"ru"}

// 	for _, language := range supportedLanguages {
// 		allCompanies, err := storage.Companies.ReadAllCompanies(language)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		for _, company := range allCompanies {

// 			for _, category := range company.Categories {

// 				cities, err := storage.Cities.ReadAllCities(language)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				for _, city := range cities {

// 					instructions, err := storage.Instructions.ReadAllInstructionsForCompany(company.ID, language)
// 					if err != nil {
// 						log.Println(err)
// 					}

// 					for _, instruction := range instructions {

// 						request := engine.InstructionOfCompany{
// 							Language: language,
// 							Company: engine.CompanyData{
// 								ID:   company.ID,
// 								Name: company.Name,
// 								IRI:  company.IRI},
// 							Category: engine.CategoryData{
// 								ID:   category.ID,
// 								Name: category.Name},
// 							City: engine.CityData{
// 								ID:   city.ID,
// 								Name: city.Name},
// 							PageInstruction: instruction.PagesInstruction[0],
// 						}

// 						data, err := json.Marshal(request)
// 						if err != nil {
// 							log.Println(err)
// 						}

// 						event := broker.EventData{
// 							Message: "Need products of category of company",
// 							Data:    string(data)}

// 						err = bro.WriteToTopic(topic, event)

// 						if err != nil {
// 							log.Println(err)
// 						}
// 					}

// 				}
// 			}

// 		}

// 	}
// }

// func handlesProductOfCategoryOfCompanyReadyEvent(productOfCategoryOfCompanyData string, storage *storage.Storage) {
// 	product := engine.ProductOfCompany{}
// 	json.Unmarshal([]byte(productOfCategoryOfCompanyData), &product)

// 	_, err := product.UpdateInStorage(storage)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
