package storage

import (
	"testing"
)

func TestIntegrationPageInstructionCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	if createdPageInstruction.ID == "" {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)
	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	pageInstructionFromStore, err := storage.Instructions.ReadPageInstructionByID(createdPageInstruction.ID)
	if err != nil {
		test.Fail()
	}

	if createdPageInstruction.ID != pageInstructionFromStore.ID {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, err := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	if err != nil {
		test.Fail()
	}

	deletedPageInstructionID, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
	if err != nil {
		test.Error(err)
	}

	if deletedPageInstructionID != createdPageInstruction.ID {
		test.Fail()
	}

	_, err = storage.Instructions.ReadPageInstructionByID(createdPageInstruction.ID)
	if err != ErrPageInstructionDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationInstructionCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	defer storage.Instructions.DeleteInstruction(instruction)
	if err != nil {
		test.Fail()
	}

	if instruction.ID == "" {
		test.Fail()
	}

	if instruction.IsActive != true {
		test.Fail()
	}

	if len(instruction.Companies) != 1 {
		test.Fatal()
	}

	if instruction.Companies[0].ID != company.ID {
		test.Fail()
	}
}

func TestIntegrationInstructionCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")

	deletedInstructionID, err := storage.Instructions.DeleteInstruction(instruction)
	if err != nil {
		test.Error(err)
	}

	if deletedInstructionID != instruction.ID {
		test.Fail()
	}

	_, err = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != ErrInstructionDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationCityCanBeAddedToInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	city, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(city)

	err = storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.Cities) != 1 {
		test.Fatal()
	}

	if updatedInstruction.Cities[0].ID != city.ID {
		test.Fail()
	}

	if updatedInstruction.Cities[0].Name != "Test city" {
		test.Fail()
	}
}

func TestIntegrationCityCanBeRemovedFromInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	city, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(city)

	err = storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if updatedInstruction.Cities[0].ID != city.ID {
		test.Fail()
	}

	if updatedInstruction.Cities[0].Name != "Test city" {
		test.Fail()
	}

	err = storage.Instructions.RemoveCityFromInstruction(instruction.ID, city.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if len(updatedInstruction.Cities) > 0 {
		test.Fatal()
	}
}

func TestIntegrationPageInstructionCanBeAddedToInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, _ := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	err = storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.PagesInstruction) != 1 {
		test.Fatal()
	}

	if updatedInstruction.PagesInstruction[0].ID != createdPageInstruction.ID {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeRemovedFromInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, _ := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	err = storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.PagesInstruction) != 1 {
		test.Fatal()
	}

	if updatedInstruction.PagesInstruction[0].ID != createdPageInstruction.ID {
		test.Fail()
	}

	err = storage.Instructions.RemovePageInstructionFromInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ = storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.PagesInstruction) > 0 {
		test.Fatal()
	}
}

func TestIntegrationCategoryCanBeAddedToInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	category, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(category)

	err = storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.Categories) != 1 {
		test.Fatal()
	}

	if updatedInstruction.Categories[0].ID != category.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeRemovedFromInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Instructions.DeleteInstruction(instruction)

	category, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(category)

	storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)

	updatedInstruction, _ := storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.Categories) != 1 {
		test.Fatal()
	}

	if updatedInstruction.Categories[0].ID != category.ID {
		test.Fail()
	}

	err = storage.Instructions.RemoveCategoryFromInstruction(instruction.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, _ = storage.Instructions.ReadInstructionByID(instruction.ID, "en")

	if len(updatedInstruction.Categories) > 0 {
		test.Fatal()
	}
}

func TestIntegrationCanGetAllInstructionsOfCompany(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, _ := storage.Instructions.CreateInstructionForCompany(company.ID, "en")

	defer storage.Instructions.DeleteInstruction(instruction)

	anotherCompany, _ := storage.Companies.CreateCompany(Company{Name: "Another test company"}, "en")
	defer storage.Companies.DeleteCompany(anotherCompany)

	anotherInstruction, _ := storage.Instructions.CreateInstructionForCompany(anotherCompany.ID, "en")
	defer storage.Instructions.DeleteInstruction(anotherInstruction)

	instructionsForCompany, err := storage.Instructions.ReadAllInstructionsForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(instructionsForCompany) != 1 {
		test.Fail()
	}

	if instructionsForCompany[0].Companies[0].ID != company.ID {
		test.Fail()
	}
}

func TestIntegrationCanGetFullInstructionForCompany(test *testing.T) {
	once.Do(prepareStorage)

	company, _ := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	defer storage.Companies.DeleteCompany(company)

	instruction, _ := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	defer storage.Instructions.DeleteInstruction(instruction)

	category, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(category)

	storage.Companies.AddCategoryToCompany(company.ID, category.ID)

	storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)

	city, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(city)

	storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, _ := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	defer storage.Instructions.DeletePageInstruction(createdPageInstruction)

	storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)

	instructionForCompany, err := storage.Instructions.ReadAllInstructionsForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(instructionForCompany) != 1 {
		test.Fail()
	}

	storage.Companies.AddLanguageOfCompanyName(company.ID, "Тестовая компания", "ru")

	secondInstruction, _ := storage.Instructions.CreateInstructionForCompany(company.ID, "ru")
	defer storage.Instructions.DeleteInstruction(secondInstruction)

	secondCity, _ := storage.Cities.CreateCity(City{Name: "Тестовый город"}, "ru")
	defer storage.Cities.DeleteCity(secondCity)

	storage.Instructions.AddCityToInstruction(secondInstruction.ID, secondCity.ID)

	instructionForCompany, err = storage.Instructions.ReadAllInstructionsForCompany(company.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(instructionForCompany) != 2 {
		test.Fail()
	}
}
