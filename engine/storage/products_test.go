package storage

import "testing"

func TestIntegrationProductCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	productForCreate := Product{Name: "Test product"}
	createdProduct, err := storage.Products.CreateProduct(productForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Products.DeleteProduct(createdProduct)

	if createdProduct.ID == "" {
		test.Fail()
	}

	existProduct, err := storage.Products.CreateProduct(productForCreate, "en")
	if err == nil || err != ErrProductAlreadyExist {
		test.Error(err)
	}

	if existProduct.ID != createdProduct.ID {
		test.Fail()
	}
}
