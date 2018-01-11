package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

var (
	// ErrCompanyCanNotBeWithoutID means that company can't be found in storage for make some operation
	ErrCompanyCanNotBeWithoutID = errors.New("company can not be without id")

	// ErrCompanyCanNotBeDeactivate means that the company can't be deactivate in database
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

	// ErrCompanyCanNotBeDeleted means that the company can't be removed from database
	ErrCompanyCanNotBeDeleted = errors.New("company can't be deleted")

	// ErrCategoryCanNotBeAddedToCompany means that the company can't be removed from database
	ErrCategoryCanNotBeAddedToCompany = errors.New("category can not be added to company")
)

// Company is a structure of Categories in database
// TODO: Добавить поддержку языка для name предиката
type Company struct {
	ID         string     `json:"uid, omitempty"`
	IRI        string     `json:"companyIri, omitempty"`
	Name       string     `json:"companyName, omitempty"`
	Categories []Category `json:"has_category, omitempty"`
	IsActive   bool       `json:"companyIsActive, omitempty"`
}

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
		companyName: string @index(exact, term) .
		companyIsActive: bool @index(bool) .
		has_category: uid @count .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = companies.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

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

	return company, nil
}

// ReadCompaniesByName is a method for get all nodes by categories name
func (companies *Companies) ReadCompaniesByName(companyName string) ([]Company, error) {
	query := fmt.Sprintf(`{
				companies(func: eq(companyName, "%v")) @filter(eq(companyIsActive, true)) {
					uid
					companyName
					companyIri
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName
					}
					companyIsActive
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
	company := Company{ID: companyID}

	if companyID == "" {
		return company, ErrCompanyCanNotBeWithoutID
	}

	query := fmt.Sprintf(`{
				companies(func: uid("%s")) @filter(has(companyName)) {
					uid
					companyName
					companyIri
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName
						belongs_to_company {
							uid
							companyName
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName
								categoryIsActive
							}
							companyIsActive
						}
						categoryIsActive
					}
					companyIsActive
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

	return company.ID, nil
}

// DeleteCompany method for remove company from database
func (companies *Companies) DeleteCompany(company Company) (string, error) {

	if company.ID == "" {
		return "", ErrCompanyCanNotBeWithoutID
	}

	deleteCompanyData, _ := json.Marshal(map[string]string{"uid": company.ID})

	mutation := dataBaseAPI.Mutation{
		DeleteJson:          deleteCompanyData,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := companies.storage.Client.NewTxn()

	var err error
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return company.ID, ErrCompanyCanNotBeDeleted
	}

	return company.ID, nil
}

// AddCategoryToCompany method for set quad of predicate about company and category
func (companies *Companies) AddCategoryToCompany(companyID, categoryID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCompanyCanNotBeAddedToCategory
	}

	forCompanyPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, companyID, "has_category", categoryID)

	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCompanyPredicate),
		CommitNow: true}

	transaction = companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeAddedToCompany
	}

	return nil
}
