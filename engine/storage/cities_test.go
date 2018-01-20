package storage

import (
	"testing"
)

func TestIntegrationCityCanBeCreated(test *testing.T) {
	once.Do(prepareStorage)

	cityForCreate := City{Name: "Moscow"}
	createdCity, err := storage.Cities.CreateCity(cityForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Cities.DeleteCity(createdCity)

	if createdCity.ID == "" {
		test.Fail()
	}
}

func TestIntegrationCityCanBeDeleted(test *testing.T) {
	once.Do(prepareStorage)

	cityForCreate := City{Name: "Moscow"}

	createdCity, err := storage.Cities.CreateCity(cityForCreate, "en")

	deletedCityID, err := storage.Cities.DeleteCity(createdCity)
	if err != nil {
		test.Error(err)
	}

	if deletedCityID != createdCity.ID {
		test.Fail()
	}

	_, err = storage.Cities.ReadCityByID(deletedCityID, ".")
	if err != ErrCityDoesNotExist {
		test.Error(err)
	}
}
