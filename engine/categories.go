package engine

import (
	"context"
	"errors"
	"fmt"
	"log"

	dataBaseClient "github.com/dgraph-io/dgraph/client"
)

var (
	// ErrCategoryDoesNotExist means that the category is in the database already
	ErrCategoryDoesNotExist = errors.New("category does not exist")

	// ErrCategoryAlreadyExist means that the category is in the database already
	ErrCategoryAlreadyExist = errors.New("category already exist")

	// ErrCategoryCantBeCreated means that the category can't be added to database
	ErrCategoryCantBeCreated = errors.New("category can't be created")

	// ErrCategoryCantBeUpdated means that the category can't be updated in database
	ErrCategoryCantBeUpdated = errors.New("category can't be updated")

	// ErrCategoryWithoutID means that the category can't be updated in database
	ErrCategoryWithoutID = errors.New("category without ID")

	// ErrCategoryCantBeDeleted means that the category is in the database already
	ErrCategoryCantBeDeleted = errors.New("category can't be deleted")
)

// CreateCategory make mategory and save it to storage
func (engine *Engine) CreateCategory(categoryName string) (Category, error) {
	category := Category{Name: categoryName}

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return category, err
	}

	defer client.Close()

	existsCategories, err := engine.ReadCategoriesByName(categoryName)
	if len(existsCategories) > 0 {
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	request := dataBaseClient.Req{}
	err = request.SetObject(&category)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeCreated
	}

	response, err := client.Run(context.Background(), &request)

	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeCreated
	}

	// fmt.Printf("Raw Response: %+v\n", proto.MarshalTextString(response))

	for _, id := range response.AssignedUids {
		category.ID = id
	}

	return category, nil
}

// ReadCategoryByID is a method for get all nodes by categories names
func (engine *Engine) ReadCategoryByID(categoryID uint64) (Category, error) {
	type categoryInDatabase struct {
		Category Category `json:"category"`
	}

	var categoryFromDatabase categoryInDatabase

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return categoryFromDatabase.Category, err
	}

	defer client.Close()

	query := fmt.Sprintf(`{
			category(func: uid(%d)) {
				_uid_
				name
			}
		}`, categoryID)

	request := &dataBaseClient.Req{}
	request.SetQuery(query)

	response, err := client.Run(context.Background(), request)
	if err != nil {
		log.Println(err)
	}

	if len(response.N[0].Children) == 0 {
		log.Println(ErrCategoryDoesNotExist)
		return categoryFromDatabase.Category, ErrCategoryDoesNotExist
	}

	err = dataBaseClient.Unmarshal(response.N, &categoryFromDatabase)
	if err != nil {
		log.Println(err)
	}

	if categoryFromDatabase.Category.Name == "" {
		return categoryFromDatabase.Category, ErrCategoryDoesNotExist
	}

	return categoryFromDatabase.Category, nil
}

// ReadCategoriesByName is a method for get all nodes by categories names
func (engine *Engine) ReadCategoriesByName(categoryName string) ([]Category, error) {

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer client.Close()

	query := fmt.Sprintf(`{
				categories(func: allofterms(name, "%v")) {
					_uid_
					name
				}
			}`, categoryName)

	request := &dataBaseClient.Req{}
	request.SetQuery(query)

	response, err := client.Run(context.Background(), request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// fmt.Printf("Raw Response: %+v\n", proto.MarshalTextString(response))

	if len(response.N[0].Children) == 0 {
		log.Println(ErrCategoryDoesNotExist)
		return nil, ErrCategoryDoesNotExist
	}

	type categoriesInDatabase struct {
		AllCategoriesFoundedByName []Category `json:"categories"`
	}

	var categories categoriesInDatabase

	err = dataBaseClient.Unmarshal(response.N, &categories)
	if err != nil {
		log.Println(err)
		return nil, ErrCategoryDoesNotExist
	}

	if len(categories.AllCategoriesFoundedByName) == 0 {
		return nil, ErrCategoryDoesNotExist
	}

	return categories.AllCategoriesFoundedByName, nil
}

// UpdateCategory method for remove categories from database
func (engine *Engine) UpdateCategory(category *Category) (*Category, error) {

	if category.ID == 0 {
		return category, ErrCategoryWithoutID
	}

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeUpdated
	}

	defer client.Close()

	request := dataBaseClient.Req{}

	err = request.SetObject(category)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeUpdated
	}

	_, err = client.Run(context.Background(), &request)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryCantBeUpdated
	}

	return category, nil
}

// DeleteCategory method for remove categories from database
func (engine *Engine) DeleteCategory(category Category) (uint64, error) {
	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return 0, ErrCategoryCantBeDeleted
	}

	defer client.Close()

	request := dataBaseClient.Req{}

	err = request.DeleteObject(&category)
	if err != nil {
		log.Println(err)
		return 0, ErrCategoryCantBeDeleted
	}

	_, err = client.Run(context.Background(), &request)
	if err != nil {
		log.Println(err)
		return 0, ErrCategoryCantBeDeleted
	}

	return category.ID, nil
}
