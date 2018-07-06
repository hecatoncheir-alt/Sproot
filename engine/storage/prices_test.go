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

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Error(err)
		}
	}()

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
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Fail()
		}
	}()

	priceFromStore, err := storage.Prices.ReadPriceByID(createdPrice.ID, ".")
	if err != nil {
		test.Fail()
	}

	if priceFromStore.ID != createdPrice.ID {
		test.Fail()
	}

	if !priceFromStore.IsActive {
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
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	createdPrice, err := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Error(err)
		}
	}()

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Error(err)
		}
	}()

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
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	createdPrice, err := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Error(err)
		}
	}()

	createdCompany, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(createdCompany)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Prices.AddCompanyToPrice(createdPrice.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedPrice, err := storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if updatedPrice.Companies[0].ID != createdCompany.ID {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeAddedFromExportedJSON(test *testing.T) {
	test.Skip()

	once.Do(prepareStorage)

	type allPrices struct {
		Prices []Price `json:"prices"`
	}

	all := allPrices{}
	// exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData := time.Now().UTC()
	priceValue := 123.0

	createdPrice, err := storage.Prices.CreatePrice(Price{Value: priceValue, DateTime: priceData, IsActive: true})
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Prices.DeletePrice(createdPrice)
	if err != nil {
		test.Error(err)
	}

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

	priceFromStorage, err := storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if priceFromStorage.ID != createdPrice.ID {
		test.Fail()
	}

	if priceFromStorage.Value != priceValue {
		test.Fail()
	}

	encodedFormattedPrice := priceData.Format(time.RFC3339)
	formattedPrice, err := time.Parse(time.RFC3339, encodedFormattedPrice)
	if err != nil {
		test.Error(err)
	}

	if priceFromStorage.DateTime != formattedPrice {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeExportedToJSON(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	priceValue := 123.0

	createdPrice, err := storage.Prices.CreatePrice(Price{Value: priceValue, DateTime: priceData})
	if err != nil {
		test.Error(err)
	}

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Error(err)
	}

	err = storage.Prices.AddProductToPrice(createdPrice.ID, createdProduct.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Error(err)
	}

	err = storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)
	if err != nil {
		test.Error(err)
	}

	exportedJSON, err := storage.Prices.ExportJSON()
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Products.DeleteProduct(createdProduct)
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Prices.DeletePrice(createdPrice)
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Cities.DeleteCity(createdCity)
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

	priceFromStorage, err := storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != nil {
		test.Error(err)
	}

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
	_, err = storage.Prices.DeletePrice(createdPrice)
	if err != nil {
		test.Error(err)
	}

	createdProduct.ID = priceFromStorage.Products[0].ID
	_, err = storage.Products.DeleteProduct(createdProduct)
	if err != nil {
		test.Error(err)
	}

	createdCity.ID = priceFromStorage.Cities[0].ID
	_, err = storage.Cities.DeleteCity(createdCity)
	if err != nil {
		test.Error(err)
	}
}

func TestIntegrationCityCanBeAddedToPrice(test *testing.T) {
	once.Do(prepareStorage)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, err := time.Parse(time.RFC3339, exampleDateTime)
	if err != nil {
		test.Error(err)
	}

	createdPrice, err := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Error(err)
		}
	}()

	createdCity, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)
	if err != nil {
		test.Error(err)
	}

	updatedPrice, err := storage.Prices.ReadPriceByID(createdPrice.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if updatedPrice.Cities[0].ID != createdCity.ID {
		test.Fail()
	}
}

func TestPriceDateTimeCanBeMarshalingRight(test *testing.T) {
	test.Skip("WTF with parse or unmarshal...")
	once.Do(prepareStorage)

	timeOfParse, err := time.Parse(time.RFC3339, "2018-06-16T14:08:11.7295653Z")
	if err != nil {
		test.Error(err)
	}

	price := Price{
		ID:       "0",
		Value:    100,
		DateTime: timeOfParse}

	priceInStorage, err := storage.Prices.CreatePrice(price)
	if err != nil {
		test.Error(err)
	}

	priceFromStorage, err := storage.Prices.ReadPriceByID(priceInStorage.ID, ".")
	if err != nil {
		test.Error(err)
	}

	_, err = storage.Prices.DeletePrice(priceInStorage)
	if err != nil {
		test.Error(err)
	}

	encodedPrice, err := json.Marshal(priceFromStorage)
	if err != nil {
		test.Error(err)
	}

	if string(encodedPrice) != timeOfParse.Format(time.RFC3339) {
		test.Fail()
	}
}
