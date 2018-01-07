package storage

import (
	"testing"
)

func TestIntegrationCategoryCanBeAddedToCompany(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()

	createdCompany, err := storage.Companies.CreateCompany(Company{Name: "Test company name"})

	defer storage.Companies.DeleteCompany(createdCompany)

	createdFirstCategory, err :=
			storage.Categories.CreateCategory(Category{Name: "First test category for company"})
	if err != nil || createdFirstCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdFirstCategory)

	updatedCompany, err := createdCompany.AddCategory(createdFirstCategory.ID)

	if updatedCompany.Categories[0].ID != createdFirstCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}

	createdSecondCategory, err :=
			storage.Categories.CreateCategory(Category{Name: "Second test category for company"})
	if err != nil || createdSecondCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdSecondCategory)

	updatedCompany, err = createdCompany.AddCategory(createdSecondCategory.ID)

	if updatedCompany.Categories[0].ID != createdFirstCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeRemovedFromCompany(test *testing.T) {}
