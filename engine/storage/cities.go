package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
	"log"
)

// City is a structure of prices in database
type City struct {
	ID       string `json:"uid"`
	Name     string `json:"cityName, omitempty"`
	IsActive bool   `json:"cityIsActive, omitempty"`
}

// NewCitiesResourceForStorage is a constructor of Prices resource
func NewCitiesResourceForStorage(storage *Storage) *Cities {
	return &Cities{storage: storage}
}

// Cities is resource of storage for CRUD operations
type Cities struct {
	storage *Storage
}

// SetUp is a method of Cities resource for prepare database client and schema.
func (cities *Cities) SetUp() (err error) {
	schema := `
		cityName: string @index(exact, term) .
		cityIsActive: bool @index(bool) .
	`
	operation := &dataBaseAPI.Operation{Schema: schema}

	err = cities.storage.Client.Alter(context.Background(), operation)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ErrCityCanNotBeCreated means that the city can't be added to database
var ErrCityCanNotBeCreated = errors.New("city can't be created")

// CreateCategory make category and save it to storage
func (cities *Cities) CreateCity(city City, language string) (City, error) {
	//existsCities, err := cities.ReadCityByName(city.Name, language)
	//if err != nil && err != ErrCitiesByNameNotFound {
	//	log.Println(err)
	//	return city, ErrCityCanNotBeCreated
	//}
	//if existsCities != nil {
	//	return existsCities[0], ErrCityAlreadyExist
	//}

	transaction := cities.storage.Client.NewTxn()

	city.IsActive = true
	encodedCity, err := json.Marshal(city)
	if err != nil {
		log.Println(err)
		return city, ErrCityCanNotBeCreated
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCity,
		CommitNow: true}

	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		log.Println(err)
		return city, ErrCityCanNotBeCreated
	}

	city.ID = assigned.Uids["blank-0"]
	if city.ID == "" {
		return city, ErrCityCanNotBeCreated
	}

	err = cities.AddLanguageOfCityName(city.ID, city.Name, language)
	if err != nil {
		return city, err
	}

	return city, nil
}

// AddLanguageOfCityName is a method for add predicate "cityName" for cityName value with new language
func (cities *Cities) AddLanguageOfCityName(cityID, name, language string) error {
	forCityNamePredicate := fmt.Sprintf(`<%s> <cityName> %s .`, cityID, "\""+name+"\""+"@"+language)

	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(forCityNamePredicate),
		CommitNow: true}

	transaction := cities.storage.Client.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}

// ErrCityCanNotBeWithoutID means that city can't be found in storage for make some operation
var ErrCityCanNotBeWithoutID = errors.New("city can not be without id")

// ErrCityCanNotBeDeleted means that the city can't be removed from database
var ErrCityCanNotBeDeleted = errors.New("city can't be deleted")

// DeleteCity method for remove category from database
func (cities *Cities) DeleteCity(city City) (string, error) {

	if city.ID == "" {
		return "", ErrCityCanNotBeWithoutID
	}

	deleteCategoryData, _ := json.Marshal(map[string]string{"uid": city.ID})

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteCategoryData,
		CommitNow:  true}

	transaction := cities.storage.Client.NewTxn()

	var err error
	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		log.Println(err)
		return city.ID, ErrCityCanNotBeDeleted
	}

	return city.ID, nil
}
