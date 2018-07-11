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

	defer func() {
		_, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

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

	defer func() {
		_, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
		if err != nil {
			test.Fail()
		}
	}()

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

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)

		if err != nil {
			test.Fail()
		}
	}()

	if instruction.ID == "" {
		test.Fail()
	}

	if !instruction.IsActive {
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

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Fail()
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

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

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	city, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(city)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Fail()
	}

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

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Fail()
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	city, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(city)
		if err != nil {
			test.Fail()
		}
	}()

	err = storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Fail()
	}

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

	updatedInstruction, err = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Fail()
	}

	if len(updatedInstruction.Cities) > 0 {
		test.Fatal()
	}
}

func TestIntegrationPageInstructionCanBeAddedToInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Fail()
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

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

	defer func() {
		_, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Fail()
	}

	if len(updatedInstruction.PagesInstruction) != 1 {
		test.Fatal()
	}

	if updatedInstruction.PagesInstruction[0].ID != createdPageInstruction.ID {
		test.Fail()
	}
}

func TestIntegrationPageInstructionCanBeRemovedFromInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	mVideoPageInstruction := PageInstruction{
		Path: "smartfony-i-svyaz/smartfony-205",
		PageInPaginationSelector: ".pagination-list .pagination-item",
		PageParamPath:            "/f/page=",
		CityParamPath:            "?cityId=",
		ItemSelector:             ".grid-view .product-tile",
		NameOfItemSelector:       ".product-tile-title",
		PriceOfItemSelector:      ".product-price-current"}

	createdPageInstruction, _ := storage.Instructions.CreatePageInstruction(mVideoPageInstruction)
	defer func() {
		_, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Error(err)
	}

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

	updatedInstruction, err = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(updatedInstruction.PagesInstruction) > 0 {
		test.Fatal()
	}
}

func TestIntegrationCategoryCanBeAddedToInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	category, err := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Categories.DeleteCategory(category)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(updatedInstruction.Categories) != 1 {
		test.Fatal()
	}

	if updatedInstruction.Categories[0].ID != category.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeRemovedFromInstruction(test *testing.T) {
	once.Do(prepareStorage)

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	category, err := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Categories.DeleteCategory(category)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	updatedInstruction, err := storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Error(err)
	}

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

	updatedInstruction, err = storage.Instructions.ReadInstructionByID(instruction.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(updatedInstruction.Categories) > 0 {
		test.Fatal()
	}
}

func TestIntegrationCanGetAllInstructionsOfCompany(test *testing.T) {
	once.Do(prepareStorage)

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	anotherCompany, err := storage.Companies.CreateCompany(Company{Name: "Another test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(anotherCompany)
		if err != nil {
			test.Error(err)
		}
	}()

	anotherInstruction, err := storage.Instructions.CreateInstructionForCompany(anotherCompany.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(anotherInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

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

	company, err := storage.Companies.CreateCompany(Company{Name: "Test company"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Companies.DeleteCompany(company)
		if err != nil {
			test.Error(err)
		}
	}()

	instruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(instruction)
		if err != nil {
			test.Error(err)
		}
	}()

	category, err := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Categories.DeleteCategory(category)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Companies.AddCategoryToCompany(company.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	err = storage.Instructions.AddCategoryToInstruction(instruction.ID, category.ID)
	if err != nil {
		test.Error(err)
	}

	city, err := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(city)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddCityToInstruction(instruction.ID, city.ID)
	if err != nil {
		test.Error(err)
	}

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
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeletePageInstruction(createdPageInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddPageInstructionToInstruction(instruction.ID, createdPageInstruction.ID)
	if err != nil {
		test.Error(err)
	}

	instructionForCompany, err := storage.Instructions.ReadAllInstructionsForCompany(company.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if len(instructionForCompany) != 1 {
		test.Fail()
	}

	err = storage.Companies.AddLanguageOfCompanyName(company.ID, "Тестовая компания", "ru")
	if err != nil {
		test.Error(err)
	}

	secondInstruction, err := storage.Instructions.CreateInstructionForCompany(company.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Instructions.DeleteInstruction(secondInstruction)
		if err != nil {
			test.Error(err)
		}
	}()

	secondCity, err := storage.Cities.CreateCity(City{Name: "Тестовый город"}, "ru")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(secondCity)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Instructions.AddCityToInstruction(secondInstruction.ID, secondCity.ID)
	if err != nil {
		test.Error(err)
	}

	instructionForCompany, err = storage.Instructions.ReadAllInstructionsForCompany(company.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	if len(instructionForCompany) != 2 {
		test.Fail()
	}
}
