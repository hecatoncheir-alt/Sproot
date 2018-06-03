package engine

import (
	"testing"

	"encoding/json"
	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
	"time"
)

func TestIntegrationEngineCanBeSetUp(test *testing.T) {
	config := configuration.New()

	engine := New(config)
	err := engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}
}

func TestIntegrationEventOfParseRequestCanBeSendToBroker(test *testing.T) {

	test.Skip("broker wtf connect")

	config := configuration.New()
	puffer := New(config)

	var err error

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

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	messages, err := puffer.Broker.ListenTopic(config.Development.SprootTopic, config.APIVersion)

	for event := range messages {
		if event.Message != "Need products of category of company" {
			test.Fail()
		}

		request := InstructionOfCompany{}
		json.Unmarshal([]byte(event.Data), &request)

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
	test.Skip("broker wtf connect")

	config := configuration.New()
	puffer := New(config)

	err := puffer.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
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

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	nameOfProduct := ""
	messages, err := puffer.Broker.ListenTopic(config.Development.SprootTopic, config.APIVersion)
	for event := range messages {

		if event.Message == "Need products of category of company" {

			request := InstructionOfCompany{}
			json.Unmarshal([]byte(event.Data), &request)

			product := ProductOfCompany{
				Language: request.Language,
				Name:     "Test product name",
				Price: PriceOfProduct{
					Value: "1200",
					City: CityData{
						ID:   request.City.ID,
						Name: request.City.Name},
				},
				Company: CompanyData{
					ID:   request.Company.ID,
					Name: request.Company.Name},
				Category: CategoryData{
					ID:   request.Category.ID,
					Name: request.Category.Name},
			}

			productJSON, err := json.Marshal(product)
			if err != nil {
				test.Fail()
			}

			event := broker.EventData{
				Message: "Product of category of company ready",
				Data:    string(productJSON)}

			err = puffer.Broker.WriteToTopic(config.Development.SprootTopic, event)

			continue
		}

		if event.Message != "Product of category of company ready" {
			test.Fail()
		}

		request := ProductOfCompany{}
		json.Unmarshal([]byte(event.Data), &request)

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

		puffer.productOfCategoryOfCompanyReadyEventHandler(event.Data)

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

func TestIntegrationPriceCanBeReturnFromParser(test *testing.T) {
	test.Skip("broker wtf connect")

	config := configuration.New()
	puffer := New(config)

	err := puffer.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
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

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	eventCount := 0
	nameOfProduct := ""
	messages, err := puffer.Broker.ListenTopic(config.Development.SprootTopic, config.APIVersion)
	for event := range messages {
		if event.Message == "Need products of category of company" {

			request := InstructionOfCompany{}
			json.Unmarshal([]byte(event.Data), &request)

			product := ProductOfCompany{
				Language: request.Language,
				Name:     "Exclusive main_test product name",
				Price: PriceOfProduct{
					Value: "1200",
					City: CityData{
						ID:   request.City.ID,
						Name: request.City.Name},
					DateTime: time.Now().UTC(),
				},
				Company: CompanyData{
					ID:   request.Company.ID,
					Name: request.Company.Name},
				Category: CategoryData{
					ID:   request.Category.ID,
					Name: request.Category.Name},
			}

			productJSON, err := json.Marshal(product)
			if err != nil {
				test.Fail()
			}

			event := broker.EventData{
				Message: "Product of category of company ready",
				Data:    string(productJSON)}

			err = puffer.Broker.WriteToTopic(config.Development.SprootTopic, event)

			continue
		}

		if event.Message != "Product of category of company ready" {
			test.Fail()
		}

		request := ProductOfCompany{}
		json.Unmarshal([]byte(event.Data), &request)

		nameOfProduct = request.Name

		puffer.productOfCategoryOfCompanyReadyEventHandler(event.Data)

		go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

		eventCount++
		if eventCount == 2 {
			break
		}
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

	_, err = puffer.Storage.Prices.DeletePrice(products[0].Prices[1])
	if err != nil {
		test.Fail()
	}
}
