package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

type Categories struct {
	storage *Storage
}

func NewCategoriesResourceForStorage(storage *Storage) *Categories {
	return &Categories{storage: storage}
}

func (categories *Categories) SetUp() (err error) {
	schema := "name: string @index(exact, term) ."
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = categories.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

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

	// ErrCategoryCanNotBeDeleted means that the category can't be deleted from database
	ErrCategoryCanNotBeDeleted = errors.New("category can't be deleted")
)

// Categories is a structure of Categories in database
type Category struct {
	ID   string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateCategory make category and save it to storage
func (categories *Categories) CreateCategory(category *Category) (Category, error) {
	existsCategories, err := categories.ReadCategoriesByName(category.Name)
	if err != nil && err != ErrCategoriesByNameNotFound {
		log.Println(err)
		return *category, ErrCategoryCanNotBeCreated
	}
	if existsCategories != nil {
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	transaction := categories.storage.Client.NewTxn()

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		return *category, ErrCategoryCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:             encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return *category, ErrCategoryCanNotBeCreated
	}

	category.ID = assigned.Uids["blank-0"]
	if category.ID == "" {
		return *category, ErrCategoryCanNotBeCreated
	}

	return *category, nil
}

// ReadCategoriesByName is a method for get all nodes by categories name
func (categories *Categories) ReadCategoriesByName(categoryName string) ([]Category, error) {
	query := fmt.Sprintf(`{
				categories(func: eq(name, "%v")) {
					uid
					name
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
				category(func: uid("%s")) {
					uid
					name
				}
			}`, categoryID)

	transaction := categories.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	type categoryInStore struct {
		Categories []Category `json:"category"`
	}

	var foundedCategory categoryInStore

	err = json.Unmarshal(response.GetJson(), &foundedCategory)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	if foundedCategory.Categories[0].Name == "" {
		return category, ErrCategoryDoesNotExist
	}

	return foundedCategory.Categories[0], nil
}

// UpdateCategory make category and save it to storage
func (categories *Categories) UpdateCategory(category *Category) (Category, error) {
	if category.ID == "" {
		return *category, ErrCategoryCanNotBeWithoutID
	}

	existsCategories, err := categories.ReadCategoriesByName(category.Name)
	if err != nil && err != ErrCategoriesByNameNotFound {
		log.Println(err)
		return *category, ErrCategoryCanNotBeCreated
	}
	if existsCategories != nil {
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	transaction := categories.storage.Client.NewTxn()

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		return *category, ErrCategoryCanNotBeUpdated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:             encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	_, err = transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return *category, ErrCategoryCanNotBeUpdated
	}

	updatedCategory, err := categories.ReadCategoryByID(category.ID)
	if err != nil {
		log.Println(err)
		return *category, ErrCategoryCanNotBeUpdated
	}

	return updatedCategory, nil
}

// DeleteCategory method for remove categories from database
func (categories *Categories) DeleteCategory(category Category) (string, error) {
	encodedCategory, err := json.Marshal(Category{ID: category.ID})
	if err != nil {
		log.Println(err)
		return "", ErrCategoryCanNotBeDeleted
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson:          encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := categories.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return "", ErrCategoryCanNotBeDeleted
	}

	return category.ID, nil
}