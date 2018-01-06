package storage

import (
	"testing"
)

func TestIntegrationCategoryCanBeAddedToCompany(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

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
}

func TestIntegrationCategoryCanBeRemovedFromCompany(test *testing.T) {}
