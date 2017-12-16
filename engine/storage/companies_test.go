package storage

import (
	"testing"
)

func TestIntegrationCompanyCanBeReadByName(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForSearch := Company{Name: "Test company"}

	companiesFromStore, err := storage.Companies.ReadCompaniesByName(companyForSearch.Name)
	if err != ErrCompaniesByNameNotFound {
		test.Fail()
	}

	if companiesFromStore != nil {
		test.Fail()
	}
	//
	//createdCategory, err := storage.Categories.CreateCategory(&categoryForSearch)
	//if err != nil {
	//	test.Error(err)
	//}
	//
	//defer storage.Categories.DeleteCategory(createdCategory)
	//
	//categoriesFromStore, err = storage.Categories.ReadCategoriesByName(createdCategory.Name)
	//if err != nil {
	//	test.Fail()
	//}
	//
	//if categoriesFromStore[0].Name != createdCategory.Name {
	//	test.Fail()
	//}
	//
	//if categoriesFromStore[0].ID == "" {
	//	test.Fail()
	//}
}

func TestIntegrationCompanyCanBeDeleted(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForTest := Company{Name: "Test company"}
	//createdCategory, err := storage.Companies.CreateCategory(&companyForTest)
	//if err != nil {
	//	test.Error(err)
	//}

	deletedCompanyID, err := storage.Companies.DeleteCompany(companyForTest)
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	if deletedCompanyID != companyForTest.ID {
		test.Fail()
	}

	//_, err = storage.Categories.ReadCategoryByID(deletedCategoryID)
	//if err != ErrCategoryDoesNotExist {
	//	test.Error(err)
	//}
}
