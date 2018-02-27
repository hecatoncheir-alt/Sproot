package main

import (
	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

func TestIntegrationEventOfParseRequestCanBeSendToBroker(test *testing.T) {

	puffer := engine.New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	err = puffer.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	companyForTest := storage.Company{Name: "М.Видео", IRI: "http://www.mvideo.ru/"}
	createdCompany, err := puffer.Storage.Companies.CreateCompany(companyForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Companies.DeleteCompany(createdCompany)

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := puffer.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Categories.DeleteCategory(createdCategory)

	err = puffer.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := puffer.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Cities.DeleteCity(createdCity)

	pageInstruction, err := puffer.Storage.Instructions.CreatePageInstruction(storage.PageInstruction{Path: "/test/"})
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Instructions.DeletePageInstruction(pageInstruction)

	instruction, err := puffer.Storage.Instructions.CreateInstructionForCompany(createdCompany.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Instructions.DeleteInstruction(instruction)

	puffer.Storage.Instructions.AddPageInstructionToInstruction(instruction.ID, pageInstruction.ID)

	err = puffer.SetUpBroker(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}

	go handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(config.Development.Channel, puffer.Broker, puffer.Storage)

	messages, err := puffer.Broker.ListenTopic(config.Development.Channel, config.Development.Channel)
	for message := range messages {

		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			test.Fail()
		}

		request := engine.InstructionOfCompany{}
		json.Unmarshal([]byte(data["Data"]), &request)

		if request.Language != "ru" {
			test.Fail()
		}

		if request.Company.Name != createdCompany.Name {
			test.Fail()
		}

		if request.Category.Name != createdCategory.Name {
			test.Fail()
		}

		if request.City.Name != createdCity.Name {
			test.Fail()
		}

		if request.PageInstruction.ID != pageInstruction.ID {
			test.Fail()
		}

		break
	}
}

func TestIntegrationProductCanBeReturnFromParser(test *testing.T) {

	puffer := engine.New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	err = puffer.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	companyForTest := storage.Company{Name: "М.Видео", IRI: "http://www.mvideo.ru/"}
	createdCompany, err := puffer.Storage.Companies.CreateCompany(companyForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Companies.DeleteCompany(createdCompany)

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := puffer.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Categories.DeleteCategory(createdCategory)

	err = puffer.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := puffer.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Cities.DeleteCity(createdCity)

	pageInstruction, err := puffer.Storage.Instructions.CreatePageInstruction(storage.PageInstruction{Path: "/test/"})
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Instructions.DeletePageInstruction(pageInstruction)

	instruction, err := puffer.Storage.Instructions.CreateInstructionForCompany(createdCompany.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	defer puffer.Storage.Instructions.DeleteInstruction(instruction)

	puffer.Storage.Instructions.AddPageInstructionToInstruction(instruction.ID, pageInstruction.ID)

	err = puffer.SetUpBroker(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}

	go handlesProductsOfCategoriesOfCompaniesMustBeParsedEvent(config.Development.Channel, puffer.Broker, puffer.Storage)

	nameOfProduct := ""
	messages, err := puffer.Broker.ListenTopic(config.Development.Channel, config.Development.Channel)
	for message := range messages {

		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] == "Need products of category of company" {

			request := engine.InstructionOfCompany{}
			json.Unmarshal([]byte(data["Data"]), &request)

			product := engine.ProductOfCompany{
				Language: request.Language,
				Name:     "Test product name",
				Price: engine.PriceOfProduct{
					Value: "1200",
					City: engine.CityData{
						ID:   request.City.ID,
						Name: request.City.Name},
				},
				Company: engine.CompanyData{
					ID:   request.Company.ID,
					Name: request.Company.Name},
				Category: engine.CategoryData{
					ID:   request.Category.ID,
					Name: request.Category.Name},
			}

			productJSON, err := json.Marshal(product)
			if err != nil {
				test.Fail()
			}

			err = puffer.Broker.WriteToTopic(config.Development.Channel, map[string]interface{}{
				"Message": "Product of category of company ready",
				"Data":    string(productJSON)})

			continue
		}

		if data["Message"] != "Product of category of company ready" {
			test.Fail()
		}

		request := engine.ProductOfCompany{}
		json.Unmarshal([]byte(data["Data"]), &request)

		if request.Name == "" {
			test.Fail()
		}

		nameOfProduct = request.Name

		if request.Price.City.Name != createdCity.Name {
			test.Fail()
		}

		if request.Price.Value == "" {
			test.Fail()
		}

		if request.Language != "ru" {
			test.Fail()
		}

		if request.Company.Name != createdCompany.Name {
			test.Fail()
		}

		if request.Category.Name != createdCategory.Name {
			test.Fail()
		}

		handlesProductOfCategoryOfCompanyReadyEvent(data["Data"], puffer.Storage)

		break
	}

	category, err := puffer.Storage.Categories.ReadCategoryByID(createdCategory.ID, "ru")
	if err != nil {
		test.Fail()
	}

	if len(category.Products) == 0 {
		test.Fatal()
	}

	if category.Products[0].Name != nameOfProduct {
		test.Fail()
	}

	products, err := puffer.Storage.Products.ReadProductsByName(nameOfProduct, "ru")

	if len(products) != 1 {
		test.Fail()
	}

	if products[0].Name != nameOfProduct {
		test.Fail()
	}

	_, err = puffer.Storage.Products.DeleteProduct(products[0])
	if err != nil {
		test.Fail()
	}

	_, err = puffer.Storage.Prices.DeletePrice(products[0].Prices[0])
	if err != nil {
		test.Fail()
	}
}
