package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
	"log"
)

type PageInstruction struct {
	ID                       string `json:"uid, omitempty"`
	Path                     string `json:"path, omitempty"`
	PageInPaginationSelector string `json:"pageInPaginationSelector, omitempty"`
	PageParamPath            string `json:"pageParamPath, omitempty"`
	CityParamPath            string `json:"cityParamPath, omitempty"`
	CityParam                string `json:"cityParam, omitempty"`
	ItemSelector             string `json:"itemSelector, omitempty"`
	NameOfItemSelector       string `json:"nameOfItemSelector, omitempty"`
	PriceOfItemSelector      string `json:"priceOfItemSelector, omitempty"`
}

// Instruction is a structure of instruction for parse
type Instruction struct {
	ID         string            `json:"uid, omitempty"`
	Language   string            `json:"instructionLanguage, omitempty"`
	IsActive   bool              `json:"instructionIsActive, omitempty"`
	Pages      []PageInstruction `json:"has_pages, omitempty"`
	Cities     []City            `json:"has_city, omitempty"`
	Companies  []Company         `json:"has_company, omitempty"`
	Categories []Categories      `json:"has_category, omitempty"`
}

// NewPricesResourceForStorage is a constructor of Prices resource
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
		instructionLanguage: string @index(exact, term) .
		instructionIsActive: bool @index(bool) .
		has_company: uid @count .
		has_city: uid @count .
		has_pages: uid @count .

		path: string @index(exact, term) .
		pageInPaginationSelector: string @index(exact, term) .
		pageParamPath: string @index(exact, term) .
		pageCityPath: string @index(exact, term) .
		cityParam: string @index(exact, term) .
		itemSelector: string @index(exact, term) .
		nameOfItemSelector: string @index(exact, term) .
		priceOfItemSelector: string @index(exact, term) .
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
					cityParam
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

	return instruction, nil
}

// ErrInstructionDoesNotExist means than the instruction does not exist in database
var ErrInstructionDoesNotExist = errors.New("instruction by id not found")

func (resource *Instructions) ReadInstructionByID(instructionID, language string) (Instruction, error) {

	query := fmt.Sprintf(`{
				instructions(func: uid("%s")) @filter(has(instructionLanguage)) {
					uid
					instructionLanguage
					instructionIsActive
					has_pages @filter(eq(cityIsActive, true)) {
						uid
						path
						pageInPaginationSelector
						pageParamPath
						cityParamPath
						cityParam
						itemSelector
						nameOfItemSelector
						priceOfItemSelector
					}
					has_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@%v
						cityIsActive
					}
					has_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@%v
						companyIsActive
					}
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@%v
						categoryIsActive
					}
				}
			}`, instructionID, language, language, language)

	instruction := Instruction{ID: instructionID}

	transaction := resource.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	type InstructionsInStorage struct {
		Instructions []Instruction `json:"instructions"`
	}

	var foundedPageInstructions InstructionsInStorage

	err = json.Unmarshal(response.GetJson(), &foundedPageInstructions)
	if err != nil {
		log.Println(err)
		return instruction, err
	}

	if len(foundedPageInstructions.Instructions) == 0 {
		return instruction, ErrInstructionDoesNotExist
	}

	return foundedPageInstructions.Instructions[0], nil
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

//func (instructions *Instructions) AddPageInstructionToInstruction(instructionID, pageInstructionID string) error {
//
//}
//
//func (instructions *Instructions) RemovePageInstructionFromInstruction(instructionID, pageInstructionID string) error {
//
//}
//
//func (instructions *Instructions) AddCategoryToInstruction(instructionID, categoryID string) error {
//
//}
//
//func (instructions *Instructions) RemoveCategoryFromInstruction(instructionID, categoryID string) error {
//
//}
