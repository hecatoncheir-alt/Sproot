package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

// Companies is resource os storage for CRUD operations
type Companies struct {
	storage *Storage
}

// NewCompaniesResourceForStorage is a constructor of Categories resource
func NewCompaniesResourceForStorage(storage *Storage) *Companies {
	return &Companies{storage: storage}
}

// SetUp is a method of Companies resource for prepare database client and schema.
func (companies *Companies) SetUp() (err error) {
	schema := `
		name: string @index(exact, term) .
		isActive: bool @index(bool) .
	`
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

	// ErrCompanyCanNotBeDeactivate means that the company can't be deactivate from database
	ErrCompanyCanNotBeDeactivate = errors.New("company can't be deactivated")

	// ErrCompaniesByNameCanNotBeFound means that the companies can't be found in database
	ErrCompaniesByNameCanNotBeFound = errors.New("companies by name can not be found")

	// ErrCompaniesByNameNotFound means than the companies does not exist in database
	ErrCompaniesByNameNotFound = errors.New("companies by name not found")

	// ErrCompanyCanNotBeCreated means that the company can't be added to database
	ErrCompanyCanNotBeCreated = errors.New("company can't be created")

	// ErrCompanyAlreadyExist means that the company is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")

	// ErrCompanyByIDCanNotBeFound means that the company can't be found in database
	ErrCompanyByIDCanNotBeFound = errors.New("company by id can not be found")

	// ErrCompanyCanNotBeUpdated means that company can't be updated
	ErrCompanyCanNotBeUpdated = errors.New("company can not be updated")

	// ErrCompanyDoesNotExist means than the company does not exist in database
	ErrCompanyDoesNotExist = errors.New("company by id not found")
)

// CreateCompany make category and save it to storage
func (companies *Companies) CreateCompany(company Company) (Company, error) {
	existsCompanies, err := companies.ReadCompaniesByName(company.Name)
	if err != nil && err != ErrCompaniesByNameNotFound {
		log.Println(err)
		return company, ErrCompanyCanNotBeCreated
	}

	if existsCompanies != nil {
		return existsCompanies[0], ErrCompanyAlreadyExist
	}

	transaction := companies.storage.Client.NewTxn()

	company.IsActive = true
	encodedCompany, err := json.Marshal(company)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:             encodedCompany,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeCreated
	}

	company.ID = assigned.Uids["blank-0"]
	if company.ID == "" {
		return company, ErrCompanyCanNotBeCreated
	}

	company.storage = companies.storage
	return company, nil
}

// ReadCompaniesByName is a method for get all nodes by categories name
func (companies *Companies) ReadCompaniesByName(companyName string) ([]Company, error) {
	query := fmt.Sprintf(`{
				companies(func: eq(name, "%v")) @filter(eq(isActive, true)) {
					uid
					name
					iri
					categories @filter(eq(isActive, true)) {
						uid
						name
					}
					isActive
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

// ReadCompanyByID is a method for get all nodes of categories by ID
func (companies *Companies) ReadCompanyByID(companyID string) (Company, error) {
	if companyID == "" {
		return Company{}, ErrCompanyCanNotBeWithoutID
	}

	company := Company{ID: companyID}

	query := fmt.Sprintf(`{
				companies(func: uid("%s")) @filter(eq(isActive, true)) {
					uid
					name
					iri
					categories @filter(eq(isActive, true)) {
						uid
						name
						companies {
							uid
							name
						}
					}
					isActive
				}
			}`, companyID)

	transaction := companies.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	type companiesInStore struct {
		Companies []Company `json:"companies"`
	}

	var foundedCompanies companiesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCompanies)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	if len(foundedCompanies.Companies) == 0 {
		return company, ErrCompanyDoesNotExist
	}

	foundedCompanies.Companies[0].storage = companies.storage
	return foundedCompanies.Companies[0], nil
}

// UpdateCompany method for change company in storage
func (companies *Companies) UpdateCompany(company Company) (Company, error) {
	if company.ID == "" {
		return company, ErrCompanyCanNotBeWithoutID
	}

	encodedCompany, err := json.Marshal(company)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeUpdated
	}

	mutation := dataBaseAPI.Mutation{
		SetJson:             encodedCompany,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeUpdated
	}

	updatedCompany, err := companies.ReadCompanyByID(company.ID)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeUpdated
	}

	updatedCompany.storage = companies.storage
	return updatedCompany, nil
}

// DeactivateCompany method for remove categories from database
func (companies *Companies) DeactivateCompany(company Company) (string, error) {
	if company.ID == "" {
		return "", ErrCompanyCanNotBeWithoutID
	}

	company.IsActive = false

	encodedCompany, err := json.Marshal(company)
	if err != nil {
		log.Println(err)
		return company.ID, ErrCompanyCanNotBeDeactivate
	}

	mutation := dataBaseAPI.Mutation{
		SetJson:             encodedCompany,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := companies.storage.Client.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return company.ID, ErrCompanyCanNotBeDeactivate
	}

	updatedCompany, err := companies.ReadCompanyByID(company.ID)
	if err != nil && err != ErrCompanyDoesNotExist {
		log.Println(err)
		return "", ErrCompanyCanNotBeDeactivate
	}

	return updatedCompany.ID, nil
}
