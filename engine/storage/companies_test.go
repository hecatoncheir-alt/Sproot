package storage

import (
	"encoding/json"
	"testing"
	"time"
)

func TestIntegrationCompanyCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	companyForTest := Company{Name: "Test company"}

	createdCompany, err := storage.Companies.CreateCompany(companyForTest, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	if createdCompany.ID == "" {
		test.Fail()
	}

	if !createdCompany.IsActive {
		test.Fail()
	}

	if createdCompany.Name != companyForTest.Name {
		test.Fail()
	}
}

// Must be run parallel with TestIntegrationNewPriceWithExistedProductCanBeCreated
func TestIntegrationCompaniesAllCanBeRead(test *testing.T) {
	test.Parallel()
	once.Do(prepareStorage)

	companyForTest := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companyForTest, "en")
	if err != nil {
		test.Fail()
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	otherCompanyForTest := Company{Name: "Other test company"}
	otherCreatedCompany, err := storage.Companies.CreateCompany(otherCompanyForTest, "en")
	if err != nil {
		test.Fail()
	}

	defer storage.Companies.DeleteCompany(otherCreatedCompany)

	companiesFromStore, err := storage.Companies.ReadAllCompanies("en")
	if err != nil {
		test.Fail()
	}

	if len(companiesFromStore) != 2 {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeReadByName(test *testing.T) {
	once.Do(prepareStorage)

	companyForTest := Company{Name: "Test company"}

	companiesFromStore, err := storage.Companies.ReadCompaniesByName(companyForTest.Name, ".")
	if err != ErrCompaniesByNameNotFound {
		test.Fail()
	}

	if companiesFromStore != nil {
		test.Fail()
	}

	createdCompany, err := storage.Companies.CreateCompany(companyForTest, "en")
	if err != nil || createdCompany.ID == "" {
		test.Fail()
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	companiesFromStore, err = storage.Companies.ReadCompaniesByName(createdCompany.Name, "en")
	if err != nil {
		test.Fail()
	}

	if len(companiesFromStore) > 1 {
		test.Fail()
	}

	if companiesFromStore[0].Name != createdCompany.Name {
		test.Fail()
	}

	if companiesFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)

	companyForSearch := Company{Name: "Test category"}

	companyFromStore, err := storage.Companies.ReadCompanyByID("0", ".")
	if err != ErrCompanyDoesNotExist {
		test.Fail()
	}

	createdCompany, err := storage.Companies.CreateCompany(companyForSearch, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	companyFromStore, err = storage.Companies.ReadCompanyByID(createdCompany.ID, ".")
	if err != nil {
		test.Fail()
	}

	if companyFromStore.Name != createdCompany.Name {
		test.Fail()
	}

	if companyFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeUpdated(test *testing.T) {
	once.Do(prepareStorage)

	updatedCompany, err := storage.Companies.UpdateCompany(Company{Name: "Updated test company"})
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	companyForCreate := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companyForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	companyForUpdate := Company{ID: createdCompany.ID, Name: "Updated test company", IsActive: createdCompany.IsActive}
	updatedCompany, err = storage.Companies.UpdateCompany(companyForUpdate)
	if err != nil {
		test.Error(err)
	}

	if updatedCompany.ID != companyForUpdate.ID {
		test.Fail()
	}

	if updatedCompany.Name != companyForUpdate.Name {
		test.Fail()
	}

	companyInStore, err := storage.Companies.ReadCompanyByID(createdCompany.ID, ".")
	if err != nil {
		test.Error(err)
	}

	if updatedCompany.Name != companyInStore.Name {
		test.Fail()
	}

	if updatedCompany.ID != companyInStore.ID {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeDeactivate(test *testing.T) {
	once.Do(prepareStorage)

	companyForTest := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companyForTest, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeleteCompany(createdCompany)

	deactivatedCompanyID, err := storage.Companies.DeactivateCompany(createdCompany)
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	if deactivatedCompanyID != createdCompany.ID {
		test.Fail()
	}

	deactivatedCompany, err := storage.Companies.ReadCompanyByID(deactivatedCompanyID, ".")
	if err != nil {
		test.Error(err)
	}

	if deactivatedCompany.IsActive {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	companyForTest := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companyForTest, "en")
	if err != nil {
		test.Error(err)
	}

	deletedCompanyID, err := storage.Companies.DeleteCompany(createdCompany)
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	if deletedCompanyID != createdCompany.ID {
		test.Fail()
	}

	_, err = storage.Companies.ReadCompanyByID(deletedCompanyID, ".")
	if err != ErrCompanyDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationCategoryCanBeAddedToCompany(test *testing.T) {
	once.Do(prepareStorage)

	var err error

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

	createdFirstCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "First test category for company"}, "en")

	defer storage.Categories.DeleteCategory(createdFirstCategory)

	err = storage.Companies.AddCategoryToCompany(createdCompany.ID, createdFirstCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if len(updatedCompany.Categories) < 1 || len(updatedCompany.Categories) > 1 {
		test.Fatal()
	}

	if updatedCompany.Categories[0].ID != createdFirstCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}

	createdSecondCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "Second test category for company"}, "en")

	defer storage.Categories.DeleteCategory(createdSecondCategory)

	err = storage.Companies.AddCategoryToCompany(createdCompany.ID, createdSecondCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, _ = storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if updatedCompany.Categories[0].ID != createdFirstCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}

	if updatedCompany.Categories[1].ID != createdSecondCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[1].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeRemovedFromCompany(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, ".")

	defer storage.Companies.DeleteCompany(createdCompany)

	createdFirstCategory, _ :=
		storage.Categories.CreateCategory(Category{Name: "First test category for company"}, ".")

	defer storage.Categories.DeleteCategory(createdFirstCategory)

	storage.Companies.AddCategoryToCompany(createdCompany.ID, createdFirstCategory.ID)

	createdSecondCategory, _ :=
		storage.Categories.CreateCategory(Category{Name: "Second test category for company"}, ".")

	defer storage.Categories.DeleteCategory(createdSecondCategory)

	storage.Companies.AddCategoryToCompany(createdCompany.ID, createdSecondCategory.ID)

	updatedCompany, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if len(updatedCompany.Categories) != 2 {
		test.Fail()
	}

	if updatedCompany.Categories[0].ID != createdFirstCategory.ID {
		test.Fail()
	}

	err = storage.Companies.RemoveCategoryFromCompany(createdCompany.ID, createdFirstCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, _ = storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if len(updatedCompany.Categories) != 1 {
		test.Fail()
	}

	if updatedCompany.Categories[0].ID != createdSecondCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCompanyCanHasNameWithManyLanguages(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(createdCompany)

	err = storage.Companies.AddLanguageOfCompanyName(createdCompany.ID, "Тестовая компания", "ru")
	if err != nil {
		test.Fail()
	}

	companyWithEnName, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, "en")
	if companyWithEnName.Name != "Test company" {
		test.Fail()
	}

	companyWithRuName, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, "ru")
	if companyWithRuName.Name != "Тестовая компания" {
		test.Fail()
	}
}

func TestIntegrationProductCanBeAddedToCompany(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")

	defer storage.Companies.DeleteCompany(createdCompany)

	createdCategory, _ :=
		storage.Categories.CreateCategory(Category{Name: "Test category for company"}, "en")

	defer storage.Categories.DeleteCategory(createdCategory)

	storage.Companies.AddCategoryToCompany(createdCompany.ID, createdCategory.ID)

	createdProductForCompany, err := storage.Products.CreateProduct(Product{Name: "Test product for company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Products.DeleteProduct(createdProductForCompany)

	storage.Categories.AddProductToCategory(createdCategory.ID, createdProductForCompany.ID)

	err = storage.Companies.AddProductToCompany(createdCompany.ID, createdProductForCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if len(updatedCompany.Categories) < 1 || len(updatedCompany.Categories) > 1 {
		test.Fatal()
	}

	if len(updatedCompany.Categories[0].Products) < 1 || len(updatedCompany.Categories[0].Products) > 1 {
		test.Fatal()
	}

	if updatedCompany.Categories[0].Products[0].ID != createdProductForCompany.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Products[0].Name != createdProductForCompany.Name {
		test.Fail()
	}

	createdProductForCategory, _ := storage.Products.CreateProduct(Product{Name: "Test product for category"}, "en")

	defer storage.Products.DeleteProduct(createdProductForCategory)

	storage.Categories.AddProductToCategory(createdCategory.ID, createdProductForCategory.ID)

	updatedCompany, err = storage.Companies.ReadCompanyByID(createdCompany.ID, ".")

	if len(updatedCompany.Categories[0].Products) < 1 || len(updatedCompany.Categories[0].Products) > 1 {
		test.Fatal()
	}

	if updatedCompany.Categories[0].Products[0].Name != createdProductForCompany.Name {
		test.Fail()
	}
}

func TestIntegrationCompaniesCanBeAddedFromExportedJSON(test *testing.T) {
	once.Do(prepareStorage)

	createdCategory, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(createdCategory)

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(createdCompany)

	storage.Companies.AddCategoryToCompany(createdCompany.ID, createdCategory.ID)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)

	storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 132.3, DateTime: priceData})
	defer storage.Prices.DeletePrice(createdPrice)

	storage.Products.AddPriceToProduct(createdProduct.ID, createdPrice.ID)

	createdCity, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(createdCity)

	storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)

	updatedCompany, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, "en")

	all := allExportedCompanies{Language: "en"}
	all.Companies = append(all.Companies, updatedCompany)

	exportedJSON, err := json.Marshal(all)
	if err != nil {
		test.Error(err)
	}

	storage.Categories.DeleteCategory(createdCategory)
	storage.Companies.DeleteCompany(createdCompany)
	storage.Products.DeleteProduct(createdProduct)
	storage.Prices.DeletePrice(createdPrice)
	storage.Cities.DeleteCity(createdCity)

	_, err = storage.Companies.ReadCompanyByID(createdCompany.ID, "en")
	if err != ErrCompanyDoesNotExist {
		test.Error(err)
	}

	err = storage.Companies.ImportJSON(exportedJSON)

	importedCompany, _ := storage.Companies.ReadCompanyByID(createdCompany.ID, "en")

	if importedCompany.Name != createdCompany.Name {
		test.Fail()
	}

	for _, importedCategory := range importedCompany.Categories {
		if importedCategory.Name != createdCategory.Name {
			test.Fail()
		}

		for _, importedProduct := range importedCategory.Products {
			if importedProduct.Name != createdProduct.Name {
				test.Fail()
			}

			for _, importedPrice := range importedProduct.Prices {
				if importedPrice.Value != createdPrice.Value {
					test.Fail()
				}

				for _, importedCity := range importedPrice.Cities {
					if importedCity.Name != createdCity.Name {
						test.Fail()
					}
				}
			}
		}
	}
}

func TestIntegrationCompaniesCanBeExportedToJSON(test *testing.T) {
	once.Do(prepareStorage)

	createdCategory, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(createdCategory)

	createdCompany, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(createdCompany)

	storage.Companies.AddCategoryToCompany(createdCompany.ID, createdCategory.ID)

	createdProduct, _ := storage.Products.CreateProduct(Product{Name: "Test product"}, "en")
	defer storage.Products.DeleteProduct(createdProduct)

	storage.Products.AddCompanyToProduct(createdProduct.ID, createdCompany.ID)

	storage.Products.AddCategoryToProduct(createdProduct.ID, createdCategory.ID)

	createdOtherProduct, _ := storage.Products.CreateProduct(Product{Name: "Other test product"}, "en")
	defer storage.Products.DeleteProduct(createdOtherProduct)
	storage.Products.AddCategoryToProduct(createdOtherProduct.ID, createdCategory.ID)

	exampleDateTime := "2017-05-01T16:27:18.543653798Z"
	priceData, _ := time.Parse(time.RFC3339, exampleDateTime)
	createdPrice, _ := storage.Prices.CreatePrice(Price{Value: 132.3, DateTime: priceData})
	defer storage.Prices.DeletePrice(createdPrice)

	storage.Products.AddPriceToProduct(createdProduct.ID, createdPrice.ID)

	createdCity, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(createdCity)

	storage.Prices.AddCityToPrice(createdPrice.ID, createdCity.ID)

	exportedJSON, err := storage.Companies.ExportJSON("en")
	if err != nil {
		test.Error(err)
	}

	exportedCompanies := allExportedCompanies{}

	err = json.Unmarshal(exportedJSON, &exportedCompanies)
	if err != nil {
		test.Error(err)
	}

	if exportedCompanies.Language != "en" {
		test.Fail()
	}

	var exportedCompany Company

	for _, company := range exportedCompanies.Companies {
		if company.ID == createdCompany.ID {
			exportedCompany = company
		}
	}

	if exportedCompany.Name != createdCompany.Name {
		test.Fail()
	}

	if len(exportedCompany.Categories) == 0 {
		test.Fatalf("Expected company has one category, but actual 0")
	}

	if exportedCompany.Categories[0].Name != createdCategory.Name {
		test.Fail()
	}

	if len(exportedCompany.Categories[0].Products) > 1 {
		test.Fatal()
	}

	if len(exportedCompany.Categories[0].Products) == 0 {
		test.Fatalf("Expected category of company has one product, but actual 0")
	}

	if exportedCompany.Categories[0].Products[0].Name != createdProduct.Name {
		test.Fail()
	}

	if exportedCompany.Categories[0].Products[0].Prices[0].Value != createdPrice.Value {
		test.Fail()
	}

	if exportedCompany.Categories[0].Products[0].Prices[0].Cities[0].Name != createdCity.Name {
		test.Fail()
	}
}
