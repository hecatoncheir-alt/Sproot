package engine

import (
	"testing"
	"reflect"
)

func TestIntegrationCategoriesCanBeDeleted(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("http", "192.168.99.100", 8080)
	if err != nil {
		test.Error(err)
	}

	testCategories := []string{"First test category", "Second test category"}
	createdCategories, err := puffer.CreateCategories(testCategories)
	if err != nil {
		if err != ErrCategoriesAlreadyExists {
			test.Error(err)
		}
	}

	if len(createdCategories) <= 2 {
		test.Fail()
	}

	deletedCategories, err := puffer.DeleteCategories(createdCategories)

	if err != nil {
		test.Error(err)
	}

	if len(deletedCategories) <= 2 {
		test.Fail()
	}

}

func TestIntegrationCategoriesCanBeCreated(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("http", "192.168.99.100", 8080)
	if err != nil {
		test.Error(err)
	}

	testCategories := []string{"First test category", "Second test category"}
	categories, err := puffer.CreateCategories(testCategories)
	if err != nil {
		if err != ErrCategoriesAlreadyExists {
			test.Error(err)
		}
	}

	if len(categories) <= 2 {
		test.Fail()
	}

	if categories[0].Name != "First test category" {
		test.Fail()
	}

	if categories[0].ID == "" {
		test.Fail()
	}

	puffer.DeleteCategories(categories)
}

func TestIntegrationCategoriesCanBeRead(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("http", "192.168.99.100", 8080)
	if err != nil {
		test.Error(err)
	}

	testCategories := []string{"First test category", "Second test category"}
	createdCategories, err := puffer.CreateCategories(testCategories)
	if err != nil {
		if err != ErrCategoriesAlreadyExists {
			test.Error(err)
		}
	}

	readCategories, err := puffer.ReadCategoriesByName(testCategories)
	if err != nil {
		test.Error(err)
	}

	if reflect.DeepEqual(createdCategories, readCategories) != true {
		test.Fail()
	}

	puffer.DeleteCategories(createdCategories)
}
