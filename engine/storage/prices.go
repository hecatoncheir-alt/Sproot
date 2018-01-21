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
	IsActive bool      `json:"priceIsActive, omitempty"`
	Cities   []City    `json:"belongs_to_city, omitempty"`
	Products []Product `json:"belongs_to_product, omitempty"`
}

// NewPricesResourceForStorage is a constructor of Prices resource
func NewPricesResourceForStorage(storage *Storage) *Prices {
	return &Prices{storage: storage}
}

// Prices is resource of storage for CRUD operations
type Prices struct {
	storage *Storage
}

// SetUp is a method of Prices resource for prepare database client and schema.
func (prices *Prices) SetUp() (err error) {
	schema := `
		pricesValue: float @index(float) .
		priceDateTime: dateTime @index(day) .
		priceIsActive: bool @index(bool) .
		belongs_to_city: uid @count .
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
					belongs_to_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@%v
						cityIsActive
					}
				}
			}`, priceID, language, language)

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

// ErrCityCanNotBeAddedToPrice means that the city can't be added to price
var ErrCityCanNotBeAddedToPrice = errors.New("city can not be added to price")

// AddCityToPrice method for set quad of predicate about price and city
func (prices *Prices) AddCityToPrice(priceID, cityID string) error {
	forPricePredicate := fmt.Sprintf(`<%s> <%s> <%s> .`, priceID, "belongs_to_city", cityID)
	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(forPricePredicate),
		CommitNow: true}

	transaction := prices.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return ErrCityCanNotBeAddedToPrice
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

		if len(exportedPrice.Products) > 0 {
			err = prices.AddProductToPrice(exportedPrice.ID, exportedPrice.Products[0].ID)
			if err != nil {
				return err
			}
		}

		if len(exportedPrice.Cities) > 0 {
			err = prices.AddCityToPrice(exportedPrice.ID, exportedPrice.Cities[0].ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ExportJSON method for export all prices belongs to product from database to json
func (prices *Prices) ExportJSON() ([]byte, error) {
	query := fmt.Sprintf(`{
				prices(func: has(belongs_to_product)) {
					uid
					priceValue
					priceDateTime
					priceCity
					priceIsActive
					belongs_to_product {
						uid
						productIsActive
					}
					belongs_to_city {
						uid
						cityIsActive
					}
				}
			}`)

	transaction := prices.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	type PricesInStore struct {
		Prices []Price `json:"prices"`
	}

	var foundedPrices PricesInStore

	err = json.Unmarshal(response.GetJson(), &foundedPrices)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	jsonForExport, err := json.Marshal(foundedPrices)
	if err != nil {
		return nil, err
	}

	return jsonForExport, nil
}
