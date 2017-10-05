package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ErrCompanyNotExists means that the company is not in the database
var ErrCompanyNotExists = errors.New("Company not exists")

// ErrCompanyAlreadyExists means that the company is in the database already
var ErrCompanyAlreadyExists = errors.New("Company already exists")

// ErrCompanyCanNotBeDeleted delete all nodes with company predicates
var ErrCompanyCanNotBeDeleted = errors.New("Company can not be deleted")

// // SaveCompany method for add triplets to graph db
func (engine *Engine) CreateCompany(company *Company) (recordID string, err error) {
	if len(company.Categories) > 0 {
		engine.CreateCategories(company.Categories)
	}

	request := fmt.Sprintf(`
		mutation {
			schema {
				name: string @index(exact, term) .
				iri: string @index(exact, term) .
			}

			set {
				_:company <name> "%v" .
				_:company <iri> "%v" .
	`, company.Name, company.IRI)

	body := bytes.NewBufferString(request)

	if len(company.Categories) > 0 {
		for _, category := range company.Categories {
			body.WriteString("_:company <has_category> ")
			body.WriteString(category + " ." + "\n")
		}
	}

	body.WriteString("}" + " \n" + "}" + "\n")

	var uid string

	req, err := http.NewRequest("POST", engine.GraphAddress+"/query", body)
	if err != nil {
		log.Fatal(err)
		return uid, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return uid, err
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return uid, err
	}

	log.Printf("Response %+v\n", string(responseData))

	var details map[string]interface{}
	json.Unmarshal(responseData, &details)

	if details["code"] == "ErrorInvalidRequest" {
		return uid, ErrCompanyCanNotBeDeleted
	}

	return "", nil
}

// GetCompany return company object of company node in graph store
func (engine *Engine) GetCompany(companyName string) (company *Company, err error) {
	return nil, nil
}

// // DeleteCompany method for delete all nodes with company name
func (engine *Engine) DeleteCompany(companyName string) error {
	body := bytes.NewBufferString(`
		mutation {
			set {

			}
		}
	`)

	req, err := http.NewRequest("POST", engine.GraphAddress+"/query", body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Response %+v\n", string(responseData))

	var details map[string]interface{}
	json.Unmarshal(responseData, &details)

	if details["code"] == "ErrorInvalidRequest" {
		return ErrCompanyCanNotBeDeleted
	}

	return nil
}

// // DeleteCompany method for delete all nodes with company name
// func (engine *Engine) DeleteCompany(companyName string) error {
// 	var err error
// 	var categories []string
// 	// var path *path.Path

// 	// regCompanyName, err := regexp.Compile(strings.ToLower(companyName))
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	categories, err = engine.GetCategoriesOfCompany(companyName)
// 	if err != nil {
// 		return err
// 	}

// 	err = engine.DeleteCategoriesOfCompany(categories, companyName)
// 	if err != nil {
// 		return err
// 	}

// 	// categories, _ = engine.GetCategoriesOfCompany(companyName)
// 	// fmt.Println(categories)

// 	// fmt.Println("In")
// 	// path = cayley.StartPath(engine.Store).Regex(regCompanyName).In(quad.String("belongs")).Out()
// 	// path.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 	// 	fmt.Println(value.String())
// 	// })

// 	// fmt.Println("Out")
// 	// path = cayley.StartPath(engine.Store).Regex(regCompanyName).OutPredicates()
// 	// path.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 	// 	fmt.Println(value.String())
// 	// })

// 	return nil
// }

// // GetCompany return company object of company node in graph store
// func (engine *Engine) GetCompany(companyName string) (company *Company, err error) {
// 	var path *path.Path

// 	regCompanyName, err := regexp.Compile(strings.ToLower(companyName))
// 	if err != nil {
// 		return nil, err
// 	}

// 	path = cayley.StartPath(engine.Store).Regex(regCompanyName)

// 	var companyInStore string
// 	path.Iterate(nil).EachValue(nil, func(value quad.Value) {
// 		companyInStore = value.String()
// 	})

// 	if companyInStore == "" {
// 		return nil, ErrCompanyNotExists
// 	}

// 	// ---

// 	return nil, ErrCompanyNotExists
// }

// // SaveCompany method for add triplets to graph db
// func (engine *Engine) SaveCompany(company *Company) (companyInStore *Company, err error) {
// 	_, err = engine.GetCompany(company.Name)
// 	if err != ErrCompanyNotExists {
// 		return nil, err
// 	}

// 	companyName := strings.ToLower(company.Name)
// 	lastChangeTime := time.Now().String()

// 	transaction := cayley.NewTransaction()
// 	transaction.AddQuad(cayley.Quad(companyName, "is", "Company name", "Company"))
// 	transaction.AddQuad(cayley.Quad(companyName, "was updated", lastChangeTime, "Time"))
// 	transaction.AddQuad(cayley.Quad(companyName, "has", company.IRI, "Company link"))
// 	transaction.AddQuad(cayley.Quad(companyName, "has", "Address", "Company"))
// 	transaction.AddQuad(cayley.Quad(companyName, "has many", "Product", "Company"))
// 	err = engine.Store.ApplyTransaction(transaction)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = engine.SaveCategoriesOfCompany(company.Categories, company.Name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	companyInStore, err = engine.GetCompany(company.Name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return companyInStore, nil
// }
