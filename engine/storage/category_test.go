package storage

import (
	"testing"
)

func TestIntegrationCompanyCanBeAddedToCategory(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()

	createdCompany, err := storage.Companies.CreateCompany(Company{Name: "Test company"})

	defer storage.Companies.DeleteCompany(createdCompany)

	createdCategory, err :=
			storage.Categories.CreateCategory(Category{Name: "Test category for company"})
	if err != nil || createdCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	updatedCategory, err := createdCategory.AddCompany(createdCompany.ID)

	if updatedCategory.Companies[0].ID != createdCompany.ID {
		test.Fail()
	}

	if updatedCategory.Companies[0].Categories[0].ID != updatedCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeRemovedFromCategory(test *testing.T) {}
