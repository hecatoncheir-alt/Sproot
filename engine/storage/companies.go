package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"bytes"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"text/template"
)

// Company is a structure of Categories in database
/* Для того что бы продукт принадлежащий компании отображался в категории принадлежащей
компании нужно иметь корректные данные belongs_to_company и belongs_to_category на гранях продукта */
type Company struct {
	ID         string     `json:"uid,omitempty"`
	IRI        string     `json:"companyIri,omitempty"`
	Name       string     `json:"companyName,omitempty"`
	Categories []Category `json:"has_category,omitempty"`
	IsActive   bool       `json:"companyIsActive"`
}

// Companies is resource of storage for CRUD operations
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
		companyName: string @lang @index(term) .
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

var (
	// ErrCompaniesByNameNotFound means than the companies does not exist in database
	ErrCompaniesByNameNotFound = errors.New("companies by name not found")

	// ErrCompanyCanNotBeCreated means that the company can't be added to database
	ErrCompanyCanNotBeCreated = errors.New("company can't be created")

	// ErrCompanyAlreadyExist means that the company is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")
)

// CreateCompany make category and save it to storage
func (companies *Companies) CreateCompany(company Company, language string) (Company, error) {
	existsCompanies, err := companies.ReadCompaniesByName(company.Name, language)
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
		SetJson:   encodedCompany,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeCreated
	}

	company.ID = assigned.Uids["blank-0"]
	if company.ID == "" {
		return company, ErrCompanyCanNotBeCreated
	}

	err = companies.AddLanguageOfCompanyName(company.ID, company.Name, language)
	if err != nil {
		return company, err
	}

	return company, nil
}

// AddLanguageOfCompanyName is a method for add predicate "companyName" for companyName value with new language
func (companies *Companies) AddLanguageOfCompanyName(companyID, name, language string) error {
	forCompanyNamePredicate := fmt.Sprintf(`<%s> <companyName> %s .`, companyID, "\""+name+"\""+"@"+language)

	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(forCompanyNamePredicate),
		CommitNow: true}

	transaction := companies.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}

var (
	// ErrCompaniesByNameCanNotBeFound means that the companies can't be found in database
	ErrCompaniesByNameCanNotBeFound = errors.New("companies by name can not be found")
)

// ReadAllCompanies is a method for get all nodes
func (companies *Companies) ReadAllCompanies(language string) ([]Company, error) {
	query := fmt.Sprintf(`{
				companies(func: eq(companyIsActive, true)) @filter(has(companyName)) {
					uid
					companyName: companyName@%v
					companyIri
					companyIsActive
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@%v
						categoryIsActive
					}
				}
			}`, language, language)

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

// ReadCompaniesByName is a method for get all nodes by categories name
func (companies *Companies) ReadCompaniesByName(companyName, language string) ([]Company, error) {
	variables := struct {
		CompanyName string
		Language    string
	}{
		CompanyName: companyName,
		Language:    language}

	queryTemplate, err := template.New("ReadCompaniesByName").Parse(`{
				companies(func: eq(companyName@{{.Language}}, "{{.CompanyName}}")) @filter(eq(companyIsActive, true)) {
					uid
					companyName: companyName@{{.Language}}
					companyIri
					companyIsActive
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIsActive
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
						}
						has_product @filter(eq(productIsActive, true)) { #TODO: belongs_to_company mast be an companyID
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							belongs_to_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
							belongs_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
								belongs_to_product @filter(eq(productIsActive, true)) {
									uid
									productName: productName@{{.Language}}
									productIri
									previewImageLink
									productIsActive
									has_price @filter(eq(priceIsActive, true)) {
										uid
										priceValue
										priceDateTime
										priceCity
										priceIsActive
									}
								}
								belongs_to_city @filter(eq(cityIsActive, true)) {
									uid
									cityName: cityName@{{.Language}}
									cityIsActive
								}
							}
						}
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	transaction := companies.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
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

var (
	// ErrCompanyCanNotBeWithoutID means that company can't be found in storage for make some operation
	ErrCompanyCanNotBeWithoutID = errors.New("company can not be without id")

	// ErrCompanyByIDCanNotBeFound means that the company can't be found in database
	ErrCompanyByIDCanNotBeFound = errors.New("company by id can not be found")

	// ErrCompanyDoesNotExist means than the company does not exist in database
	ErrCompanyDoesNotExist = errors.New("company by id not found")
)

// ReadCompanyByID is a method for get all nodes of categories by ID
func (companies *Companies) ReadCompanyByID(companyID, language string) (Company, error) {
	company := Company{ID: companyID}

	if companyID == "" {
		return company, ErrCompanyCanNotBeWithoutID
	}

	variables := struct {
		CompanyID string
		Language  string
	}{
		CompanyID: companyID,
		Language:  language}

	queryTemplate, err := template.New("ReadCompanyByID").Parse(`{
				companies(func: uid("{{.CompanyID}}")) @filter(has(companyName)) {
					uid
					companyName: companyName@{{.Language}}
					companyIri
					companyIsActive
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIsActive
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
						}
						has_product @filter(uid_in(belongs_to_company, {{.CompanyID}}) AND eq(productIsActive, true)) {
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							belongs_to_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
							belongs_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
								belongs_to_product @filter(eq(productIsActive, true)) {
									uid
									productName: productName@{{.Language}}
									productIri
									previewImageLink
									productIsActive
									has_price @filter(eq(priceIsActive, true)) {
										uid
										priceValue
										priceDateTime
										priceCity
										priceIsActive
									}
								}
								belongs_to_city @filter(eq(cityIsActive, true)) {
									uid
									cityName: cityName@{{.Language}}
									cityIsActive
								}
							}
						}
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return company, err
	}

	transaction := companies.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
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

var (
	// ErrCompanyCanNotBeUpdated means that company can't be updated
	ErrCompanyCanNotBeUpdated = errors.New("company can not be updated")
)

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
		SetJson:   encodedCompany,
		CommitNow: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeUpdated
	}

	updatedCompany, err := companies.ReadCompanyByID(company.ID, ".")
	if err != nil {
		log.Println(err)
		return company, ErrCompanyCanNotBeUpdated
	}

	return updatedCompany, nil
}

var (
	// ErrCompanyCanNotBeDeactivate means that the company can't be deactivate in database
	ErrCompanyCanNotBeDeactivate = errors.New("company can't be deactivated")
)

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

var (
	// ErrCategoryCanNotBeAddedToCompany means that the category can't be added to company
	ErrCategoryCanNotBeAddedToCompany = errors.New("category can not be added to company")

	// ErrCompanyCanNotBeDeleted means that the company can't be removed from database
	ErrCompanyCanNotBeDeleted = errors.New("company can't be deleted")
)

// DeleteCompany method for remove company from database
func (companies *Companies) DeleteCompany(company Company) (string, error) {

	if company.ID == "" {
		return "", ErrCompanyCanNotBeWithoutID
	}

	deleteCompanyData, _ := json.Marshal(map[string]string{"uid": company.ID})

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteCompanyData,
		CommitNow:  true}

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

var (
	// ErrCategoryCanNotBeRemovedFromCompany means that the company can't be removed from database
	ErrCategoryCanNotBeRemovedFromCompany = errors.New("category can not be removed from company")
)

// RemoveCategoryFromCompany method for delete quad of predicate about company and category
func (companies *Companies) RemoveCategoryFromCompany(companyID, categoryID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		DelNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCompanyCanNotBeRemovedFromCategory
	}

	forCompanyPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, companyID, "has_category", categoryID)

	mutation = dataBaseAPI.Mutation{
		DelNquads: []byte(forCompanyPredicate),
		CommitNow: true}

	transaction = companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeRemovedFromCompany
	}

	return nil
}

var (
	// ErrProductCanNotBeAddedToCompany means that the product can't be added to company
	ErrProductCanNotBeAddedToCompany = errors.New("product can not be added to company")
)

// AddProductToCompany method for set quad of predicate about company and product
func (companies *Companies) AddProductToCompany(companyID, productID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction := companies.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToCompany
	}

	return nil
}

type allExportedCompanies struct {
	Language  string    `json:"language"`
	Companies []Company `json:"companies"`
}

// ImportJSON is a method for add companies, categories of companies, products of categories,
// prices of products and cities of prices to database.
func (companies *Companies) ImportJSON(exportedCompanies []byte) error {

	var allCompaniesInJSON allExportedCompanies

	err := json.Unmarshal(exportedCompanies, &allCompaniesInJSON)
	if err != nil {
		return err
	}

	language := allCompaniesInJSON.Language

	for _, exportedCompany := range allCompaniesInJSON.Companies {

		encodedCompany, err := json.Marshal(exportedCompany)
		if err != nil {
			log.Println(err)
			return err
		}

		mutation := &dataBaseAPI.Mutation{
			SetJson:   encodedCompany,
			CommitNow: true}

		transaction := companies.storage.Client.NewTxn()
		_, err = transaction.Mutate(context.Background(), mutation)
		if err != nil {
			log.Println(err)
			return err
		}

		err = companies.AddLanguageOfCompanyName(exportedCompany.ID, exportedCompany.Name, language)
		if err != nil {
			return err
		}

		for _, exportedCategory := range exportedCompany.Categories {
			encodedCategory, err := json.Marshal(exportedCategory)
			if err != nil {
				log.Println(err)
				return err
			}

			mutation := &dataBaseAPI.Mutation{
				SetJson:   encodedCategory,
				CommitNow: true}

			transaction := companies.storage.Client.NewTxn()
			_, err = transaction.Mutate(context.Background(), mutation)
			if err != nil {
				log.Println(err)
				return err
			}

			err = companies.storage.Categories.AddLanguageOfCategoryName(exportedCategory.ID, exportedCategory.Name, language)
			if err != nil {
				log.Println(err)
				return err
			}

			err = companies.storage.Categories.AddCompanyToCategory(exportedCategory.ID, exportedCompany.ID)
			if err != nil {
				return err
			}

			for _, exportedProduct := range exportedCategory.Products {
				encodedProduct, err := json.Marshal(exportedProduct)
				if err != nil {
					log.Println(err)
					return err
				}

				mutation := &dataBaseAPI.Mutation{
					SetJson:   encodedProduct,
					CommitNow: true}

				transaction := companies.storage.Client.NewTxn()
				_, err = transaction.Mutate(context.Background(), mutation)
				if err != nil {
					log.Println(err)
					return err
				}

				err = companies.storage.Products.AddLanguageOfProductName(exportedProduct.ID, exportedProduct.Name, language)
				if err != nil {
					return err
				}

				err = companies.storage.Products.AddCategoryToProduct(exportedProduct.ID, exportedCategory.ID)
				if err != nil {
					return err
				}

				err = companies.storage.Products.AddCompanyToProduct(exportedProduct.ID, exportedCompany.ID)
				if err != nil {
					return err
				}

				for _, exportedPrice := range exportedProduct.Prices {
					encodedPrice, err := json.Marshal(exportedPrice)
					if err != nil {
						log.Println(err)
						return err
					}

					mutation := &dataBaseAPI.Mutation{
						SetJson:   encodedPrice,
						CommitNow: true}

					transaction := companies.storage.Client.NewTxn()
					_, err = transaction.Mutate(context.Background(), mutation)
					if err != nil {
						log.Println(err)
						return err
					}

					err = companies.storage.Products.AddPriceToProduct(exportedProduct.ID, exportedPrice.ID)
					if err != nil {
						log.Println(err)
						return err
					}

					for _, exportedCity := range exportedPrice.Cities {
						encodedCity, err := json.Marshal(exportedCity)
						if err != nil {
							log.Println(err)
							return err
						}

						mutation := &dataBaseAPI.Mutation{
							SetJson:   encodedCity,
							CommitNow: true}

						transaction := companies.storage.Client.NewTxn()
						_, err = transaction.Mutate(context.Background(), mutation)
						if err != nil {
							log.Println(err)
							return err
						}

						err = companies.storage.Cities.AddLanguageOfCityName(exportedCity.ID, exportedCity.Name, language)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// ExportJSON is a method for export companies, categories of companies, products of categories,
// prices of products and cities of prices to database.
func (companies *Companies) ExportJSON(language string) ([]byte, error) {
	query := fmt.Sprintf(`{
				companies(func: has(companyName)) {
					uid
				}
			}`)

	transaction := companies.storage.Client.NewTxn()
	responseWithCompaniesIDs, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	type allCompaniesWithIDOnly struct {
		CompaniesWithIDOnly []Company `json:"companies"`
	}

	var allCompaniesIDs allCompaniesWithIDOnly

	err = json.Unmarshal(responseWithCompaniesIDs.GetJson(), &allCompaniesIDs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	type allExportedCompanies struct {
		Language  string    `json:"language"`
		Companies []Company `json:"companies"`
	}

	foundedCompanies := allExportedCompanies{Language: language}

	for _, companyID := range allCompaniesIDs.CompaniesWithIDOnly {
		company, err := companies.ReadCompanyByID(companyID.ID, language)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		foundedCompanies.Companies = append(foundedCompanies.Companies, company)
	}

	jsonForExport, err := json.Marshal(foundedCompanies)
	if err != nil {
		return nil, err
	}

	return jsonForExport, nil
}
