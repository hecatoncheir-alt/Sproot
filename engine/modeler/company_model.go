package modeler

import "github.com/hecatoncheir/Sproot/engine/storage"

func setCompanyModel(store *storage.Storage) {
	companyForCreate:=storage.Company{

	}

	store.Companies.ReadCompaniesByName()

	store.Companies.CreateCompany()
}
