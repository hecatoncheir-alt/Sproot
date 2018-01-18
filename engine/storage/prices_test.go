package storage

import (
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

	if updatedPrice.Product[0].ID != createdProduct.ID {
		test.Fail()
	}
}
