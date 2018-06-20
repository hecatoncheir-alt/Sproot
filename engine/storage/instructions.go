package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"bytes"
	"text/template"

	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
)

// PageInstruction is a structure for parser of products
type PageInstruction struct {
	ID                         string `json:"uid,omitempty"`
	Path                       string `json:"path,omitempty"`
	PageInPaginationSelector   string `json:"pageInPaginationSelector,omitempty"`
	PreviewImageOfItemSelector string `json:"previewImageOfSelector,omitempty"`
	PageParamPath              string `json:"pageParamPath,omitempty"`
	CityParamPath              string `json:"cityParamPath,omitempty"`
	ItemSelector               string `json:"itemSelector,omitempty"`
	NameOfItemSelector         string `json:"nameOfItemSelector,omitempty"`
	LinkOfItemSelector         string `json:"linkOfItemSelector,omitempty"`
	CityInCookieKey            string `json:"cityInCookieKey,omitempty"`
	CityIDForCookie            string `json:"cityIdForCookie,omitempty"`
	PriceOfItemSelector        string `json:"priceOfItemSelector,omitempty"`
}

// Instruction is a structure of instruction for parse
type Instruction struct {
	ID               string            `json:"uid,omitempty"`
	Language         string            `json:"instructionLanguage,omitempty"`
	IsActive         bool              `json:"instructionIsActive"`
	PagesInstruction []PageInstruction `json:"has_page,omitempty"`
	Cities           []City            `json:"has_city,omitempty"`
	Companies        []Company         `json:"has_company,omitempty"`
	Categories       []Category        `json:"has_category,omitempty"`
}

// NewInstructionsResourceForStorage is a constructor of Prices resource
func NewInstructionsResourceForStorage(storage *Storage) *Instructions {
	return &Instructions{storage: storage}
}

// Instructions is resource of storage for CRUD operations
type Instructions struct {
	storage *Storage
}

// SetUp is a method of Instructions resource for prepare database client and schema.
func (resource *Instructions) SetUp() (err error) {
	schema := `
		instructionLanguage: string @index(term) .
		instructionIsActive: bool @index(bool) .
		has_company: uid @count .
		has_city: uid @count .
		has_page: uid @count .

		path: string @index(term) .
		pageInPaginationSelector: string @index(term) .
		previewImageOfSelector: string @index(term) .
		pageParamPath: string @index(term) .
		pageCityPath: string @index(term) .
		itemSelector: string @index(term) .
		nameOfItemSelector: string @index(term) .
		priceOfItemSelector: string @index(term) .
		cityInCookieKey: string @index(term) .
		cityIdForCookie: string @index(term) .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = resource.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (resource *Instructions) CreatePageInstruction(pageInstruction PageInstruction) (PageInstruction, error) {
	transaction := resource.storage.Client.NewTxn()

	encodedPageInstruction, err := json.Marshal(pageInstruction)
	if err != nil {
		log.Println(err)
		return pageInstruction, err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedPageInstruction,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return pageInstruction, err
	}

	pageInstruction.ID = assigned.Uids["blank-0"]

	return pageInstruction, nil
}

// ErrPageInstructionDoesNotExist means than the page instruction does not exist in database
var ErrPageInstructionDoesNotExist = errors.New("page instruction by id not found")

func (resource *Instructions) ReadPageInstructionByID(pageInstructionID string) (PageInstruction, error) {
	pageInstruction := PageInstruction{ID: pageInstructionID}

	query := fmt.Sprintf(`{
				pageInstructions(func: uid("%s")) @filter(has(path)) {
					uid
					path
					pageInPaginationSelector
					pageParamPath
					cityParamPath
					itemSelector
					nameOfItemSelector
					priceOfItemSelector
				}
			}`, pageInstructionID)

	transaction := resource.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return pageInstruction, err
	}

	type PageInstructionsInStorage struct {
		PageInstructions []PageInstruction `json:"pageInstructions"`
	}

	var foundedPageInstructions PageInstructionsInStorage

	err = json.Unmarshal(response.GetJson(), &foundedPageInstructions)
	if err != nil {
		log.Println(err)
		return pageInstruction, err
	}

	if len(foundedPageInstructions.PageInstructions) == 0 {
		return pageInstruction, ErrPageInstructionDoesNotExist
	}

	return foundedPageInstructions.PageInstructions[0], nil
}

func (resource *Instructions) DeletePageInstruction(pageInstruction PageInstruction) (string, error) {
	deletePageInstructionData, err := json.Marshal(map[string]string{"uid": pageInstruction.ID})
	if err != nil {
		log.Println(err)
		return pageInstruction.ID, err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deletePageInstructionData,
		CommitNow:  true}

	transaction := resource.storage.Client.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return pageInstruction.ID, err
	}

	return pageInstruction.ID, nil
}

func (resource *Instructions) CreateInstructionForCompany(companyID, language string) (Instruction, error) {
	instruction := Instruction{IsActive: true, Language: language}

	transaction := resource.storage.Client.NewTxn()

	encodedInstruction, err := json.Marshal(instruction)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedInstruction,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	instruction.ID = assigned.Uids["blank-0"]

	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instruction.ID, "has_company", companyID)
	mutation = &dataBaseAPI.Mutation{
		SetNquads: []byte(predicate),
		CommitNow: true}

	transaction = resource.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return instruction, err
	}

	updatedInstruction, err := resource.ReadInstructionByID(instruction.ID, language)
	if err != nil {
		return updatedInstruction, err
	}

	return updatedInstruction, nil
}

// ErrInstructionDoesNotExist means than the instruction does not exist in database
var ErrInstructionDoesNotExist = errors.New("instruction by id not found")

func (resource *Instructions) ReadInstructionByID(instructionID, language string) (Instruction, error) {
	variables := struct {
		InstructionID string
		Language      string
	}{
		InstructionID: instructionID,
		Language:      language}

	queryTemplate, err := template.New("ReadInstructionByID").Parse(`{
				instructions(func: uid("{{.InstructionID}}")) @filter(has(instructionLanguage)) {
					uid
					instructionLanguage
					instructionIsActive
					has_page {
						uid
						path
						pageInPaginationSelector
						pageParamPath
						cityParamPath
						itemSelector
						nameOfItemSelector
						priceOfItemSelector
					}
					has_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					has_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIsActive
					}
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
					}
				}
			}`)

	instruction := Instruction{ID: instructionID}

	if err != nil {
		log.Println(err)
		return instruction, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	transaction := resource.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	type InstructionsInStorage struct {
		Instructions []Instruction `json:"instructions"`
	}

	var foundedInstructions InstructionsInStorage

	err = json.Unmarshal(response.GetJson(), &foundedInstructions)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	if len(foundedInstructions.Instructions) == 0 {
		return instruction, ErrInstructionDoesNotExist
	}

	return foundedInstructions.Instructions[0], nil
}

func (resource *Instructions) DeleteInstruction(instruction Instruction) (string, error) {
	deleteInstructionData, err := json.Marshal(map[string]string{"uid": instruction.ID})
	if err != nil {
		log.Println(err)
		return instruction.ID, err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteInstructionData,
		CommitNow:  true}

	transaction := resource.storage.Client.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return instruction.ID, err
	}

	return instruction.ID, nil
}

// ErrCityCanNotBeAddedToInstruction means that the city can't be added to instruction
var ErrCityCanNotBeAddedToInstruction = errors.New("city can not be added to instruction")

func (resource *Instructions) AddCityToInstruction(instructionID, cityID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_city", cityID)
	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCityCanNotBeAddedToInstruction
	}

	return nil
}

// ErrCityCanNotBeRemovedFromInstruction means that the city can't be removed from instruction
var ErrCityCanNotBeRemovedFromInstruction = errors.New("city can not be removed from instruction")

func (resource *Instructions) RemoveCityFromInstruction(instructionID, cityID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_city", cityID)
	mutation := dataBaseAPI.Mutation{
		DelNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCityCanNotBeRemovedFromInstruction
	}

	return nil
}

// ErrPageInstructionCanNotBeAddedToInstruction means that the page instruction can't be added to instruction
var ErrPageInstructionCanNotBeAddedToInstruction = errors.New("page instruction can not be added to instruction")

func (resource *Instructions) AddPageInstructionToInstruction(instructionID, pageInstructionID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_page", pageInstructionID)
	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrPageInstructionCanNotBeAddedToInstruction
	}

	return nil
}

// ErrPageInstructionCanNotBeRemovedFromInstruction means that the page instruction can't be removed from instruction
var ErrPageInstructionCanNotBeRemovedFromInstruction = errors.New("page instruction can not be removed from instruction")

func (resource *Instructions) RemovePageInstructionFromInstruction(instructionID, pageInstructionID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_page", pageInstructionID)
	mutation := dataBaseAPI.Mutation{
		DelNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrPageInstructionCanNotBeRemovedFromInstruction
	}

	return nil
}

// ErrCategoryCanNotBeAddedToInstruction means that the category can't be added to instruction
var ErrCategoryCanNotBeAddedToInstruction = errors.New("category can not be added to instruction")

func (resource *Instructions) AddCategoryToInstruction(instructionID, categoryID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_category", categoryID)
	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeAddedToInstruction
	}

	return nil
}

func (resource *Instructions) RemoveCategoryFromInstruction(instructionID, categoryID string) error {
	predicate := fmt.Sprintf(`<%s> <%s> <%s> .`, instructionID, "has_category", categoryID)
	mutation := dataBaseAPI.Mutation{
		DelNquads: []byte(predicate),
		CommitNow: true}

	transaction := resource.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCategoryCanNotBeAddedToInstruction
	}

	return nil
}

var ErrInstructionsForCompanyDoesNotExist = errors.New("instructions can not be founded for company")

func (resource *Instructions) ReadAllInstructionsForCompany(companyID, language string) ([]Instruction, error) {

	variables := struct {
		CompanyID string
		Language  string
	}{
		CompanyID: companyID,
		Language:  language}

	queryTemplate, err := template.New("ReadAllInstructionsForCompany").Parse(`{
				instructions(func: has(has_company))
				@filter(eq(instructionIsActive, true) AND uid_in(has_company, {{.CompanyID}})) {
					uid
					instructionLanguage
					instructionIsActive
					has_page {
						uid
						path
						pageInPaginationSelector
						pageParamPath
						cityParamPath
						itemSelector
						nameOfItemSelector
						priceOfItemSelector
					}
					has_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIsActive
					}
					has_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	transaction := resource.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	type InstructionsFromStorage struct {
		Instructions []Instruction `json:"instructions"`
	}

	var foundedInstructions InstructionsFromStorage

	err = json.Unmarshal(response.GetJson(), &foundedInstructions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(foundedInstructions.Instructions) == 0 {
		return nil, ErrInstructionsForCompanyDoesNotExist
	}

	return foundedInstructions.Instructions, nil
}

//TODO
//func (resource *Instructions) ReadInstructionOfCategoryForCompany(companyID, categoryID, language string) ([]Instruction, error) {
//
//}
