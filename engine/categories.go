package engine

import (
	"strings"

	"github.com/cayleygraph/cayley"
)

// DeleteCategoriesOfCompany is method for delete categories from company
func (engine *Engine) DeleteCategoriesOfCompany(categories []string, companyName string) error {

	return nil
}

// SaveCategoriesOfCompany method for add categories to company
func (engine *Engine) SaveCategoriesOfCompany(categories []string, companyName string) error {
	var err error
	companyName = strings.ToLower(companyName)

	_, err = engine.GetCompany(companyName)
	if err != ErrCompanyNotExists {
		return err
	}

	// TODO: Нужно получить список категорий и добавлять только нужные
	for category := range categories {
		transaction := cayley.NewTransaction()
		transaction.AddQuad(cayley.Quad(category, "is", "Category name", "Category"))
		transaction.AddQuad(cayley.Quad(category, "belongs", companyName, "Category"))
		engine.Store.ApplyTransaction(transaction)
	}

	return nil
}
