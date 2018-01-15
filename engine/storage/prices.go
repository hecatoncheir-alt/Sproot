package storage

import (
	"context"
	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
	"log"
	"time"
)

// Price is a structure of prices in database
type Price struct {
	ID       string    `json:"uid"`
	Value    float64   `json:"priceValue, omitempty"`
	DateTime time.Time `json:"priceDateTime, omitempty"`
	City     string    `json:"priceCity, omitempty"`
	IsActive bool      `json:"priceIsActive, omitempty"`
	Product  Product   `json:"belongs_to_product, omitempty"`
}

// NewPricesResourceForStorage is a constructor of Prices resource
func NewPricesResourceForStorage(storage *Storage) *Prices {
	return &Prices{storage: storage}
}

// Products is resource of storage for CRUD operations
type Prices struct {
	storage *Storage
}

// SetUp is a method of Companies resource for prepare database client and schema.
func (prices *Prices) SetUp() (err error) {
	schema := `
		pricesValue: float @index(float) .
		priceDateTime: dateTime @index(day) .
		priceCity: string @index(term) .
		priceIsActive: bool @index(bool) .
		belongs_to_product: uid @count .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = prices.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
