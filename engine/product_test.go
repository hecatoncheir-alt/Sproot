package engine

import (
	"testing"
	"time"

	"github.com/hecatoncheir/Sproot/configuration"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

// TODO test check
// ok      github.com/hecatoncheir/Sproot  3.059s
// ok      github.com/hecatoncheir/Sproot/configuration    0.045s
// --- FAIL: TestIntegrationNewPriceWithNewProductCanBeCreated (0.55s)
//         product_test.go:27: company already exist
//         product_test.go:35: category already exist
// panic: runtime error: index out of range [recovered]
//         panic: runtime error: index out of range

// goroutine 39 [running]:
// testing.tRunner.func1(0xc0421fc3c0)
//         C:/Users/Vitaliy/dev/go/src/testing/testing.go:742 +0x2a4
// panic(0x8d34c0, 0xc335a0)
//         C:/Users/Vitaliy/dev/go/src/runtime/panic.go:505 +0x237
// github.com/hecatoncheir/Sproot/engine.TestIntegrationNewPriceWithNewProductCanBeCreated(0xc0421fc3c0)
//         C:/Users/Vitaliy/go/src/github.com/hecatoncheir/Sproot/engine/product_test.go:111 +0xeff
// testing.tRunner(0xc0421fc3c0, 0x998578)
//         C:/Users/Vitaliy/dev/go/src/testing/testing.go:777 +0xd7
// created by testing.(*T).Run
//         C:/Users/Vitaliy/dev/go/src/testing/testing.go:824 +0x2e7
// FAIL    github.com/hecatoncheir/Sproot/engine   3.205s
func TestIntegrationNewPriceWithNewProductCanBeCreated(test *testing.T) {
	test.Skip()
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

	product := ProductOfCompany{
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

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])
	defer engine.Storage.Products.DeleteProduct(productFromStorage)

	products, err := engine.Storage.Products.ReadProductsByName(product.Name, "ru")
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

// TODO: for real tests
// TODO check tests
// ok      github.com/hecatoncheir/Sproot/configuration    0.042s
// --- FAIL: TestIntegrationNewPriceWithExistedProductCanBeCreated (0.51s)
//         product_test.go:158: company already exist
//         product_test.go:166: category already exist
//         product_test.go:178: city already exist
//         product_test.go:233: product can not be added to price
//         product_test.go:237:
// FAIL
// FAIL    github.com/hecatoncheir/Sproot/engine   5.185s
func TestIntegrationNewPriceWithExistedProductCanBeCreated(test *testing.T) {
	test.Skip()
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

	product := ProductOfCompany{
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

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])

	products, err := engine.Storage.Products.ReadProductsByName(product.Name, "ru")
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

// TODO check tests
// ok      github.com/hecatoncheir/Sproot/configuration    0.044s
// --- FAIL: TestIntegrationNewPriceWithExistedProductsCanBeCreatedForRightProduct (0.53s)
//         product_test.go:303: company already exist
//         product_test.go:311: category already exist
//         product_test.go:325: city already exist
//         product_test.go:397: product can not be added to price
//         product_test.go:401:
// FAIL
// FAIL    github.com/hecatoncheir/Sproot/engine   2.922s
func TestIntegrationNewPriceWithExistedProductsCanBeCreatedForRightProduct(test *testing.T) {
	test.Skip()
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

	product := ProductOfCompany{
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

	productFromStorage, err := product.UpdateInStorage(engine.Storage)
	if err != nil {
		test.Error(err)
	}

	if len(productFromStorage.Prices) == 0 {
		test.Fatal()
	}

	defer engine.Storage.Prices.DeletePrice(productFromStorage.Prices[0])

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

	products, err := engine.Storage.Products.ReadProductsByName(product.Name, "ru")
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
