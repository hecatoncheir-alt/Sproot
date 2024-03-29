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

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Fail()
		}
	}()

	if createdCity.ID == "" {
		test.Fail()
	}
}

// Must be run parallel with TestIntegrationNewPriceWithExistedProductsCanBeCreatedForRightProduct
func TestIntegrationAllCitiesCanBeRead(test *testing.T) {
	test.Parallel()
	once.Do(prepareStorage)

	cityForTest := City{Name: "Test city"}
	createdCity, err := storage.Cities.CreateCity(cityForTest, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Fail()
		}
	}()

	otherCityForTest := City{Name: "Other test city"}
	otherCreatedCity, err := storage.Cities.CreateCity(otherCityForTest, "en")
	if err != nil {
		test.Fail()
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(otherCreatedCity)
		if err != nil {
			test.Fail()
		}
	}()

	citiesFromStore, err := storage.Cities.ReadAllCities("en")
	if err != nil {
		test.Fail()
	}

	if len(citiesFromStore) != 2 {
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

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Fail()
		}
	}()

	citiesFromStore, err = storage.Cities.ReadCitiesByName(createdCity.Name, "en")
	if err != nil {
		test.Fail()
	}

	if len(citiesFromStore) > 1 || len(citiesFromStore) < 1 {
		test.Fatal()
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

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Fail()
		}
	}()

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
	if err != nil {
		test.Error(err)
	}

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

	testCityName := "Test city"
	testCityRuName := "Тестовый город"

	createdCity, err := storage.Cities.CreateCity(City{Name: testCityName}, "en")
	if err != nil {
		test.Error(err)
	}

	defer func() {
		_, err := storage.Cities.DeleteCity(createdCity)
		if err != nil {
			test.Error(err)
		}
	}()

	err = storage.Cities.AddLanguageOfCityName(createdCity.ID, testCityRuName, "ru")
	if err != nil {
		test.Fail()
	}

	cityWithEnName, err := storage.Cities.ReadCityByID(createdCity.ID, "en")
	if err != nil {
		test.Error(err)
	}

	if cityWithEnName.Name != testCityName {
		test.Fail()
	}

	cityWithRuName, err := storage.Cities.ReadCityByID(createdCity.ID, "ru")
	if err != nil {
		test.Error(err)
	}

	if cityWithRuName.Name != testCityRuName {
		test.Fail()
	}
}
