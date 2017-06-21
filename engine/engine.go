package engine

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"

	// sql driver
	_ "github.com/cayleygraph/cayley/graph/sql"
)

// ErrCompanyNotExists means that the company is not in the database
var ErrCompanyNotExists = errors.New("Company not exists")

// ErrCompanyAlreadyExists means that the company is in the database already
var ErrCompanyAlreadyExists = errors.New("Company already exists")

// Engine is a main object of engine pkg
type Engine struct {
	Store       *cayley.Handle
	SQLDataBase *sql.DB
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// DatabaseSetUp is a method for setup SQL database for graph engine
func (engine *Engine) DatabaseSetUp(user, host string, port int, ssl string, baseName string) (*sql.DB, *cayley.Handle, error) {
	var err error
	baseType := "postgres"

	dbAddr := fmt.Sprintf("%vql://%v@%v:%v?sslmode=%v", baseType, user, host, port, ssl)

	db, err := sql.Open(baseType, dbAddr)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	engine.SQLDataBase = db

	createDbQuery := fmt.Sprintf("create database if not exists %v", baseName)

	_, err = db.Query(createDbQuery)
	if err != nil {
		log.Fatal(err)
		return db, nil, err
	}

	tableAddr := fmt.Sprintf("%vql://%v@%v:%v/%v?sslmode=%v", baseType, user, host, port, baseName, ssl)

	err = graph.InitQuadStore("sql", tableAddr, graph.Options{"flavor": "cockroach"})
	if err != nil {
		log.Println(err)
	}

	store, err := cayley.NewGraph("sql", tableAddr, graph.Options{"flavor": "cockroach"})
	if err != nil {
		log.Fatal(err)
		return db, nil, err
	}

	// defer store.Close()

	engine.Store = store

	return db, store, nil
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
	transaction.AddQuad(cayley.Quad(companyName, "was added", companyAddTime, "Time of adding a company"))
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

func (engine *Engine) GetProductOfCompany(product *Product) {}

func (engine *Engine) SaveProductOfCompany(product *Product) {
	// Проверить наличие продукта по имени
}

func (engine *Engine) GetPricesOfProductsByName(productName string) ([]*PriceOfProduct, error) {
	return []*PriceOfProduct{&PriceOfProduct{Name: productName}}, nil
}

// SavePriceForProductOfCompany method for save subject, predicate and object in graph database
func (engine *Engine) SavePriceForProductOfCompany(item *Item) (*PriceOfProduct, error) {
	price := PriceOfProduct{Name: item.Name}
	return &price, nil
}
