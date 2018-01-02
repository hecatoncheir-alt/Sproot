package storage

import (
	"testing"
	"fmt"
)

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	storage := New(databaseHost, databasePort)
	err := storage.SetUp()
	if err != nil {
		test.Fail()
	}
}

// TODO
func TestIntegrationStorageCanGetCompanyWithCategories(test *testing.T) {
	test.Skip()
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	createdCompanyForOneCategory, _ := storage.Companies.CreateCompany(Company{Name: "Company with categories"})
	fmt.Println(createdCompanyForOneCategory)

	//firstCategoryForCompany := Category{Name: "First category for company", Companies: []Company{createdCompanyForOneCategory}}
	//createdFirstCategoryForCompany, _ := storage.Categories.CreateCategory(firstCategoryForCompany)

	//updatedCompany, _ := createdCompanyForOneCategory.AddCategory(createdFirstCategoryForCompany.ID)
	//
	//if len(updatedCompany.Categories) > 1 {
	//	test.Fail()
	//}
	//
	//if updatedCompany.Categories[0].Companies[0].ID != createdCompanyForOneCategory.ID {
	//	test.Fail()
	//}
	//
	//secondCategoryForCompany := Category{Name: "Second category for company", Companies: []Company{createdCompanyForOneCategory}}
	//createdSecondCategoryForCompany, _ := storage.Categories.CreateCategory(secondCategoryForCompany)
	//
	//updatedCompany, _ = createdCompanyForOneCategory.AddCategory(createdSecondCategoryForCompany.ID)
	//
	//if len(updatedCompany.Categories) < 2 {
	//	test.Fail()
	//}
	//
	//if updatedCompany.Categories[1].Companies[0].ID != createdCompanyForOneCategory.ID {
	//	test.Fail()
	//}
	//
	////updatedCompanyWithOneCategory, _ := storage.Companies.UpdateCompany(createdCompanyForOneCategory)
	//fmt.Println(updatedCompany)

	storage.DeleteAll()
}
