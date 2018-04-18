package storage

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIntegrationPriceCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	priceForCreate := Price{Value: 21.440, DateTime: dateTime}
	createdPrice, err := storage.Prices.CreatePrice(priceForCreate)
	if err != nil {
		test.Error(err)
	}

	defer storage.Prices.DeletePrice(createdPrice)

	if createdPrice.ID == "" {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	priceForCreate := Price{Value: 21.440, DateTime: dateTime}
	createdPrice, err := storage.Prices.CreatePrice(priceForCreate)

	defer storage.Prices.DeletePrice(createdPrice)

	priceFromStore, err := storage.Prices.ReadPriceByID(createdPrice.ID, ".")
	if err != nil {
		test.Fail()
	}

	if priceFromStore.ID != createdPrice.ID {
		test.Fail()
	}

	if priceFromStore.IsActive == false {
		test.Fail()
	}

	if priceFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	priceForCreate := Price{Value: 22.440, DateTime: dateTime}

	createdPrice, err := storage.Prices.CreatePrice(priceForCreate)
	if err != nil {
		test.Error(err)
	}

	deletedPriceID, err := storage.Prices.DeletePrice(createdPrice)
	if err != nil {
		test.Error(err)
	}

	if deletedPriceID != createdPrice.ID {
		test.Fail()
	}

	_, err = storage.Prices.ReadPriceByID(deletedPriceID, ".")
	if err != ErrPriceDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationProductCanBeAddedToPrice(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	defer storage.Prices.DeletePrice(createdPrice)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	err := storage.Prices.AddProductToPrice(createdPrice.ID, createdProduct.ID)
	if err != nil {
		test.Error(err)
	}

	updatedPrice, _ := storage.Prices.ReadPriceByID(createdPrice.ID, "en")

	if updatedPrice.Products[0].ID != createdProduct.ID {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeAddedToPrice(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	defer storage.Prices.DeletePrice(createdPrice)

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeactivateCompany(createdCompany)

	err := storage.Prices.AddCompanyToPrice(createdPrice.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedPrice, _ := storage.Prices.ReadPriceByID(createdPrice.ID, "en")

	if updatedPrice.Companies[0].ID != createdCompany.ID {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeAddedFromExportedJSON(test *testing.T) {
	once.Do(prepareStorage)

	type allPrices struct {
		Prices []Price `json:"prices"`
	}

	all := allPrices{}

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData, _ := time.Parse(time.RFC3339, exampleDateTime)
	priceValue := 123.0

	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: priceValue, DateTime: priceData, IsActive: true})
	storage.Prices.DeletePrice(createdPrice)

	all.Prices = append(all.Prices, createdPrice)

	exportedJSON, err := json.Marshal(all)
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != ErrPriceDoesNotExist {
		test.Error(err)
	}

	err = storage.Prices.ImportJSON(exportedJSON)
	if err != nil {
		test.Error(err)
	}

	priceFromStorage, _ := storage.Prices.ReadPriceByID(createdPrice.ID, "en")

	if priceFromStorage.ID != createdPrice.ID {
		test.Fail()
	}

	if priceFromStorage.Value != priceValue {
		test.Fail()
	}

	if priceFromStorage.DateTime != priceData {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeExportedToJSON(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData, _ := time.Parse(time.RFC3339, exampleDateTime)
	priceValue := 123.0

	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: priceValue, DateTime: priceData})

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")

	storage.Prices.AddProductToPrice(createdPrice.ID, createdProduct.ID)

	createdCity, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Error(err)
	}

	storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)

	exportedJSON, err := storage.Prices.ExportJSON()
	if err != nil {
		test.Error(err)
	}

	storage.Products.DeleteProduct(createdProduct)
	storage.Prices.DeletePrice(createdPrice)
	storage.Cities.DeleteCity(createdCity)

	_, err = storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != ErrPriceDoesNotExist {
		test.Error(err)
	}

	err = storage.Prices.ImportJSON(exportedJSON)
	if err != nil {
		test.Error(err)
	}

	priceFromStorage, _ := storage.Prices.ReadPriceByID(createdPrice.ID, "en")

	if priceFromStorage.ID != createdPrice.ID {
		test.Fail()
	}

	if len(priceFromStorage.Products) != 1 {
		test.Fatal()
	}

	if priceFromStorage.Products[0].ID != createdProduct.ID {
		test.Fail()
	}

	if priceFromStorage.Cities[0].ID != createdCity.ID {
		test.Fail()
	}

	createdPrice.ID = priceFromStorage.ID
	storage.Prices.DeletePrice(createdPrice)

	createdProduct.ID = priceFromStorage.Products[0].ID
	storage.Products.DeleteProduct(createdProduct)

	createdCity.ID = priceFromStorage.Cities[0].ID
	storage.Cities.DeleteCity(createdCity)
}

func TestIntegrationCityCanBeAddedToPrice(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	defer storage.Prices.DeletePrice(createdPrice)

	createdCity, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(createdCity)

	err := storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)
	if err != nil {
		test.Error(err)
	}

	updatedPrice, _ := storage.Prices.ReadPriceByID(createdPrice.ID, "en")

	if updatedPrice.Cities[0].ID != createdCity.ID {
		test.Fail()
	}
}
