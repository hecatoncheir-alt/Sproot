package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// ErrPriceCanNotBeCreated means that the price can't be added to database
var ErrPriceCanNotBeCreated = errors.New("price can't be created")

// CreatePrice is a method for make product and save it to storage
func (prices *Prices) CreatePrice(price Price) (Price, error) {
	transaction := prices.storage.Client.NewTxn()

	price.IsActive = true
	encodedPrice, err := json.Marshal(price)
	if err != nil {
		log.Println(err)
		return price, ErrPriceCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedPrice,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return price, ErrPriceCanNotBeCreated
	}

	price.ID = assigned.Uids["blank-0"]
	if price.ID == "" {
		return price, ErrPriceCanNotBeCreated
	}

	return price, nil
}

// ErrPriceCanNotBeDeleted means that the price can't be deleted from database
var ErrPriceCanNotBeDeleted = errors.New("price can't be deleted")

// DeletePrice method for remove price from database
func (prices *Prices) DeletePrice(price Price) (string, error) {
	var err error

	deletePriceData, err := json.Marshal(map[string]string{"uid": price.ID})
	if err != nil {
		log.Println(err)
		return price.ID, ErrPriceCanNotBeDeleted
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deletePriceData,
		CommitNow:  true}

	transaction := prices.storage.Client.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return price.ID, ErrPriceCanNotBeDeleted
	}

	return price.ID, nil
}

// ErrPriceByIDCanNotBeFound means that the price can't be found in database
var ErrPriceByIDCanNotBeFound = errors.New("price by id can not be found")

// ErrPriceDoesNotExist means than the price does not exist in database
var ErrPriceDoesNotExist = errors.New("price by id not found")

// ReadPriceByID is a method for get all nodes of prices by ID
func (prices *Prices) ReadPriceByID(priceID, language string) (Price, error) {
	price := Price{ID: priceID}

	query := fmt.Sprintf(`{
				prices(func: uid("%s")) {
					uid
					priceValue
					priceDateTime
					priceCity
					priceIsActive
					belongs_to_product @filter(eq(productIsActive, true)) {
						uid
						productName: productName@%v
						productIri
						previewImageLink
						productIsActive
					}
				}
			}`, priceID, language)

	transaction := prices.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return price, ErrPriceByIDCanNotBeFound
	}

	type PricesInStore struct {
		Prices []Price `json:"prices"`
	}

	var foundedPrices PricesInStore

	err = json.Unmarshal(response.GetJson(), &foundedPrices)
	if err != nil {
		log.Println(err)
		return price, ErrPriceByIDCanNotBeFound
	}

	if len(foundedPrices.Prices) == 0 {
		return price, ErrPriceDoesNotExist
	}

	return foundedPrices.Prices[0], nil
}
