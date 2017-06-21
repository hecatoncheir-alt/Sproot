package engine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
)

// ErrCompanyNotExists means that the company is not in the database
var ErrCompanyNotExists = errors.New("Company not exists")

// ErrCompanyAlreadyExists means that the company is in the database already
var ErrCompanyAlreadyExists = errors.New("Company already exists")

// ErrCompanyCanNotBeDeleted delete all nodes with company predicates
var ErrCompanyCanNotBeDeleted = errors.New("Company can not be deleted")

// DeleteCompany method for delete all nodes with company name
func (engine *Engine) DeleteCompany(companyName string) error {
	var path *path.Path

	regCompanyName, err := regexp.Compile(strings.ToLower(companyName))
	if err != nil {
		return err
	}

	fmt.Println("In")
	path = cayley.StartPath(engine.Store).Regex(regCompanyName).InPredicates()
	path.Iterate(nil).EachValue(nil, func(value quad.Value) {
		fmt.Println(value.String())
	})

	fmt.Println("Out")
	path = cayley.StartPath(engine.Store).Regex(regCompanyName).OutPredicates()
	path.Iterate(nil).EachValue(nil, func(value quad.Value) {
		fmt.Println(value.String())
	})

	return nil
}

// GetCompany return company object of company node in graph store
func (engine *Engine) GetCompany(companyName string) (company *Company, err error) {
	var path *path.Path

	regCompanyName, err := regexp.Compile(strings.ToLower(companyName))
	if err != nil {
		return nil, err
	}

	path = cayley.StartPath(engine.Store).Regex(regCompanyName)

	var companyInStore string
	path.Iterate(nil).EachValue(nil, func(value quad.Value) {
		companyInStore = value.String()
	})

	if companyInStore == "" {
		return nil, ErrCompanyNotExists
	}

	// ---

	return nil, ErrCompanyNotExists
}

// SaveCompany method for add triplets to graph db
func (engine *Engine) SaveCompany(company *Company) (companyInStore *Company, err error) {
	_, err = engine.GetCompany(company.Name)
	if err != ErrCompanyNotExists {
		return nil, err
	}

	companyName := strings.ToLower(company.Name)
	companyAddTime := time.Now().String()

	transaction := cayley.NewTransaction()
	transaction.AddQuad(cayley.Quad(companyName, "is", "Company name", "Company"))
	transaction.AddQuad(cayley.Quad(companyName, "was added", companyAddTime, "Time"))
	transaction.AddQuad(cayley.Quad(companyName, "has", company.IRI, "Company link"))
	transaction.AddQuad(cayley.Quad(companyName, "has", "Address", "Company"))
	transaction.AddQuad(cayley.Quad(companyName, "has many", "Product", "Company"))
	err = engine.Store.ApplyTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = engine.SaveCategoriesOfCompany(company.Categories, company.Name)
	if err != nil {
		return nil, err
	}

	companyInStore, err = engine.GetCompany(company.Name)
	if err != nil {
		return nil, err
	}

	return companyInStore, nil
}
