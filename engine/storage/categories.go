package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"text/template"

	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
)

var (

	// ErrCategoryCanNotBeUpdated means that category can't be updated
	ErrCategoryCanNotBeUpdated = errors.New("category can not be updated")

	// ErrCategoryCanNotBeDeactivate means that the category can't be deactivate in database
	ErrCategoryCanNotBeDeactivate = errors.New("category can't be deactivate")

	// ErrCompanyCanNotBeRemovedFromCategory means that the company can't be removed from database
	ErrCompanyCanNotBeRemovedFromCategory = errors.New("company can not be removed from category")

	// ErrProductCanNotBeAddedToCategory means that the product can't be added to category
	ErrProductCanNotBeAddedToCategory = errors.New("product can not be added to category")
)

// Category is a structure of Categories in database
type Category struct {
	ID        string    `json:"uid,omitempty"`
	Name      string    `json:"categoryName,omitempty"`
	IsActive  bool      `json:"categoryIsActive"`
	Companies []Company `json:"belongs_to_company,omitempty"`
	Products  []Product `json:"has_product,omitempty"`
}

// Categories is resource os storage for CRUD operations
type Categories struct {
	storage *Storage
}

// NewCategoriesResourceForStorage is a constructor of Categories resource
func NewCategoriesResourceForStorage(storage *Storage) *Categories {
	return &Categories{storage: storage}
}

// SetUp is a method of Categories resource for prepare database client and schema.
func (categories *Categories) SetUp() (err error) {
	schema := `
		categoryName: string @lang @index(term).
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
		has_product: uid .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = categories.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

var (
	// ErrCategoriesByNameNotFound means than the categories does not exist in database
	ErrCategoriesByNameNotFound = errors.New("categories by name not found")

	// ErrCategoryCanNotBeCreated means that the category can't be added to database
	ErrCategoryCanNotBeCreated = errors.New("category can't be created")

	// ErrCategoryAlreadyExist means that the category is in the database already
	ErrCategoryAlreadyExist = errors.New("category already exist")
)

// CreateCategory make category and save it to storage
func (categories *Categories) CreateCategory(category Category, language string) (Category, error) {
	existsCategories, err := categories.ReadCategoriesByName(category.Name, language)
	if err != nil && err != ErrCategoriesByNameNotFound {
		log.Println(err)
		return category, ErrCategoryCanNotBeCreated
	}
	if existsCategories != nil {
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	transaction := categories.storage.Client.NewTxn()

	category.IsActive = true
	encodedCategory, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCategory,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeCreated
	}

	category.ID = assigned.Uids["blank-0"]
	if category.ID == "" {
		return category, ErrCategoryCanNotBeCreated
	}

	err = categories.AddLanguageOfCategoryName(category.ID, category.Name, language)
	if err != nil {
		return category, err
	}

	createdCategory, err := categories.ReadCategoryByID(category.ID, language)
	if err != nil {
		return category, err
	}

	return createdCategory, nil
}

// AddLanguageOfCategoryName is a method for add predicate "categoryName" for companyName value with new language
func (categories *Categories) AddLanguageOfCategoryName(categoryID, name, language string) error {
	forCategoryNamePredicate := fmt.Sprintf(`<%s> <categoryName> %s .`, categoryID, "\""+name+"\""+"@"+language)

	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryNamePredicate),
		CommitNow: true}

	transaction := categories.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}

var (
	// ErrCategoriesByNameCanNotBeFound means that the category can't be found in database
	ErrCategoriesByNameCanNotBeFound = errors.New("categories by name can not be found")
)

// ReadCategoriesByName is a method for get all nodes by categories name
func (categories *Categories) ReadCategoriesByName(categoryName, language string) ([]Category, error) {

	variables := struct {
		CategoryName string
		Language     string
	}{
		CategoryName: categoryName,
		Language:     language}

	queryTemplate, err := template.New("ReadCategoriesByName").Parse(`{
				categories(func: eq(categoryName@{{.Language}}, "{{.CategoryName}}"))
				@filter(eq(categoryIsActive, true)) {
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
							belong_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_product @filter(eq(productIsActive, true)) {
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
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	transaction := categories.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	type categoriesInStorage struct {
		AllCategoriesFoundedByName []Category `json:"categories"`
	}

	var foundedCategories categoriesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCategories)
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	if len(foundedCategories.AllCategoriesFoundedByName) == 0 {
		return nil, ErrCategoriesByNameNotFound
	}

	return foundedCategories.AllCategoriesFoundedByName, nil
}

var (
	// ErrCategoryByIDCanNotBeFound means that the category can't be found in database
	ErrCategoryByIDCanNotBeFound = errors.New("category by id can not be found")

	// ErrCategoryDoesNotExist means than the categories does not exist in database
	ErrCategoryDoesNotExist = errors.New("categories by id not found")
)

// ReadCategoryByID is a method for get all nodes of categories by ID
func (categories *Categories) ReadCategoryByID(categoryID, language string) (Category, error) {

	variables := struct {
		CategoryID string
		Language   string
	}{
		CategoryID: categoryID,
		Language:   language}

	queryTemplate, err := template.New("ReadCategoryByID").Parse(`{
				categories(func: uid("{{.CategoryID}}")) @filter(has(categoryName)) {
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
							belong_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_product @filter(eq(productIsActive, true)) {
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
					}
				}
			}`)

	category := Category{ID: categoryID}

	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	transaction := categories.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	type categoriesInStore struct {
		Categories []Category `json:"categories"`
	}

	var foundedCategories categoriesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCategories)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	if len(foundedCategories.Categories) == 0 {
		return category, ErrCategoryDoesNotExist
	}

	return foundedCategories.Categories[0], nil
}

// UpdateCategory method for change category in storage
func (categories *Categories) UpdateCategory(category Category) (Category, error) {
	if category.ID == "" {
		return category, ErrCategoryCanNotBeWithoutID
	}

	transaction := categories.storage.Client.NewTxn()

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeUpdated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCategory,
		CommitNow: true}

	_, err = transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeUpdated
	}

	updatedCategory, err := categories.ReadCategoryByID(category.ID, ".")
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeUpdated
	}

	return updatedCategory, nil
}

// DeactivateCategory method for remove categories from database
func (categories *Categories) DeactivateCategory(category Category) (string, error) {
	if category.ID == "" {
		return "", ErrCategoryCanNotBeWithoutID
	}

	categoryForUpdate := Category{
		ID:        category.ID,
		Name:      category.Name,
		Companies: category.Companies,
		IsActive:  false}

	updatedCategory, err := categories.UpdateCategory(categoryForUpdate)
	if err != nil {
		log.Println(err)
		return "", ErrCategoryCanNotBeDeactivate
	}

	return updatedCategory.ID, nil
}

var (
	// ErrCategoryCanNotBeWithoutID means that category can't be found in storage for make some operation
	ErrCategoryCanNotBeWithoutID = errors.New("category can not be without id")

	// ErrCategoryCanNotBeDeleted means that the category can't be removed from database
	ErrCategoryCanNotBeDeleted = errors.New("category can't be deleted")
)

// DeleteCategory method for remove category from database
func (categories *Categories) DeleteCategory(category Category) (string, error) {

	if category.ID == "" {
		return "", ErrCategoryCanNotBeWithoutID
	}

	deleteCategoryData, _ := json.Marshal(map[string]string{"uid": category.ID})

	mutation := dataBaseAPI.Mutation{
		DeleteJson:          deleteCategoryData,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := categories.storage.Client.NewTxn()

	var err error
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return category.ID, ErrCategoryCanNotBeDeleted
	}

	return category.ID, nil
}

var (
	// ErrCompanyCanNotBeAddedToCategory means that the company can't be added to category
	ErrCompanyCanNotBeAddedToCategory = errors.New("company can not be added to category")
)

// AddCompanyToCategory method for set quad of predicate about category and company
func (categories *Categories) AddCompanyToCategory(categoryID, companyID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forCompanyPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, companyID, "has_category", categoryID)

	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCompanyPredicate),
		CommitNow: true}

	transaction := categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeAddedToCompany
	}

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction = categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCompanyCanNotBeAddedToCategory
	}

	return nil
}

// RemoveCompanyFromCategory method for delete quad of predicate about category and company
func (categories *Categories) RemoveCompanyFromCategory(categoryID, companyID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forCompanyPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, companyID, "has_category", categoryID)

	mutation = dataBaseAPI.Mutation{
		DelNquads: []byte(forCompanyPredicate),
		CommitNow: true}

	transaction := categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeRemovedFromCompany
	}

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "belongs_to_company", companyID)
	mutation = dataBaseAPI.Mutation{
		DelNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction = categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCompanyCanNotBeRemovedFromCategory
	}

	return nil
}

// AddProductToCategory method for set quad of predicate about category and product
func (categories *Categories) AddProductToCategory(categoryID, productID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "belongs_to_category", categoryID)

	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction := categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToCategory
	}

	forCategoryPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, categoryID, "has_product", productID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryPredicate),
		CommitNow: true}

	transaction = categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToCategory
	}

	return nil
}
