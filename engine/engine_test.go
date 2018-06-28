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

func TestIntegrationProductCanBeReturnFromParser(test *testing.T) {
	test.Skip("Test it late")

	config := configuration.New()
	puffer := New(config)

	err := puffer.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	err = puffer.SetUpBroker(config.Development.EventBus.Host, config.Development.EventBus.Port)
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

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	nameOfProduct := ""

	otherBroker := broker.New(config.APIVersion, config.ServiceName)
	err = otherBroker.Connect(config.Development.EventBus.Host, 8181)
	if err != nil {
		test.Error(err)
	}

	for event := range otherBroker.InputChannel {

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

			go puffer.Broker.Write(event)

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
		test.Error(err)
	}

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
	if err != nil {
		test.Error(err)
	}

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

	//err = puffer.SetUpBroker(config.Development.EventBus.Host, config.Development.EventBus.Port)
	err = puffer.SetUpBroker(config.Development.EventBus.Host, 8181)
	if err != nil {
		test.Error(err)
	}

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	nameOfProduct := ""

	for event := range puffer.Broker.OutputChannel {
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

			go func() {
				puffer.Broker.InputChannel <- event
			}()

			break
		}
	}

	for event := range puffer.Broker.InputChannel {
		if event.Message != "Product of category of company ready" {
			test.Fail()
		}

		request := ProductOfCompany{}
		json.Unmarshal([]byte(event.Data), &request)

		nameOfProduct = request.Name

		puffer.productOfCategoryOfCompanyReadyEventHandler(event.Data)

		go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

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
	if err != nil {
		test.Error(err)
	}

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
