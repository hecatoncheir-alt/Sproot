package storage

import "log"

// Category is a structure of Categories in database
type Category struct {
	storage *Storage

	ID        string    `json:"uid,omitempty"`
	Name      string    `json:"name,omitempty"`
	IsActive  bool      `json:"isActive, omitempty"`
	Companies []Company `json:"companies, omitempty"`
}

// AddCompany method update category companies property in storage
func (category *Category) AddCompany(companyID string) (Category, error) {
	companyForApply, err := category.storage.Companies.ReadCompanyByID(companyID)
	if err != nil {
		log.Println(err)
		return *category, err
	}

	companyForApply.Categories = append(companyForApply.Categories, *category)
	updatedCompany, err := category.storage.Companies.UpdateCompany(companyForApply)
	if err != nil {
		log.Println(err)
		return *category, err
	}

	category.Companies = append(category.Companies, updatedCompany)

	updateCategory, err := category.storage.Categories.UpdateCategory(*category)
	if err != nil {
		log.Println(err)
		return *category, err
	}

	updatedCategory, err := category.storage.Categories.ReadCategoryByID(updateCategory.ID)
	if err != nil {
		log.Println(err)
		return *category, err
	}

	return updatedCategory, nil
}

// TODO
func (category *Category) RemoveCompany(companyID string) (*Category, error) {

	return category, nil
}
