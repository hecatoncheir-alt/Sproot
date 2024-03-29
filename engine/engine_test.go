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
	// TODO
	// Что-то не так при обновлении на новую версия базы данных 1.0.6
	test.Skip()

	config := configuration.New()

	engine := New(config)
	err := engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}
}

func TestIntegrationProductCanBeReturnFromParser(test *testing.T) {
	// TODO
	// Что-то не так при обновлении на новую версия базы данных 1.0.6
	// Отдельно тест проходит
	test.Skip()

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

	defer func() {
		_, err := puffer.Storage.Companies.DeleteCompany(createdCompany)
		if err != nil {
			test.Error(err)
		}
	}()

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := puffer.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Categories.DeleteCategory(createdCategory)
		if err != nil {
			test.Error(err)
		}
	}()

	err = puffer.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := puffer.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Error(err)
		}
	}()

	pageInstruction, err := puffer.Storage.Instructions.CreatePageInstruction(storage.PageInstruction{Path: "/test/"})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Instructions.DeletePageInstruction(pageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := puffer.Storage.Instructions.CreateInstructionForCompany(createdCompany.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	err = puffer.Storage.Instructions.AddPageInstructionToInstruction(instruction.ID, pageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	nameOfProduct := ""

	puffer.Broker = broker.New(config.APIVersion, config.ServiceName)

	for event := range puffer.Broker.OutputChannel {

		if event.Message == "Need products of category of company" {

			request := InstructionOfCompany{}
			err := json.Unmarshal([]byte(event.Data), &request)
			if err != nil {
				test.Error(err)
			}

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
		err := json.Unmarshal([]byte(event.Data), &request)
		if err != nil {
			test.Error(err)
		}

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

		close(puffer.Broker.InputChannel)
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
	// TODO
	// Что-то не так при обновлении на новую версия базы данных 1.0.6
	// Отдельно тест проходит
	test.Skip()

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

	defer func() {
		_, err := puffer.Storage.Companies.DeleteCompany(createdCompany)
		if err != nil {
			test.Error(err)
		}
	}()

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := puffer.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Categories.DeleteCategory(createdCategory)
		if err != nil {
			test.Error(err)
		}
	}()

	err = puffer.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := puffer.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Error(err)
		}
	}()

	pageInstruction, err := puffer.Storage.Instructions.CreatePageInstruction(storage.PageInstruction{Path: "/test/"})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Instructions.DeletePageInstruction(pageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := puffer.Storage.Instructions.CreateInstructionForCompany(createdCompany.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := puffer.Storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	err = puffer.Storage.Instructions.AddPageInstructionToInstruction(instruction.ID, pageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	puffer.Broker = broker.New(config.APIVersion, config.ServiceName)

	go puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

	nameOfProduct := ""

	for event := range puffer.Broker.OutputChannel {
		if event.Message == "Need products of category of company" {

			request := InstructionOfCompany{}
			err := json.Unmarshal([]byte(event.Data), &request)
			if err != nil {
				test.Error(err)
			}

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
		err = json.Unmarshal([]byte(event.Data), &request)
		if err != nil {
			test.Error(err)
		}

		nameOfProduct = request.Name

		puffer.productOfCategoryOfCompanyReadyEventHandler(event.Data)

		go func() {
			puffer.productsOfCategoriesOfCompaniesMustBeParsedEventHandler(config.Development.SprootTopic)

			close(puffer.Broker.InputChannel)
		}()
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
