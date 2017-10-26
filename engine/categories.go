package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var (
	// ErrCategoriesAlreadyExists means that the categories is in the database already
	ErrCategoriesAlreadyExists = errors.New("categories already exists")

	// ErrCategoriesCanBeCreated means that the categories can't be added to database
	ErrCategoriesCanBeCreated = errors.New("categories can't be created")

	// ErrCategoriesCantBeDeleted means that the categories is in the database already
	ErrCategoriesCantBeDeleted = errors.New("categories can't be deleted")
)

// DeleteCategories method for remove categories from database
func (engine *Engine) DeleteCategories(categories []Category) ([]Category, error) {
	request := bytes.NewBufferString(`
		mutation {
			delete {
	`)

	for _, category := range categories {
		request.WriteString("<" + string(category.ID) + "> * * .\n")
	}

	request.WriteString("}\n" + "}\n")

	req, err := http.NewRequest("POST", engine.GraphAddress+"/query", request)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	var details GraphResponse
	json.Unmarshal(responseData, &details)

	if details.Data.Code != "Success" {
		return categories, ErrCategoriesCantBeDeleted
	}

	return categories, nil
}

// ReadCategoriesByName is a method for get all nodes by categories names
func (engine *Engine) ReadCategoriesByName(categoriesNames []string) (map[string][]Category, error) {
	categoriesByName := map[string][]Category{}

	for _, categoryName := range categoriesNames {

		request := fmt.Sprintf(`{
			categories(func:allofterms(name, "%v")) {
				name
				_uid_
			}
		}`, categoryName)

		req, err := http.NewRequest("POST", engine.GraphAddress+"/query", bytes.NewBufferString(request))
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		resp.Body.Close()

		var details map[string]map[string][]map[string]interface{}
		json.Unmarshal(responseData, &details)

		if len(details["data"]["categories"]) == 0 {
			continue
		}

		categoriesInDatabase := []Category{}

		for _, categoryInDatabase := range details["data"]["categories"] {
			name := categoryInDatabase["name"].(string)
			id := categoryInDatabase["_uid_"].(string)

			id64, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			category := Category{Name: name, ID: id64}
			categoriesInDatabase = append(categoriesInDatabase, category)
		}

		categoriesByName[categoryName] = append(categoriesByName[categoryName], categoriesInDatabase...)
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

	client, request, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Fatal(err)
		return createdCategories, err
	}

	//request.SetObject()

	return createdCategories, nil
}

// CreateCategories is a method for add node for each category in database
func (engine *Engine) CreateCategoriesByHttp(categoriesNames []string) ([]Category, error) {

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

	buf := bytes.NewBufferString(`
		mutation {
			schema {
				name: string @index(exact, term) .
			}

			set {
		`)

	for index, category := range categoriesNames {
		buf.WriteString("_:category-" + strconv.Itoa(index) + " <name> ")
		buf.WriteString("\"" + category + "\"" + " ." + "\n")
	}

	buf.WriteString("}" + " \n" + "}" + "\n")

	req, err := http.NewRequest("POST", engine.GraphAddress+"/query", buf)
	if err != nil {
		log.Fatal(err)
		return createdCategories, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return createdCategories, err
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return createdCategories, err
	}

	var details GraphResponse
	json.Unmarshal(responseData, &details)

	if details.Data.Code == "ErrorInvalidRequest" {
		return createdCategories, ErrCategoriesCanBeCreated
	}

	if details.Data.Message == "Done" {
		for index, name := range categoriesNames {
			idOfCreatedCategory := details.Data.Uids["category-"+strconv.Itoa(index)]
			id64, err := strconv.ParseUint(idOfCreatedCategory, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			category := Category{
				Name: name,
				ID:   id64,
			}
			createdCategories = append(createdCategories, category)
		}
	}

	return createdCategories, nil
}
