package storage

import (
	"testing"
)

func TestIntegrationCompanyCanBeCreated(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForTest := Company{Name: "Test company"}

	createdCompany, err := storage.Companies.CreateCompany(companyForTest)
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	if createdCompany.ID == "" {
		test.Fail()
	}

	if createdCompany.IsActive != true {
		test.Fail()
	}

	if createdCompany.Name != companyForTest.Name {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeReadByName(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForTest := Company{Name: "Test company"}

	companiesFromStore, err := storage.Companies.ReadCompaniesByName(companyForTest.Name)
	if err != ErrCompaniesByNameNotFound {
		test.Fail()
	}

	if companiesFromStore != nil {
		test.Fail()
	}

	createdCompany, err := storage.Companies.CreateCompany(companyForTest)
	if err != nil || createdCompany.ID == "" {
		test.Fail()
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	companiesFromStore, err = storage.Companies.ReadCompaniesByName(createdCompany.Name)
	if err != nil {
		test.Fail()
	}

	if companiesFromStore[0].Name != createdCompany.Name {
		test.Fail()
	}

	if companiesFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeReadById(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForSearch := Company{Name: "Test category"}

	companiesFromStore, err := storage.Companies.ReadCompaniesByName(companyForSearch.Name)
	if err != ErrCompaniesByNameNotFound {
		test.Fail()
	}

	companyFromStore, err := storage.Companies.ReadCompanyByID("0")
	if err != ErrCompanyDoesNotExist {
		test.Fail()
	}

	if companiesFromStore != nil {
		test.Fail()
	}

	createdCompany, err := storage.Companies.CreateCompany(companyForSearch)
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	companyFromStore, err = storage.Companies.ReadCompanyByID(createdCompany.ID)
	if err != nil {
		test.Fail()
	}

	if companyFromStore.Name != createdCompany.Name {
		test.Fail()
	}

	if companyFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeUpdated(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	updatedCompany, err := storage.Companies.UpdateCompany(Company{Name: "Updated test company"})
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	companytForCreate := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companytForCreate)
	if err != nil {
		test.Error(err)
	}

	defer storage.Companies.DeactivateCompany(createdCompany)

	companyForUpdate := Company{ID: createdCompany.ID, Name: "Updated test company", IsActive: createdCompany.IsActive}
	updatedCompany, err = storage.Companies.UpdateCompany(companyForUpdate)
	if err != nil {
		test.Error(err)
	}

	if updatedCompany.Name != companyForUpdate.Name {
		test.Fail()
	}

	companyInStore, err := storage.Companies.ReadCompanyByID(createdCompany.ID)
	if err != nil {
		test.Error(err)
	}

	if updatedCompany.Name != companyInStore.Name {
		test.Fail()
	}
}

func TestIntegrationCompanyCanBeDeactivate(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	companyForTest := Company{Name: "Test company"}
	createdCompany, err := storage.Companies.CreateCompany(companyForTest)
	if err != nil {
		test.Error(err)
	}

	deletedCompanyID, err := storage.Companies.DeactivateCompany(createdCompany)
	if err != nil {
		if err != ErrCompanyCanNotBeWithoutID {
			test.Error(err)
		}
	}

	if deletedCompanyID != companyForTest.ID {
		test.Fail()
	}

	_, err = storage.Companies.ReadCompanyByID(deletedCompanyID)
	if err != ErrCompanyDoesNotExist {
		test.Error(err)
	}
}
