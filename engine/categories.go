package engine

import (
	"fmt"
	"strings"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
)

// GetCategoriesOfCompany is a method for get all categories of company
func (engine *Engine) GetCategoriesOfCompany(companyName string) (categories []string, err error) {
	// var err error
	companyName = strings.ToLower(companyName)
	path := cayley.StartPath(engine.Store, quad.String(companyName)).LabelContext("Category").In("belongs").Out("is")

	path.Iterate(nil).EachValue(engine.Store, func(value quad.Value) {
		fmt.Println(value.String())
		// categories = append(categories, value.String())
	})

	return categories, nil
}

// DeleteCategoriesOfCompany is method for delete categories from company
func (engine *Engine) DeleteCategoriesOfCompany(categories []string, companyName string) error {

	// path := cayley.StartPath(engine.Store, quad.String(companyName))
	categoriesOfCompany, _ := engine.GetCategoriesOfCompany(companyName)
	fmt.Println(categoriesOfCompany)
	// sort.Strings()

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
