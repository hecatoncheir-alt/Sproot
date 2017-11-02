package engine

import (
	"errors"
	"log"
)

// ErrCompanyNotExists means that the company is not in the database
var ErrCompanyNotExists = errors.New("company not exists")

// ErrCompanyAlreadyExists means that the company is in the database already
var ErrCompanyAlreadyExists = errors.New("company already exists")

// ErrCompanyCanNotBeDeleted delete all nodes with company predicates
var ErrCompanyCanNotBeDeleted = errors.New("company can not be deleted")

// CreateCompany method for add triplets to graph db
func (engine *Engine) CreateCompany(company *Company) (createdCompany Company, err error) {
	categoriesForCompany := []Category{}

	if len(company.Categories) > 0 {
		categoriesForCreate := []string{}
		for _, category := range company.Categories {
			if int(category.ID) != 0 {
				categoriesForCompany = append(categoriesForCompany, category)
				break
			}

			categoriesForCreate = append(categoriesForCreate, category.Name)
		}

		categoriesForCompany, err = engine.CreateCategories(categoriesForCreate)
		if err != nil {
			log.Fatal(err)
			return Company{}, err
		}
	}

	//request := fmt.Sprintf(`
	//	mutation {
	//		schema {
	//			name: string @index(exact, term) .
	//			iri: string @index(exact, term) .
	//		}
	//
	//		set {
	//			_:company <name> "%v" .
	//			_:company <iri> "%v" .
	//`, company.Name, company.IRI)
	//
	//body := bytes.NewBufferString(request)
	//
	//if len(company.Categories) > 0 {
	//	for _, category := range company.Categories {
	//		body.WriteString("_:company <has_category> ")
	//		body.WriteString(category + " ." + "\n")
	//	}
	//}
	//
	//body.WriteString("}" + " \n" + "}" + "\n")
	//
	//var uid string
	//
	//req, err := http.NewRequest("POST", engine.GraphAddress+"/query", body)
	//if err != nil {
	//	log.Fatal(err)
	//	return uid, err
	//}
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//	return uid, err
	//}
	//
	//defer resp.Body.Close()
	//
	//responseData, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//	return uid, err
	//}
	//
	//log.Printf("Response %+v\n", string(responseData))
	//
	//var details map[string]interface{}
	//json.Unmarshal(responseData, &details)
	//
	//if details["code"] == "ErrorInvalidRequest" {
	//	return uid, ErrCompanyCanNotBeDeleted
	//}

	return Company{}, nil
}

// ReadCompany return company object of company node in graph store
func (engine *Engine) ReadCompany(companyName string) (company Company, err error) {
	return company, nil
}

// UpdateCompany return company object of company node in graph store
func (engine *Engine) UpdateCompany(company *Company) (updatedCompany Company, err error) {
	return updatedCompany, nil
}

// DeleteCompany method for delete all nodes with company name
func (engine *Engine) DeleteCompany(company *Company) (deletedRecordId uint64, err error) {
	//body := bytes.NewBufferString(`
	//	mutation {
	//		set {
	//
	//		}
	//	}
	//`)
	//
	//req, err := http.NewRequest("POST", engine.GraphAddress+"/query", body)
	//if err != nil {
	//	log.Fatal(err)
	//	return ,err
	//}
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//	return err
	//}
	//
	//defer resp.Body.Close()
	//
	//responseData, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//	return err
	//}
	//
	//log.Printf("Response %+v\n", string(responseData))
	//
	//var details map[string]interface{}
	//json.Unmarshal(responseData, &details)
	//
	//if details["code"] == "ErrorInvalidRequest" {
	//	return ErrCompanyCanNotBeDeleted
	//}
	//
	//return nil
	return deletedRecordId, nil
}
