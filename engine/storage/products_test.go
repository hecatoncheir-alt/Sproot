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

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Error(err)
		}
	}()

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

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Fail()
		}
	}()

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

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Error(err)
		}
	}()

	productFromStore, err = storage.Products.ReadProductByID(createdProduct.ID, ".")
	if err != nil {
		test.Fail()
	}

	if productFromStore.Name != createdProduct.Name {
		test.Fail()
	}

	if !productFromStore.IsActive {
		test.Fail()
	}

	if productFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationProductCanHasNameWithManyLanguages(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Products.AddLanguageOfProductName(createdProduct.ID, "Тестовый продукт", "ru")
	if err != nil {
		test.Fail()
	}

	productWithEnName, err := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if err != nil {
		test.Fail()
	}

	if productWithEnName.Name != "Test product" {
		test.Fail()
	}

	productWithRuName, err := storage.Products.ReadProductByID(createdProduct.ID, "ru")
	if err != nil {
		test.Fail()
	}

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

	createdCategory, err := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Categories.DeleteCategory(createdCategory)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)
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

	createdCompany, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(createdCompany)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedProduct, err := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if err != nil {
		test.Fail()
	}

	if len(updatedProduct.Companies) < 1 || len(updatedProduct.Companies) > 1 {
		test.Fatal()
	}

	if updatedProduct.Companies[0].ID != createdCompany.ID {
		test.Fail()
	}
}

func TestIntegrationPriceCanBeAddedToProduct(test *testing.T) {
	once.Do(prepareStorage)

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

	err = storage.Products.AddPriceToProduct(createdProduct.ID, createdPrice.ID)
	if err != nil {
		test.Error(err)
	}

	updatedProduct, err := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(updatedProduct.Prices) < 1 || len(updatedProduct.Prices) > 1 {
		test.Fatal()
	}

	if updatedProduct.Prices[0].ID != createdPrice.ID {
		test.Fail()
	}

	if updatedProduct.Prices[0].Products[0].ID != createdProduct.ID {
		test.Fail()
	}
}

func TestIntegrationProductCanHasPrice(test *testing.T) {
	once.Do(prepareStorage)
	createdPrice, err := storage.Prices.CreatePrice(Price{Value: 31.2})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdPrice)
		if err != nil {
			test.Error(err)
		}
	}()

	createdCity, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)
	if err != nil {
		test.Error(err)
	}

	createdProduct, err := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Products.AddPriceToProduct(createdProduct.ID, createdPrice.ID)
	if err != nil {
		test.Error(err)
	}

	productInStore, err := storage.Products.ReadProductByID(createdProduct.ID, "en")
	if err != nil {
		test.Fail()
	}

	if productInStore.Prices[0].ID != createdPrice.ID {
		test.Fail()
	}

	if productInStore.Prices[0].Cities[0].ID != createdCity.ID {
		test.Fail()
	}
}

func TestIntegrationProductsTotalCountCanBeReturned(test *testing.T) {
	once.Do(prepareStorage)

	createdProduct1, err := storage.Products.CreateProduct(Product{Name: "Первый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct1)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct2, err := storage.Products.CreateProduct(Product{Name: "Второй тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct2)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct3, err := storage.Products.CreateProduct(Product{Name: "Третий тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct3)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct4, err := storage.Products.CreateProduct(Product{Name: "Четвёртый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct4)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct5, err := storage.Products.CreateProduct(Product{Name: "Пятый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct5)
		if err != nil {
			test.Fail()
		}
	}()

	counterOfFoundedProducts, err := storage.Products.ReadTotalCountOfProductsByName("тестовый", "ru")
	if err != nil {
		test.Error(err)
	}

	if counterOfFoundedProducts == 0 {
		test.Errorf("Expected 5, actual: %v", counterOfFoundedProducts)
	}
}

func TestIntegrationProductCanBeSearchedByNameWithPaginationAndCounter(test *testing.T) {
	once.Do(prepareStorage)

	createdProduct1, err := storage.Products.CreateProduct(Product{Name: "Первый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct1)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct2, err := storage.Products.CreateProduct(Product{Name: "Второй тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct2)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct3, err := storage.Products.CreateProduct(Product{Name: "Третий тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct3)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct4, err := storage.Products.CreateProduct(Product{Name: "Четвёртый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct4)
		if err != nil {
			test.Fail()
		}
	}()

	createdProduct5, err := storage.Products.CreateProduct(Product{Name: "Пятый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct5)
		if err != nil {
			test.Fail()
		}
	}()

	foundedProductsForFirstPage, err := storage.Products.ReadProductsByNameWithPagination("тестовый", "ru", 1, 2)
	if err != nil {
		test.Error(err)
	}

	if foundedProductsForFirstPage.TotalProductsFound != 5 {
		test.Errorf("Expected 5 products, actual: %v", foundedProductsForFirstPage.TotalProductsFound)
	}

	if foundedProductsForFirstPage.TotalProductsForOnePage != 2 {
		test.Errorf("Expected 2 products on one page, actual: %v", foundedProductsForFirstPage.TotalProductsForOnePage)
	}

	if foundedProductsForFirstPage.CurrentPage != 1 {
		test.Errorf("Expected page 1, actual: %v", foundedProductsForFirstPage.CurrentPage)
	}

	if foundedProductsForFirstPage.Products[0].Name != "Первый тестовый продукт" {
		test.Errorf("Expected \"Первый тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[0].Name)
	}

	if foundedProductsForFirstPage.Products[1].Name != "Второй тестовый продукт" {
		test.Errorf("Expected \"Второй тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[1].Name)
	}

	foundedProductsForFirstPage, err = storage.Products.ReadProductsByNameWithPagination("тестовый", "ru", 2, 2)
	if err != nil {
		test.Error(err)
	}

	if foundedProductsForFirstPage.Products[0].Name != "Третий тестовый продукт" {
		test.Errorf("Expected \"Третий тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[0].Name)
	}

	if foundedProductsForFirstPage.Products[1].Name != "Четвёртый тестовый продукт" {
		test.Errorf("Expected \"Четвёртый тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[1].Name)
	}

	foundedProductsForFirstPage, err = storage.Products.ReadProductsByNameWithPagination("тестовый", "ru", 3, 2)
	if err != nil {
		test.Error(err)
	}

	if foundedProductsForFirstPage.CurrentPage != 3 {
		test.Errorf("Expected page 3, actual: %v", foundedProductsForFirstPage.CurrentPage)
	}

	if len(foundedProductsForFirstPage.Products) != 1 {
		test.Errorf("Expected 1 product for one page, actual: %v", len(foundedProductsForFirstPage.Products))
	}

	if len(foundedProductsForFirstPage.Products) != 1 {
		test.Errorf("Expected one product. actual: %v", len(foundedProductsForFirstPage.Products))
	}

	if foundedProductsForFirstPage.Products[0].Name != "Пятый тестовый продукт" {
		test.Errorf("Expected \"Пятый тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[0].Name)
	}
}

func TestPricesOfProductMustBeSortedByDate(test *testing.T) {
	once.Do(prepareStorage)

	createdProduct1, err := storage.Products.CreateProduct(Product{Name: "Первый тестовый продукт"}, "ru")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Products.DeleteProduct(createdProduct1)
		if err != nil {
			test.Error(err)
		}
	}()

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

	err = storage.Products.AddPriceToProduct(createdProduct1.ID, createdPrice.ID)
	if err != nil {
		test.Error(err)
	}

	exampleSecondDateTime := "2017-06-01T16:27:18.543653798Z"
	secondDateTime, err := time.Parse(time.RFC3339, exampleSecondDateTime)
	if err != nil {
		test.Error(err)
	}

	createdSecondPrice, err := storage.Prices.CreatePrice(Price{Value: 124, DateTime: secondDateTime})
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Prices.DeletePrice(createdSecondPrice)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Products.AddPriceToProduct(createdProduct1.ID, createdSecondPrice.ID)
	if err != nil {
		test.Error(err)
	}

	foundedProductsForFirstPage, err := storage.Products.ReadProductsByNameWithPagination("тестовый", "ru", 1, 1)
	if err != nil {
		test.Error(err)
	}

	if len(foundedProductsForFirstPage.Products) != 1 {
		test.Fatalf(err.Error())
	}

	if len(foundedProductsForFirstPage.Products[0].Prices) != 2 {
		test.Fatalf(err.Error())
	}

	if foundedProductsForFirstPage.Products[0].Prices[0].Value != 124 {
		test.Fatalf(err.Error())
	}
}
