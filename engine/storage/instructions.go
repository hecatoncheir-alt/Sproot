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
	ID        string            `json:"uid, omitempty"`
	Language  string            `json:"instructionLanguage, omitempty"`
	IsActive  bool              `json:"instructionIsActive, omitempty"`
	Pages     []PageInstruction `json:"has_pages, omitempty"`
	Cities    []City            `json:"has_city, omitempty"`
	Companies []Company         `json:"has_company, omitempty"`
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
func (instructions *Instructions) SetUp() (err error) {
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

	err = instructions.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (instructions *Instructions) CreatePageInstruction(pageInstruction PageInstruction) (PageInstruction, error) {
	transaction := instructions.storage.Client.NewTxn()

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

func (instructions *Instructions) ReadPageInstructionByID(pageInstructionID string) (PageInstruction, error) {
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

	transaction := instructions.storage.Client.NewTxn()
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

func (instructions *Instructions) DeletePageInstruction(pageInstruction PageInstruction) (string, error) {
	deletePageInstructionData, err := json.Marshal(map[string]string{"uid": pageInstruction.ID})
	if err != nil {
		log.Println(err)
		return pageInstruction.ID, err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deletePageInstructionData,
		CommitNow:  true}

	transaction := instructions.storage.Client.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return pageInstruction.ID, err
	}

	return pageInstruction.ID, nil
}
