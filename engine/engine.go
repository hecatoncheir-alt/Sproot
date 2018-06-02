package engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

// Engine is a main object of engine pkg
type Engine struct {
	Configuration *configuration.Configuration
	Storage       *storage.Storage
	Broker        *broker.Broker
}

// New is a constructor for Engine
func New(config *configuration.Configuration) *Engine {
	engine := Engine{Configuration: config}
	return &engine
}

// SetUpStorage for make connect to database and prepare client for requests
func (engine *Engine) SetUpStorage(host string, port int) error {
	engine.Storage = storage.New(host, port)
	err := engine.Storage.SetUp()
	if err != nil {
		return err
	}

	return nil
}

// SetUpBroker for make connect to broker and prepare client for requests
func (engine *Engine) SetUpBroker(host string, port int) error {
	bro := broker.New(engine.Configuration.APIVersion, engine.Configuration.ServiceName)
	engine.Broker = bro

	err := bro.Connect(host, port)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) SubscribeOnEvents(inputTopic string) {

	channel, err := engine.Broker.ListenTopic(engine.Configuration.Production.SprootTopic, engine.Configuration.APIVersion)
	if err != nil {
		log.Fatal(err)
	}

	for event := range channel {
		data := broker.EventData{}
		json.Unmarshal(event, &data)

		if data.APIVersion != engine.Configuration.APIVersion {
			continue
		}

		log.Println(fmt.Sprintf("Received message: '%v'", data.Message))

		if data.Message == "Need items by name" {
			details := storage.ProductsByNameForPage{}
			json.Unmarshal([]byte(data.Data), &details)

			go engine.productsByNameAndPaginationHandler(
				details, data.ClientID, data.APIVersion, engine.Configuration.Production.InitialTopic)
		}

		if data.Message == "Product of category of company ready" {
			go engine.productOfCategoryOfCompanyReadyEventHandler(data.Data)
		}

		if data.Message == "Products of categories of companies must be parsed" {
			go engine.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(engine.Configuration.Production.HecatoncheirTopic)
		}
	}
}

func (engine *Engine) productsByNameAndPaginationHandler(
	details storage.ProductsByNameForPage, clientID, APIVersion string, outputTopic string) {
	productsForPage, err := engine.Storage.Products.ReadProductsByNameWithPagination(
		details.SearchedName, details.Language, details.CurrentPage, details.TotalProductsForOnePage)

	if err != nil && err != storage.ErrProductsByNameNotFound {
		log.Println(err)
	}

	if err != nil && err == storage.ErrProductsByNameNotFound {
		data, err := json.Marshal(productsForPage)
		if err != nil {
			log.Println(err)
		}

		event := broker.EventData{
			Message:    "Items by name not found",
			Data:       string(data),
			APIVersion: APIVersion,
			ClientID:   clientID}

		err = engine.Broker.WriteToTopic(outputTopic, event)
		if err != nil {
			log.Println(err)
		}
	}

	if err == nil {
		data, err := json.Marshal(productsForPage)
		if err != nil {
			log.Println(err)
		}

		event := broker.EventData{
			Message:    "Items by name ready",
			Data:       string(data),
			APIVersion: APIVersion,
			ClientID:   clientID}

		err = engine.Broker.WriteToTopic(outputTopic, event)
		if err != nil {
			log.Println(err)
		}
	}

}

func (engine *Engine) productOfCategoryOfCompanyReadyEventHandler(productOfCategoryOfCompanyData string) {
	product := ProductOfCompany{}
	json.Unmarshal([]byte(productOfCategoryOfCompanyData), &product)

	_, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		log.Println(err)
	}
}

func (engine *Engine) productsOfCategoriesOfCompaniesMustBeParsedEventHandler(outputTopic string) {
	supportedLanguages := []string{"ru"}

	for _, language := range supportedLanguages {
		allCompanies, err := engine.Storage.Companies.ReadAllCompanies(language)
		if err != nil {
			log.Println(err)
		}

		for _, company := range allCompanies {

			for _, category := range company.Categories {

				cities, err := engine.Storage.Cities.ReadAllCities(language)
				if err != nil {
					log.Println(err)
				}

				for _, city := range cities {

					instructions, err := engine.Storage.Instructions.ReadAllInstructionsForCompany(company.ID, language)
					if err != nil {
						log.Println(err)
					}

					for _, instruction := range instructions {

						request := InstructionOfCompany{
							Language: language,
							Company: CompanyData{
								ID:   company.ID,
								Name: company.Name,
								IRI:  company.IRI},
							Category: CategoryData{
								ID:   category.ID,
								Name: category.Name},
							City: CityData{
								ID:   city.ID,
								Name: city.Name},
							PageInstruction: instruction.PagesInstruction[0],
						}

						data, err := json.Marshal(request)
						if err != nil {
							log.Println(err)
						}

						event := broker.EventData{
							Message: "Need products of category of company",
							Data:    string(data)}

						err = engine.Broker.WriteToTopic(outputTopic, event)

						if err != nil {
							log.Println(err)
						}
					}

				}
			}

		}

	}
}
