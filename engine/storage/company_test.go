package storage

import (
	"fmt"
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
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	createdCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "Test category for company"})
	if err != nil || createdCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeactivateCategory(createdCategory)

	updatedCompany, err := createdCompany.AddCategory(createdCategory.ID)

	if updatedCompany.Categories[0].ID != createdCategory.ID {
		test.Fail()
	}

	if updatedCompany.Categories[0].Companies[0].ID != updatedCompany.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeRemovedFromCompany(test *testing.T) {
	test.Skip()
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	createdCompany, err := storage.Companies.CreateCompany(Company{Name: "Test company name"})
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	firstCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "First category for company"})
	if err != nil || firstCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeactivateCategory(firstCategory)

	updatedCompany, err := createdCompany.AddCategory(firstCategory.ID)

	fmt.Println(updatedCompany)

	secondCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "Second category for company"})
	if err != nil || secondCategory.ID == "" {
		test.Error(err)
	}

}
