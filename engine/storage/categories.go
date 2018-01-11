package storage

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"fmt"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

var (
	// ErrCategoryCanNotBeCreated means that the category can't be added to database
	ErrCategoryCanNotBeCreated = errors.New("category can't be created")

	// ErrCategoryAlreadyExist means that the category is in the database already
	ErrCategoryAlreadyExist = errors.New("category already exist")

	// ErrCategoriesByNameCanNotBeFound means that the category can't be found in database
	ErrCategoriesByNameCanNotBeFound = errors.New("categories by name can not be found")

	// ErrCategoriesByNameNotFound means than the categories does not exist in database
	ErrCategoriesByNameNotFound = errors.New("categories by name not found")

	// ErrCategoryByIDCanNotBeFound means that the category can't be found in database
	ErrCategoryByIDCanNotBeFound = errors.New("category by id can not be found")

	// ErrCategoryDoesNotExist means than the categories does not exist in database
	ErrCategoryDoesNotExist = errors.New("categories by id not found")

	// ErrCategoryCanNotBeWithoutID means that category can't be found in storage for make some operation
	ErrCategoryCanNotBeWithoutID = errors.New("category can not be without id")

	// ErrCategoryCanNotBeUpdated means that category can't be updated
	ErrCategoryCanNotBeUpdated = errors.New("category can not be updated")

	// ErrCategoryCanNotBeDeactivate means that the category can't be deactivate in database
	ErrCategoryCanNotBeDeactivate = errors.New("category can't be deactivate")

	// ErrCategoryCanNotBeDeleted means that the category can't be removed from database
	ErrCategoryCanNotBeDeleted = errors.New("category can't be deleted")
)

// Category is a structure of Categories in database
type Category struct {
	ID        string    `json:"uid, omitempty"`
	Name      string    `json:"categoryName, omitempty"`
	IsActive  bool      `json:"categoryIsActive, omitempty"`
	Companies []Company `json:"belongs_to_company, omitempty"`
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
		categoryName: string @index(exact, term) .
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = categories.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// CreateCategory make category and save it to storage
func (categories *Categories) CreateCategory(category Category) (Category, error) {
	existsCategories, err := categories.ReadCategoriesByName(category.Name)
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
		SetJson:             encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCanNotBeCreated
	}

	category.ID = assigned.Uids["blank-0"]
	if category.ID == "" {
		return category, ErrCategoryCanNotBeCreated
	}

	return category, nil
}

// ReadCategoriesByName is a method for get all nodes by categories name
func (categories *Categories) ReadCategoriesByName(categoryName string) ([]Category, error) {
	query := fmt.Sprintf(`{
				categories(func: eq(categoryName, "%v")) @filter(eq(categoryIsActive, true)) {
					uid
					categoryName
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName
						companyCategories {
							uid
							categoryName
						}
					}
					categoryIsActive
				}
			}`, categoryName)

	transaction := categories.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
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

// ReadCategoryByID is a method for get all nodes of categories by ID
func (categories *Categories) ReadCategoryByID(categoryID string) (Category, error) {
	category := Category{ID: categoryID}

	query := fmt.Sprintf(`{
				categories(func: uid("%s")) @filter(has(categoryName)) {
					uid
					categoryName
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName
						companyCategories {
							uid
							categoryName
						}
					}
					categoryIsActive
				}
			}`, categoryID)

	transaction := categories.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
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

	updatedCategory, err := categories.ReadCategoryByID(category.ID)
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

/// DeleteCategory method for remove category from database
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

//TODO
func (categories *Categories) AddCompanyToCategory(categoryID, companyID string) error {
	//category, _ := categories.ReadCategoryByID(categoryID)
	//company, _ := categories.storage.Companies.ReadCompanyByID(companyID)
	//category.Companies = append(category.Companies, company)
	//time.Sleep(time.Second * 2)
	//categories.UpdateCategory(*category)
	return nil
}
