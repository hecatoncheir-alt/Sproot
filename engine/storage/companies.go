package storage

import (
	"context"
	"log"
	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
	"encoding/json"
	"errors"
	"fmt"
)

type Companies struct {
	storage *Storage
}

func NewCompaniesResourceForStorage(storage *Storage) *Companies {
	return &Companies{storage: storage}
}

func (companies *Companies) SetUp() (err error) {
	schema := "name: string @index(exact, term) ."
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = companies.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

var (
	// ErrCompanyCanNotBeWithoutID means that company can't be found in storage for make some operation
	ErrCompanyCanNotBeWithoutID = errors.New("company can not be without id")

	// ErrCompanyCanNotBeDeleted means that the company can't be deleted from database
	ErrCompanyCanNotBeDeleted = errors.New("company can't be deleted")

	// ErrCompaniesByNameCanNotBeFound means that the companies can't be found in database
	ErrCompaniesByNameCanNotBeFound = errors.New("companies by name can not be found")

	// ErrCompaniesByNameNotFound means than the companies does not exist in database
	ErrCompaniesByNameNotFound = errors.New("companies by name not found")

	// ErrCompanyCanNotBeCreated means that the category can't be added to database
	ErrCompanyCanNotBeCreated = errors.New("company can't be created")

	// ErrCompanyAlreadyExist means that the category is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")
)

// Company is a structure of Categories in database
type Company struct {
	ID         string      `json:"uid,omitempty"`
	IRI        string      `json:"iri, omitempty"`
	Name       string      `json:"name,omitempty"`
	Categories []*Category `json:"category, omitempty"`
}

// CreateCategory make category and save it to storage
func (companies *Companies) CreateCompany(company *Company) (Company, error) {
	existsCompanies, err := categories.ReadCategoriesByName(category.Name)
	if err != nil && err != ErrCompaniesByNameNotFound {
		log.Println(err)
		return *category, ErrCompanyCanNotBeCreated
	}

	if existsCompanies != nil {
		return existsCompanies[0], ErrCompanyAlreadyExist
	}
	//
	//transaction := categories.storage.Client.NewTxn()
	//
	//encodedCategory, err := json.Marshal(category)
	//if err != nil {
	//	log.Println(err)
	//	return *category, ErrCategoryCanNotBeCreated
	//}
	//
	//mutation := &dataBaseAPI.Mutation{
	//	SetJson:             encodedCategory,
	//	CommitNow:           true,
	//	IgnoreIndexConflict: true}
	//
	//assigned, err := transaction.Mutate(context.Background(), mutation)
	//if err != nil {
	//	log.Println(err)
	//	return *category, ErrCategoryCanNotBeCreated
	//}
	//
	//category.ID = assigned.Uids["blank-0"]
	//if category.ID == "" {
	//	return *category, ErrCategoryCanNotBeCreated
	//}
	//
	return *company, nil
}

// ReadCompaniesByName is a method for get all nodes by categories name
func (companies *Companies) ReadCompaniesByName(companyName string) ([]Company, error) {
	query := fmt.Sprintf(`{
				companies(func: eq(name, "%v")) {
					uid
					name
				}
			}`, companyName)

	transaction := companies.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	type companiesInStorage struct {
		AllCompaniesFoundedByName []Company `json:"companies"`
	}

	var foundedCompanies companiesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCompanies)
	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	if len(foundedCompanies.AllCompaniesFoundedByName) == 0 {
		return nil, ErrCompaniesByNameNotFound
	}

	return foundedCompanies.AllCompaniesFoundedByName, nil
}

// DeleteCompany method for remove categories from database
func (companies *Companies) DeleteCompany(company Company) (string, error) {
	if company.ID == "" {
		return "", ErrCompanyCanNotBeWithoutID
	}

	encodedCompany, err := json.Marshal(Company{ID: company.ID})
	if err != nil {
		log.Println(err)
		return "", ErrCompanyCanNotBeDeleted
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson:          encodedCompany,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return "", ErrCompanyCanNotBeDeleted
	}

	return company.ID, nil
}
