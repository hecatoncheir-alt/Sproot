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
	Product  []Product `json:"belongs_to_product, omitempty"`
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
				prices(func: uid("%s")) @filter(has(priceValue)) {
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

// ErrProductCanNotBeAddedToPrice means that the product can't be added to price
var ErrProductCanNotBeAddedToPrice = errors.New("product can not be added to price")

// AddProductToPrice method for set quad of predicate about price and product
func (prices *Prices) AddProductToPrice(priceID, productID string) error {
	var err error
	var mutation dataBaseAPI.Mutation

	forPricePredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, priceID, "belongs_to_product", productID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forPricePredicate),
		CommitNow: true}

	transaction := prices.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToPrice
	}

	forProductPredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, productID, "has_price", priceID)
	mutation = dataBaseAPI.Mutation{
		SetNquads: []byte(forProductPredicate),
		CommitNow: true}

	transaction = prices.storage.Client.NewTxn()
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrProductCanNotBeAddedToPrice
	}

	return nil
}

// ImportJSON is a method for add prices to database
func (prices *Prices) ImportJSON(exportedPrices []byte) error {

	type allPrices struct {
		Prices []Price `json:"prices"`
	}

	var allPricesInJSON allPrices

	err := json.Unmarshal(exportedPrices, &allPricesInJSON)
	if err != nil {
		return err
	}

	for _, exportedPrice := range allPricesInJSON.Prices {
		_, err := prices.CreatePrice(exportedPrice)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO
func (prices *Prices) ExportJSON() string {
	var jsonForExport string

	return jsonForExport
}
