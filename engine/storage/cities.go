package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
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
		cityName: string @index(term) .
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

var (
	// ErrCityCanNotBeCreated means that the city can't be added to database
	ErrCityCanNotBeCreated = errors.New("city can't be created")

	// ErrCityAlreadyExist means that the city is in the database already
	ErrCityAlreadyExist = errors.New("city already exist")
)

// CreateCity make category and save it to storage
func (cities *Cities) CreateCity(city City, language string) (City, error) {
	existsCities, err := cities.ReadCitiesByName(city.Name, language)
	if err != nil && err != ErrCitiesByNameNotFound {
		log.Println(err)
		return city, ErrCityCanNotBeCreated
	}
	if existsCities != nil {
		return existsCities[0], ErrCityAlreadyExist
	}

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

var (
	// ErrCitiesByNameCanNotBeFound means that the cities can't be found in database
	ErrCitiesByNameCanNotBeFound = errors.New("cities by name can not be found")

	// ErrCitiesByNameNotFound means than the cities does not exist in database
	ErrCitiesByNameNotFound = errors.New("cities by name not found")
)

// ReadALlCities is a method for get all nodes
func (cities *Cities) ReadAllCities(language string) ([]City, error) {
	query := fmt.Sprintf(`{
				cities(func: eq(cityIsActive, true)) @filter(has(cityName)) {
					uid
					cityName: cityName@%v
					cityIsActive
				}
			}`, language)

	transaction := cities.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	type citiesInStorage struct {
		AllCitiesFoundedByName []City `json:"cities"`
	}

	var foundedCities citiesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCities)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	if len(foundedCities.AllCitiesFoundedByName) == 0 {
		return nil, ErrCitiesByNameNotFound
	}

	return foundedCities.AllCitiesFoundedByName, nil
}

// ReadCitiesByName is a method for get all nodes by city name
func (cities *Cities) ReadCitiesByName(cityName, language string) ([]City, error) {
	query := fmt.Sprintf(`{
				cities(func: eq(cityName@%v, "%v")) @filter(eq(cityIsActive, true)) {
					uid
					cityName: cityName@%v
					cityIsActive
				}
			}`, language, cityName, language)

	transaction := cities.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	type citiesInStorage struct {
		AllCitiesFoundedByName []City `json:"cities"`
	}

	var foundedCities citiesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCities)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	if len(foundedCities.AllCitiesFoundedByName) == 0 {
		return nil, ErrCitiesByNameNotFound
	}

	return foundedCities.AllCitiesFoundedByName, nil
}

var (
	// ErrCityCanNotBeWithoutID means that city can't be found in storage for make some operation
	ErrCityCanNotBeWithoutID = errors.New("city can not be without id")

	// ErrCityByIDCanNotBeFound means that the city can't be found in database
	ErrCityByIDCanNotBeFound = errors.New("city by id can not be found")

	// ErrCityDoesNotExist means than the city does not exist in database
	ErrCityDoesNotExist = errors.New("city by id not found")
)

// ReadCityByID is a method for get all nodes of categories by ID
func (cities *Cities) ReadCityByID(cityID, language string) (City, error) {
	city := City{ID: cityID}

	if cityID == "" {
		return city, ErrCityCanNotBeWithoutID
	}

	query := fmt.Sprintf(`{
				cities(func: uid("%s")) @filter(has(cityName)) {
					uid
					cityName: cityName@%v
					cityIsActive
				}
			}`, cityID, language)

	transaction := cities.storage.Client.NewTxn()
	response, err := transaction.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	type citiesInStore struct {
		Cities []City `json:"cities"`
	}

	var foundedCities citiesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCities)
	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	if len(foundedCities.Cities) == 0 {
		return city, ErrCityDoesNotExist
	}

	return foundedCities.Cities[0], nil
}

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
