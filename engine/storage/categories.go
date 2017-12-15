package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dgraph-io/dgraph/y"
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
	// ErrCategoryCantBeCreated means that the category can't be added to database
	ErrCategoryCantBeCreated = errors.New("category can't be created")

	// ErrCategoryAlreadyExist means that the category is in the database already
	ErrCategoryAlreadyExist = errors.New("category already exist")

	// ErrCategoriesByNameCanNotBeFound means that the category can't be found in database
	ErrCategoriesByNameCanNotBeFound = errors.New("categories by name can not be found")

	// ErrCategoriesByNameNotFound means than the categories does not exist in database
	ErrCategoriesByNameNotFound = errors.New("categories by name not found")

	// ErrCategoryCantBeDeleted means that the category can't be deleted from database
	ErrCategoryCantBeDeleted = errors.New("category can't be deleted")
)

// Category is a structure of Category in database
type Category struct {
	ID   string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}

// CreateCategory make category and save it to storage
func (categories *Categories) CreateCategory(categoryName string) (Category, error) {
	category := Category{Name: categoryName}

	existsCategories, err := categories.ReadCategoriesByName(categoryName)
	if err != nil && err != ErrCategoriesByNameNotFound {
		log.Println(err)
		return category, ErrCategoryCantBeCreated
	}
	if existsCategories != nil {
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	transaction := categories.storage.Client.NewTxn()

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:             encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeCreated
	}

	category.ID = assigned.Uids["blank-0"]
	if category.ID == "" {
		return category, ErrCategoryCantBeCreated
	}

	return category, nil
}

// ReadCategoriesByName is a method for get all nodes by categories names
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

// DeleteCategory method for remove categories from database
func (categories *Categories) DeleteCategory(category Category) (string, error) {
	encodedCategory, err := json.Marshal(Category{ID: category.ID})
	if err != nil {
		log.Println(err)
		return "", ErrCategoryCantBeDeleted
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson:          encodedCategory,
		CommitNow:           true,
		IgnoreIndexConflict: true}

	transaction := categories.storage.Client.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		if err == y.ErrAborted {
			log.Println(err)
		} else {
			log.Println(err)
			return "", ErrCategoryCantBeDeleted
		}
	}

	categoryID := assigned.Uids["blank-0"]
	if categoryID == "" {
		return "", ErrCategoryCantBeDeleted
	}

	return categoryID, nil
}
