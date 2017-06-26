package engine

import (
	"fmt"
	"strings"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/iterator"
	"github.com/cayleygraph/cayley/quad"
)

// GetCategoriesOfCompany is a method for get all categories of company
func (engine *Engine) GetCategoriesOfCompany(companyName string) (categories []string, err error) {
	// var err error
	companyName = strings.ToLower(companyName)
	// it := iterator.NewAnd(engine.Store,
	// 	engine.Store.QuadIterator(quad.Object, engine.Store.ValueOf(quad.String(companyName))),
	// 	engine.Store.QuadIterator(quad.Predicate, engine.Store.ValueOf(quad.String("belongs"))))

	// defer it.Close()

	// for it.Next() {
	// 	f := engine.Store.Quad(it.Result()).String()
	// 	fmt.Println(f)
	// }

	path := cayley.StartPath(engine.Store, quad.String(companyName)).LabelContext("Category").In("belongs")

	path.Iterate(nil).EachValue(engine.Store, func(value quad.Value) {
		categories = append(categories, value.String())
	})

	return categories, nil
}

// DeleteCategoriesOfCompany is method for delete categories from company
func (engine *Engine) DeleteCategoriesOfCompany(categories []string, companyName string) error {
	var err error

	companyName = strings.ToLower(companyName)
	c, _ := engine.GetCategoriesOfCompany(companyName)
	fmt.Println(c)

	for _, category := range categories {
		it := iterator.NewAnd(engine.Store,
			// engine.Store.QuadIterator(quad.Predicate, engine.Store.ValueOf(quad.String("belongs"))),
			engine.Store.QuadIterator(quad.Subject, engine.Store.ValueOf(quad.String(category))))

		defer it.Close()
		fmt.Println(it.Next())

		for it.Next() {
			f := engine.Store.Quad(it.Result()).String()
			fmt.Println(f)
			// store.RemoveQuad(store.Quad(it.Result()))
			// fmt.Println("removed")
		}

		// path := cayley.StartPath(engine.Store, quad.String(companyName)).LabelContext("Category").In("belongs")

		// path.Iterate(nil).EachValue(engine.Store, func(value quad.Value) {
		// 	cat := value.String()
		// 	if cat == category {
		// 		value := engine.Store.ValueOf(value)
		// 		da := engine.Store.Quad(value)
		// 		fmt.Println("aaaaaaaa")
		// 		fmt.Println(da)
		// 		// err = engine.Store.RemoveQuad(engine.Store.Quad(engine.Store.ValueOf(value)))
		// 	}
		// })
	}
	c, _ = engine.GetCategoriesOfCompany(companyName)
	fmt.Println(c)

	// sort.Strings()

	return err
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
	for _, category := range categories {
		transaction := cayley.NewTransaction()
		transaction.AddQuad(cayley.Quad(category, "is", "Category name", "Category"))
		transaction.AddQuad(cayley.Quad(category, "belongs", companyName, "Category"))
		engine.Store.ApplyTransaction(transaction)
	}

	return nil
}
