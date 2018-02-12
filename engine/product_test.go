package engine

import (
	"encoding/json"
	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
	"testing"
	"time"
)

func TestIntegrationNewPriceWithNewProductCanBeCreated(test *testing.T) {
	engine := New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	err = engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	companyForTest := storage.Company{Name: "М.Видео", IRI: "http://www.mvideo.ru/"}
	createdCompany, err := engine.Storage.Companies.CreateCompany(companyForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Companies.DeleteCompany(createdCompany)

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := engine.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Categories.DeleteCategory(createdCategory)

	err = engine.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := engine.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Cities.DeleteCity(createdCity)

	parseTime, err := time.Parse(time.RFC3339, "2018-02-10T08:34:35.6055814Z")
	if err != nil {
		test.Error(err)
	}

	testProductOfCompany := ProductOfCompany{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg",
		Price: PriceOfProduct{
			Value:    "46990",
			DateTime: parseTime,
			City: CityData{
				ID:   createdCity.ID,
				Name: createdCity.Name},
		},
		Language: "ru",
		Company: CompanyData{
			ID:   createdCompany.ID,
			IRI:  createdCompany.IRI,
			Name: createdCompany.Name},
		Category: CategoryData{
			ID:   createdCategory.ID,
			Name: createdCategory.Name},
	}

	productWithPriceJSON, err := json.Marshal(testProductOfCompany)
	if err != nil {
		test.Error(err)
	}

	product := ProductOfCompany{}
	json.Unmarshal([]byte(productWithPriceJSON), &product)

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])
	defer engine.Storage.Products.DeleteProduct(productFromStorage)

	products, err := engine.Storage.Products.ReadProductsByName(testProductOfCompany.Name, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(products) != 1 {
		test.Fatal()
	}

	if len(products[0].Prices) != 1 {
		test.Fatal()
	}

	if products[0].Prices[0].Value != 46990 {
		test.Fail()
	}

	if products[0].Prices[0].Cities[0].ID != createdCity.ID {
		test.Fail()
	}

	if products[0].Companies[0].ID != createdCompany.ID {
		test.Fail()
	}

	if products[0].Categories[0].ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationNewPriceWithExistedProductCanBeCreated(test *testing.T) {
	engine := New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	err = engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	companyForTest := storage.Company{Name: "М.Видео", IRI: "http://www.mvideo.ru/"}
	createdCompany, err := engine.Storage.Companies.CreateCompany(companyForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Companies.DeleteCompany(createdCompany)

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := engine.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Categories.DeleteCategory(createdCategory)

	err = engine.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	createdCity, err := engine.Storage.Cities.CreateCity(storage.City{Name: "Москва"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Cities.DeleteCity(createdCity)

	productForTest := storage.Product{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg"}

	createdProduct, err := engine.Storage.Products.CreateProduct(productForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Products.DeleteProduct(createdProduct)

	err = engine.Storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	err = engine.Storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	parseTime, err := time.Parse(time.RFC3339, "2018-02-10T08:34:35.6055814Z")
	if err != nil {
		test.Error(err)
	}

	testProductOfCompany := ProductOfCompany{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg",
		Price: PriceOfProduct{
			Value:    "46990",
			DateTime: parseTime,
			City: CityData{
				ID:   createdCity.ID,
				Name: createdCity.Name},
		},
		Language: "ru",
		Company: CompanyData{
			ID:   createdCompany.ID,
			IRI:  createdCompany.IRI,
			Name: createdCompany.Name},
		Category: CategoryData{
			ID:   createdCategory.ID,
			Name: createdCategory.Name},
	}

	productWithPriceJSON, err := json.Marshal(testProductOfCompany)
	if err != nil {
		test.Error(err)
	}

	product := ProductOfCompany{}
	json.Unmarshal([]byte(productWithPriceJSON), &product)

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])
	defer engine.Storage.Products.DeleteProduct(productFromStorage)

	products, err := engine.Storage.Products.ReadProductsByName(testProductOfCompany.Name, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(products) != 1 {
		test.Fail()
	}

	if products[0].ID != createdProduct.ID {
		test.Fail()
	}

	if len(products[0].Prices) != 1 {
		test.Fatal()
	}

	if products[0].Prices[0].Value != 46990 {
		test.Fail()
	}

	if products[0].Prices[0].Cities[0].ID != createdCity.ID {
		test.Fail()
	}

	if products[0].Companies[0].ID != createdCompany.ID {
		test.Fail()
	}

	if products[0].Categories[0].ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationNewPriceWithExistedProductsCanBeCreatedForRightProduct(test *testing.T) {
	engine := New()

	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	err = engine.SetUpStorage(config.Development.Database.Host, config.Development.Database.Port)
	if err != nil {
		test.Error(err)
	}

	companyForTest := storage.Company{Name: "М.Видео", IRI: "http://www.mvideo.ru/"}
	createdCompany, err := engine.Storage.Companies.CreateCompany(companyForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Companies.DeleteCompany(createdCompany)

	categoryForTest := storage.Category{Name: "Смартфоны"}
	createdCategory, err := engine.Storage.Categories.CreateCategory(categoryForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Categories.DeleteCategory(createdCategory)

	err = engine.Storage.Categories.AddCompanyToCategory(createdCategory.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	cityForTest := storage.City{Name: "Москва"}

	createdCity, err := engine.Storage.Cities.CreateCity(cityForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Cities.DeleteCity(createdCity)

	productForTest := storage.Product{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg"}

	createdProduct, err := engine.Storage.Products.CreateProduct(productForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Products.DeleteProduct(createdProduct)

	err = engine.Storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	err = engine.Storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	otherProductForTest := storage.Product{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Белый бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg"}

	otherCreatedProduct, err := engine.Storage.Products.CreateProduct(otherProductForTest, "ru")
	if err != nil {
		test.Error(err)
	}

	err = engine.Storage.Products.AddCategoryToProduct(otherCreatedProduct.ID, createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	defer engine.Storage.Products.DeleteProduct(otherCreatedProduct)

	parseTime, err := time.Parse(time.RFC3339, "2018-02-10T08:34:35.6055814Z")
	if err != nil {
		test.Error(err)
	}

	testProductOfCompany := ProductOfCompany{
		Name:             "Смартфон Samsung Galaxy S8 64Gb Черный бриллиант",
		IRI:              "http://www.mvideo.ru//products/smartfon-samsung-galaxy-s8-64gb-chernyi-brilliant-30027818",
		PreviewImageLink: "img.mvideo.ru/Pdb/30027818m.jpg",
		Price: PriceOfProduct{
			Value:    "46990",
			DateTime: parseTime,
			City: CityData{
				ID:   createdCity.ID,
				Name: createdCity.Name},
		},
		Language: "ru",
		Company: CompanyData{
			ID:   createdCompany.ID,
			IRI:  createdCompany.IRI,
			Name: createdCompany.Name},
		Category: CategoryData{
			ID:   createdCategory.ID,
			Name: createdCategory.Name},
	}

	productWithPriceJSON, err := json.Marshal(testProductOfCompany)
	if err != nil {
		test.Error(err)
	}

	product := ProductOfCompany{}
	json.Unmarshal([]byte(productWithPriceJSON), &product)

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])
	defer engine.Storage.Products.DeleteProduct(productFromStorage)

	categoryWithProducts, err := engine.Storage.Categories.ReadCategoryByID(createdCategory.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(categoryWithProducts.Products) != 2 {
		test.Fail()
	}

	productWithoutPrice, err := engine.Storage.Products.ReadProductByID(otherCreatedProduct.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(productWithoutPrice.Prices) > 0 {
		test.Fail()
	}

	products, err := engine.Storage.Products.ReadProductsByName(testProductOfCompany.Name, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(products) != 1 {
		test.Fail()
	}

	if products[0].ID != createdProduct.ID {
		test.Fail()
	}

	if len(products[0].Prices) != 1 {
		test.Fatal()
	}

	if products[0].Prices[0].Value != 46990 {
		test.Fail()
	}

	if products[0].Prices[0].Cities[0].ID != createdCity.ID {
		test.Fail()
	}

	if products[0].Companies[0].ID != createdCompany.ID {
		test.Fail()
	}

	if products[0].Categories[0].ID != createdCategory.ID {
		test.Fail()
	}

	_, err = engine.Storage.Cities.DeleteCity(createdCity)
	if err != nil {
		test.Error(err)
	}
}
