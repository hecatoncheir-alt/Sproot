package storage

import (
	"testing"
	"time"
)

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

	productsFromStore, err := storage.Products.ReadProductsByName(productForSearch.Name, "en")
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

	productsFromStore, err = storage.Products.ReadProductsByName(createdProduct.Name, "en")
	if err != nil {
		test.Fail()
	}

	if len(productsFromStore) < 1 || len(productsFromStore) > 1 {
		test.Fatal()
	}

	if productsFromStore[0].Name != createdProduct.Name {
		test.Fail()
	}

	if productsFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationProductCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)

	productForSearch := Product{Name: "Test product"}

	productsFromStore, err := storage.Products.ReadProductsByName(productForSearch.Name, "en")
	if err != ErrProductsByNameNotFound {
		test.Fail()
	}

	productFromStore, err := storage.Products.ReadProductByID("0", ".")
	if err != ErrProductDoesNotExist {
		test.Fail()
	}

	if productsFromStore != nil {
		test.Fail()
	}

	createdProduct, err := storage.Products.CreateProduct(productForSearch, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Products.DeleteProduct(createdProduct)

	productFromStore, err = storage.Products.ReadProductByID(createdProduct.ID, ".")
	if err != nil {
		test.Fail()
	}

	if productFromStore.Name != createdProduct.Name {
		test.Fail()
	}

	if productFromStore.IsActive == false {
		test.Fail()
	}

	if productFromStore.ID == "" {
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

func TestIntegrationCategoryCanBeAddedToProduct(test *testing.T) {
	once.Do(prepareStorage)

	createdCategory, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(createdCategory)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	err := storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedProduct, _ := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if len(updatedProduct.Categories) < 1 || len(updatedProduct.Categories) > 1 {
		test.Fatal()
	}

	if updatedProduct.Categories[0].ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeAddedToProduct(test *testing.T) {
	once.Do(prepareStorage)

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(createdCompany)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	err := storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedProduct, _ := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if len(updatedProduct.Companies) < 1 || len(updatedProduct.Companies) > 1 {
		test.Fatal()
	}

	if updatedProduct.Companies[0].ID != createdCompany.ID {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeAddedToProduct(test *testing.T) {
	once.Do(prepareStorage)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	dateTime, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 123, DateTime: dateTime})
	defer storage.Prices.DeletePrice(createdPrice)

	err := storage.Products.AddPriceToProduct(createdProduct.ID, createdPrice.ID)
	if err != nil {
		test.Error(err)
	}

	updatedProduct, _ := storage.Products.ReadProductByID(createdProduct.ID, "en")

	if len(updatedProduct.Prices) < 1 || len(updatedProduct.Prices) > 1 {
		test.Fatal()
	}

	if updatedProduct.Prices[0].ID != createdPrice.ID {
		test.Fail()
	}

	if updatedProduct.Prices[0].Product[0].ID != createdProduct.ID {
		test.Fail()
	}
}
