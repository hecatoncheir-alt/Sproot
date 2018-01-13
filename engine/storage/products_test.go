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

func TestIntegrationProductsCanBeReadByName(test *testing.T) {
	once.Do(prepareStorage)

	productForSearch := Product{Name: "Test product"}

	productsFromStore, err := storage.Products.ReadProductsByName(productForSearch.Name, ".")
	if err != ErrProductsByNameNotFound {
		test.Fail()
	}

	if productsFromStore != nil {
		test.Fail()
	}

	createdProduct, err := storage.Products.CreateProduct(productForSearch, "en")
	if err != nil || createdProduct.ID == "" {
		test.Error(err)
	}

	defer storage.Products.DeleteProduct(createdProduct)

	productsFromStore, err = storage.Products.ReadProductsByName(createdProduct.Name, ".")
	if err != nil {
		test.Fail()
	}

	if productsFromStore[0].Name != createdProduct.Name {
		test.Fail()
	}

	if productsFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationProductCanHasNameWithManyLanguages(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")

	defer storage.Products.DeleteProduct(createdProduct)

	err = storage.Products.AddLanguageOfProductName(createdProduct.ID, "Тестовый продукт", "ru")
	if err != nil {
		test.Fail()
	}

	productWithEnName, _ := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if productWithEnName.Name != "Test product" {
		test.Fail()
	}

	productWithRuName, _ := storage.Products.ReadProductByID(createdProduct.ID, "ru")
	if productWithRuName.Name != "Тестовый продукт" {
		test.Fail()
	}
}

func TestIntegrationProductCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	productForCreate := Product{Name: "Test product"}
	createdProduct, err := storage.Products.CreateProduct(productForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	deletedProductID, err := storage.Products.DeleteProduct(createdProduct)
	if err != nil {
		test.Error(err)
	}

	if deletedProductID != createdProduct.ID {
		test.Fail()
	}

	_, err = storage.Products.ReadProductByID(deletedProductID, ".")
	if err != ErrProductDoesNotExist {
		test.Error(err)
	}
}
