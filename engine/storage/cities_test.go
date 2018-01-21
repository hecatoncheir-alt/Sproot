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

func TestIntegrationCitiesCanBeReadByName(test *testing.T) {
	once.Do(prepareStorage)

	cityForTest := City{Name: "Test city"}

	citiesFromStore, err := storage.Cities.ReadCitiesByName(cityForTest.Name, ".")
	if err != ErrCitiesByNameNotFound {
		test.Fail()
	}

	if citiesFromStore != nil {
		test.Fail()
	}

	createdCity, err := storage.Cities.CreateCity(cityForTest, "en")
	if err != nil || createdCity.ID == "" {
		test.Fail()
	}

	defer storage.Cities.DeleteCity(createdCity)

	citiesFromStore, err = storage.Cities.ReadCitiesByName(createdCity.Name, "en")
	if err != nil {
		test.Fail()
	}

	if len(citiesFromStore) > 1 {
		test.Fail()
	}

	if citiesFromStore[0].Name != createdCity.Name {
		test.Fail()
	}

	if citiesFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationCityCanBeReadById(test *testing.T) {
	once.Do(prepareStorage)

	cityFromStore, err := storage.Cities.ReadCityByID("0", ".")
	if err != ErrCityDoesNotExist {
		test.Fail()
	}

	cityForCreate := City{Name: "Moscow"}

	createdCity, err := storage.Cities.CreateCity(cityForCreate, "en")
	if err != nil {
		test.Error(err)
	}

	defer storage.Cities.DeleteCity(createdCity)

	cityFromStore, err = storage.Cities.ReadCityByID(createdCity.ID, ".")
	if err != nil {
		test.Fail()
	}

	if cityFromStore.Name != createdCity.Name {
		test.Fail()
	}

	if cityFromStore.ID == "" {
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

func TestIntegrationCityCanHasNameWithManyLanguages(test *testing.T) {
	once.Do(prepareStorage)

	createdCity, _ := storage.Cities.CreateCity(City{Name: "Test city"}, "en")
	defer storage.Cities.DeleteCity(createdCity)

	err := storage.Cities.AddLanguageOfCityName(createdCity.ID, "Тестовый город", "ru")
	if err != nil {
		test.Fail()
	}

	cityWithEnName, _ := storage.Cities.ReadCityByID(createdCity.ID, "en")
	if cityWithEnName.Name != "Test city" {
		test.Fail()
	}

	cityWithRuName, _ := storage.Cities.ReadCityByID(createdCity.ID, "ru")
	if cityWithRuName.Name != "Тестовый город" {
		test.Fail()
	}
}