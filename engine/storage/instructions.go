package storage

import (
	"context"
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
	ID       string `json:"uid, omitempty"`
	Language string `json:"instructionLanguage, omitempty"`

	IsActive bool `json:"instructionIsActive, omitempty"`

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
