package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sort"
	"fmt"
)

// ErrCategoriesAlreadyExists means that the categories is in the database already
var ErrCategoriesAlreadyExists = errors.New("categories already exists")

// ErrCategoriesCanBeCreated means that the categories can't be added to database
var ErrCategoriesCanBeCreated = errors.New("categories can't be created")

// ErrCategoriesCantBeDeleted means that the categories is in the database already
var ErrCategoriesCantBeDeleted = errors.New("categories can't be deleted")

func (engine *Engine) DeleteCategories(categories []Category) ([]Category, error) {
	request := bytes.NewBufferString(`
		mutation {
			delete {
	`)

	for _, category := range categories {
		request.WriteString("<" + category.ID + "> * * .\n")
	}

	request.WriteString("}\n" + "}\n")
	fmt.Println(request.String())

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

		fmt.Println(request)

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

		if details["data"]["Code"] == nil {
			continue
		}

		categoriesInDatabase := []Category{}

		for _, categoryInDatabase := range details["data"]["categories"] {
			name := categoryInDatabase["name"].(string)
			id := categoryInDatabase["_uid_"].(string)

			category := Category{Name: name, ID: id}
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
			category := Category{
				Name: name,
				ID:   idOfCreatedCategory,
			}
			createdCategories = append(createdCategories, category)
		}
	}

	return createdCategories, nil
}

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/cayleygraph/cayley"
// 	"github.com/cayleygraph/cayley/quad"
// )

// // GetCategoriesOfCompany is a method for get all categories of company
// func (engine *Engine) GetCategoriesOfCompany(companyName string) (categories []string, err error) {
// 	// var err error
// 	companyName = strings.ToLower(companyName)
// 	// it := iterator.NewAnd(engine.Store,
// 	// 	engine.Store.QuadIterator(quad.Object, engine.Store.ValueOf(quad.String(companyName))),
// 	// 	engine.Store.QuadIterator(quad.Predicate, engine.Store.ValueOf(quad.String("belongs"))))

// 	// defer it.Close()

// 	// for it.Next() {
// 	// 	f := engine.Store.Quad(it.Result()).String()
// 	// 	fmt.Println(f)
// 	// }

// 	path := cayley.StartPath(engine.Store, quad.String(companyName)).LabelContext("Category").In("belongs")

// 	path.Iterate(nil).EachValue(engine.Store, func(value quad.Value) {
// 		categories = append(categories, value.String())
// 	})

// 	return categories, nil
// }

// // DeleteCategoriesOfCompany is method for delete categories from company
// func (engine *Engine) DeleteCategoriesOfCompany(categories []string, companyName string) error {
// 	var err error
// 	var store *cayley.Handle
// 	store = engine.Store

// 	companyName = strings.ToLower(companyName)
// 	c, _ := engine.GetCategoriesOfCompany(companyName)
// 	fmt.Println(c)

// 	for _, category := range categories {
// 		fmt.Println(category)
// 		for _, direction := range []quad.Direction{quad.Subject, quad.Predicate} {
// 			it := store.QuadIterator(direction, store.ValueOf(quad.String(companyName)))
// 			for it.Next() {
// 				store.RemoveQuad(store.Quad(it.Result()))
// 			}
// 			it.Close()
// 		}

// 	}

// 	// it := iterator.NewAnd(engine.Store,
// 	// 	engine.Store.QuadIterator(quad.Predicate, engine.Store.ValueOf(quad.String("belongs"))),
// 	// 	engine.Store.QuadIterator(quad.Object, engine.Store.ValueOf(quad.String(companyName))))

// 	// defer it.Close()

// 	// for it.Next() {
// 	// 	res := engine.Store.Quad(it.Result()).String()
// 	// 	subject := strings.Split(res, "--")[0]

// 	// 	for _, category := range categories {
// 	// 		da := strings.EqualFold(subject, category)
// 	// 		fmt.Println(da)

// 	// 		if category == subject {
// 	// 			fmt.Println("da")
// 	// 			engine.Store.RemoveQuad(engine.Store.Quad(it.Result()))
// 	// 		}
// 	// 	}
// 	// }

// 	c, _ = engine.GetCategoriesOfCompany(companyName)
// 	fmt.Println(c)

// 	// sort.Strings()

// 	return err
// }

// // SaveCategoriesOfCompany method for add categories to company
// func (engine *Engine) SaveCategoriesOfCompany(categories []string, companyName string) error {
// 	var err error
// 	companyName = strings.ToLower(companyName)

// 	_, err = engine.GetCompany(companyName)
// 	if err != ErrCompanyNotExists {
// 		return err
// 	}

// 	// TODO: Нужно получить список категорий и добавлять только нужные
// 	for _, category := range categories {
// 		transaction := cayley.NewTransaction()
// 		transaction.AddQuad(cayley.Quad(category, "is", "Category name", "Category"))
// 		transaction.AddQuad(cayley.Quad(category, "belongs", companyName, "Category"))
// 		engine.Store.ApplyTransaction(transaction)
// 	}

// 	return nil
// }
