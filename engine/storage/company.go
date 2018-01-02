package storage

import (
	"log"
)

// Company is a structure of Categories in database
type Company struct {
	storage *Storage

	ID         string     `json:"uid,omitempty"`
	IRI        string     `json:"iri, omitempty"`
	Name       string     `json:"name,omitempty"`
	Categories []Category `json:"categories, omitempty"`
	IsActive   bool       `json:"isActive, omitempty"`
}

// AddCategory append company to category and category to company
func (company *Company) AddCategory(categoryID string) (Company, error) {
	categoryForApply, err := company.storage.Categories.ReadCategoryByID(categoryID)
	if err != nil {
		log.Println(err)
		return *company, err
	}

	categoryForApply.Companies = append(categoryForApply.Companies, *company)
	updatedCategory, err := company.storage.Categories.UpdateCategory(categoryForApply)
	if err != nil {
		log.Println(err)
		return *company, err
	}

	company.Categories = append(company.Categories, updatedCategory)

	updateCompany, err := company.storage.Companies.UpdateCompany(*company)
	if err != nil {
		log.Println(err)
		return *company, err
	}

	updatedCompany, err := company.storage.Companies.ReadCompanyByID(updateCompany.ID)
	if err != nil {
		log.Println(err)
		return *company, err
	}

	return updatedCompany, nil
}

// TODO
func (company *Company) RemoveCategory(categoryID string) (Company, error) {
	var indexOfCategoryForRemove int

	for i, category := range company.Categories {
		if category.ID == categoryID {
			indexOfCategoryForRemove = i
		}
	}

	company.Categories = append(company.Categories[:indexOfCategoryForRemove], company.Categories[indexOfCategoryForRemove+1:]...)

	updatedCompany, _ := company.storage.Companies.UpdateCompany(*company)
	return updatedCompany, nil
}
