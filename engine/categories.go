package engine

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"context"
	dataBaseClient "github.com/dgraph-io/dgraph/client"
)

var (
	// ErrCategoriesAlreadyExists means that the categories is in the database already
	ErrCategoryDoesNotExist = errors.New("category does not exist")

	// ErrCategoriesAlreadyExists means that the categories is in the database already
	ErrCategoriesAlreadyExists = errors.New("categories already exists")

	// ErrCategoryCanBeCreated means that the categories can't be added to database
	ErrCategoryCantBeCreated = errors.New("category can't be created")

	// ErrCategoryCantBeDeleted means that the categories is in the database already
	ErrCategoryCantBeDeleted = errors.New("category can't be deleted")
)

// DeleteCategories method for remove categories from database
func (engine *Engine) DeleteCategories(categories []Category) ([]uint64, error) {
	deletedIDs := []uint64{}

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return nil, ErrCategoryCantBeDeleted
	}

	defer client.Close()

	for _, category := range categories {
		request := dataBaseClient.Req{}
		err = request.DeleteObject(category)
		if err != nil {
			log.Println(err)
			return nil, ErrCategoryCantBeDeleted
		}

		_, err := client.Run(context.Background(), &request)
		if err != nil {
			log.Println(err)
			return nil, ErrCategoryCantBeDeleted
		}

		deletedIDs = append(deletedIDs, category.ID)
	}

	return deletedIDs, nil
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

	return categoryFromDatabase.Category, nil
}

// ReadCategoriesByName is a method for get all nodes by categories names
func (engine *Engine) ReadCategoriesByName(categoriesNames []string) (map[string][]Category, error) {
	categoriesByName := map[string][]Category{}

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return categoriesByName, err
	}

	defer client.Close()

	for _, categoryName := range categoriesNames {
		query := fmt.Sprintf(`{
			category(func: allofterms(name, "%v")) {
				_uid_
				name
			}
		}`, categoryName)

		request := &dataBaseClient.Req{}
		request.SetQuery(query)

		response, err := client.Run(context.Background(), request)
		if err != nil {
			log.Println(err)
		}

		if len(response.N[0].Children) == 0 {
			log.Println(ErrCategoryDoesNotExist)
			continue
		}

		type categoriesInDatabase struct {
			InDatabase []Category `dgraph:"category"`
		}
		var categories categoriesInDatabase

		err = dataBaseClient.Unmarshal(response.N[0].Children, &categories)
		if err != nil {
			log.Println(err)
		}

		categoriesByName[categoryName] = append(categoriesByName[categoryName], categories.InDatabase...)
	}

	return categoriesByName, nil
}

// CreateCategories is a method for add node for each category in database
func (engine *Engine) CreateCategories(categoriesNames []string) ([]Category, error) {
	if !sort.StringsAreSorted(categoriesNames) {
		sort.Strings(categoriesNames)
	}

	var createdCategories []Category

	existCategoriesByName, err := engine.ReadCategoriesByName(categoriesNames)
	if err != nil {
		log.Fatal(err)
		return createdCategories, err
	}

	if len(existCategoriesByName) > 0 {
		for _, categoryName := range categoriesNames {
			for _, existCategory := range existCategoriesByName[categoryName] {
				createdCategories = append(createdCategories, existCategory)
			}

			index := sort.SearchStrings(categoriesNames, categoryName)
			categoriesNames = append(categoriesNames[:index], categoriesNames[index+1:]...)
		}
	}

	if len(categoriesNames) == 0 {
		return createdCategories, ErrCategoriesAlreadyExists
	}

	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return createdCategories, err
	}

	defer client.Close()

	request := &dataBaseClient.Req{}

	request.SetSchema(`
				name: string @index(exact, term) .
	`)

	for _, categoryName := range categoriesNames {
		category := Category{Name: categoryName}

		err = request.SetObject(&category)
		if err != nil {
			log.Println(err)
			return createdCategories, ErrCategoryCantBeCreated
		}

		response, err := client.Run(context.Background(), request)

		if err != nil {
			log.Println(err)
			return createdCategories, ErrCategoryCantBeCreated
		}

		for _, id := range response.AssignedUids {
			category.ID = id
			createdCategories = append(createdCategories, category)
		}
	}

	return createdCategories, nil
}
