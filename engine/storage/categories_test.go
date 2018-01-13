package storage

import (
	"testing"
)

func TestIntegrationCategoryCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	if createdCategory.ID == "" {
		test.Fail()
	}

	existCategory, err := storage.Categories.CreateCategory(categoryForCreate, "en")
	if err == nil || err != ErrCategoryAlreadyExist {
		test.Error(err)
	}

	if existCategory.ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeReadByName(test *testing.T) {
	once.Do(prepareStorage)

	categoryForSearch := Category{Name: "Test category"}

	categoriesFromStore, err := storage.Categories.ReadCategoriesByName(categoryForSearch.Name, ".")
	if err != ErrCategoriesByNameNotFound {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := storage.Categories.CreateCategory(categoryForSearch, "en")
	if err != nil || createdCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	categoriesFromStore, err = storage.Categories.ReadCategoriesByName(createdCategory.Name, ".")
	if err != nil {
		test.Fail()
	}

	if categoriesFromStore[0].Name != createdCategory.Name {
		test.Fail()
	}

	if categoriesFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)

	categoryForSearch := Category{Name: "Test category"}

	categoriesFromStore, err := storage.Categories.ReadCategoriesByName(categoryForSearch.Name, "en")
	if err != ErrCategoriesByNameNotFound {
		test.Fail()
	}

	categoryFromStore, err := storage.Categories.ReadCategoryByID("0", ".")
	if err != ErrCategoryDoesNotExist {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := storage.Categories.CreateCategory(categoryForSearch, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	categoryFromStore, err = storage.Categories.ReadCategoryByID(createdCategory.ID, ".")
	if err != nil {
		test.Fail()
	}

	if categoryFromStore.Name != createdCategory.Name {
		test.Fail()
	}

	if categoryFromStore.IsActive == false {
		test.Fail()
	}

	if categoryFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeUpdated(test *testing.T) {
	once.Do(prepareStorage)

	updatedCategory, err := storage.Categories.UpdateCategory(Category{Name: "Updated test category"})
	if err != nil {
		if err != ErrCategoryCanNotBeWithoutID {
			test.Error(err)
		}
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	categoryForUpdate := Category{
		ID:       createdCategory.ID,
		Name:     "Updated test category",
		IsActive: createdCategory.IsActive}

	updatedCategory, err = storage.Categories.UpdateCategory(categoryForUpdate)
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryForUpdate.Name {
		test.Fail()
	}

	categoryInStore, err := storage.Categories.ReadCategoryByID(createdCategory.ID, ".")
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryInStore.Name {
		test.Fail()
	}

	if updatedCategory.IsActive != categoryInStore.IsActive {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeDeactivate(test *testing.T) {
	once.Do(prepareStorage)

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	categoryInStore, err := storage.Categories.ReadCategoryByID(createdCategory.ID, ".")
	if err != nil {
		test.Error(err)
	}

	if categoryInStore.IsActive != true {
		test.Fail()
	}

	updatedCategoryID, err := storage.Categories.DeactivateCategory(categoryInStore)
	if err != nil {
		test.Error(err)
	}

	categoryInStore, err = storage.Categories.ReadCategoryByID(updatedCategoryID, ".")
	if err != nil {
		test.Error(err)
	}

	if categoryInStore.IsActive != false {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeDeleted(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	deletedCategoryID, err := storage.Categories.DeleteCategory(createdCategory)
	if err != nil {
		test.Error(err)
	}

	if deletedCategoryID != createdCategory.ID {
		test.Fail()
	}

	_, err = storage.Categories.ReadCategoryByID(deletedCategoryID, ".")
	if err != ErrCategoryDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationCompanyCanBeAddedToCategory(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")

	defer storage.Categories.DeleteCategory(createdCategory)

	createdFirstCompany, _ := storage.Companies.CreateCompany(Company{Name: "First test company for category"}, "en")

	defer storage.Companies.DeleteCompany(createdFirstCompany)

	err = storage.Categories.AddCompanyToCategory(createdCategory.ID, createdFirstCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCategory, _ := storage.Categories.ReadCategoryByID(createdCategory.ID, ".")

	if updatedCategory.Companies[0].ID != createdFirstCompany.ID {
		test.Fail()
	}

	if updatedCategory.Companies[0].Categories[0].ID != createdCategory.ID {
		test.Fail()
	}

	createdSecondCompany, _ := storage.Companies.CreateCompany(Company{Name: "Second test company for category"}, "en")

	defer storage.Companies.DeleteCompany(createdSecondCompany)

	err = storage.Categories.AddCompanyToCategory(createdCategory.ID, createdSecondCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCategory, _ = storage.Categories.ReadCategoryByID(createdCategory.ID, ".")

	if updatedCategory.Companies[0].ID != createdFirstCompany.ID {
		test.Fail()
	}

	if updatedCategory.Companies[0].Categories[0].ID != updatedCategory.ID {
		test.Fail()
	}

	if updatedCategory.Companies[1].ID != createdSecondCompany.ID {
		test.Fail()
	}

	if updatedCategory.Companies[1].Categories[0].ID != updatedCategory.ID {
		test.Fail()
	}
}
func TestIntegrationCompanyCanBeRemovedFromCategory(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCategory, _ :=
		storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")

	defer storage.Categories.DeleteCategory(createdCategory)

	createdFirstCompany, _ := storage.Companies.CreateCompany(Company{Name: "First test company for category"}, "en")

	defer storage.Companies.DeleteCompany(createdFirstCompany)

	storage.Categories.AddCompanyToCategory(createdCategory.ID, createdFirstCompany.ID)

	createdSecondCompany, _ := storage.Companies.CreateCompany(Company{Name: "Second test company for category"}, "en")

	defer storage.Companies.DeleteCompany(createdSecondCompany)

	storage.Categories.AddCompanyToCategory(createdCategory.ID, createdSecondCompany.ID)

	updatedCategory, _ := storage.Categories.ReadCategoryByID(createdCategory.ID, ".")

	if len(updatedCategory.Companies) != 2 {
		test.Fail()
	}

	if updatedCategory.Companies[0].ID != createdFirstCompany.ID {
		test.Fail()
	}

	err = storage.Categories.RemoveCompanyFromCategory(createdCategory.ID, createdFirstCompany.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCategory, _ = storage.Categories.ReadCategoryByID(createdCategory.ID, ".")

	if len(updatedCategory.Companies) != 1 {
		test.Fail()
	}

	if updatedCategory.Companies[0].ID != createdSecondCompany.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanHasNameWithManyLanguages(test *testing.T) {
	once.Do(prepareStorage)

	var err error

	createdCategory, _ := storage.Categories.CreateCategory(Category{Name: "Test category"}, "en")
	defer storage.Categories.DeleteCategory(createdCategory)

	err = storage.Categories.AddLanguageOfCategoryName(createdCategory.ID, "Тестовая категория", "ru")
	if err != nil {
		test.Fail()
	}

	categoryWithEnName, _ := storage.Categories.ReadCategoryByID(createdCategory.ID, "en")
	if categoryWithEnName.Name != "Test category" {
		test.Fail()
	}

	categoryWithRuName, _ := storage.Categories.ReadCategoryByID(createdCategory.ID, "ru")
	if categoryWithRuName.Name != "Тестовая категория" {
		test.Fail()
	}
}
