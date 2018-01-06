package storage

import (
	"testing"
	"fmt"
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

// TODO last
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

	_, err = createdCompany.AddCategory(firstCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, err := storage.Companies.ReadCompanyByID(createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	if len(updatedCompany.Categories) != 1 {
		test.Fail()
	}

	if updatedCompany.Categories[0].ID != firstCategory.ID {
		test.Fail()
	}

	secondCategory, err :=
		storage.Categories.CreateCategory(Category{Name: "Second category for company"})
	if err != nil || secondCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeactivateCategory(secondCategory)

	_, err = updatedCompany.AddCategory(secondCategory.ID)
	if err != nil {
		test.Error(err)
	}

	updatedCompany, err = storage.Companies.ReadCompanyByID(createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	if len(updatedCompany.Categories) != 2 {
		test.Fail()
	}

	fmt.Println(updatedCompany)
}
