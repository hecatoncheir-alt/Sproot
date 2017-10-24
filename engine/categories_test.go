package engine

import (
	"testing"
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

	if len(createdCategories) < 2 {
		test.Fail()
	}

	deletedCategories, err := puffer.DeleteCategories(createdCategories)

	if err != nil {
		test.Error(err)
	}

	if len(deletedCategories) < 2 {
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

	defer puffer.DeleteCategories(categories)

	if len(categories) < 2 {
		test.Fatal("Created categories count must be two")
	}

	if categories[0].Name != "First test category" {
		test.Fail()
	}

	if categories[0].ID == "" {
		test.Fail()
	}

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

	// Created categories must be deleted
	defer puffer.DeleteCategories(createdCategories)

	if len(createdCategories) < 2 {
		test.Fatal("Created categories count must be two")
	}

	readCategories, err := puffer.ReadCategoriesByName(testCategories)
	if err != nil {
		test.Error(err)
	}

	if len(readCategories) < 2 {
		test.Fatal("Two created categories and two categories must be in database")
	}

	for _, categories := range readCategories {
		for _, readCategory := range categories {
			for _, createdCategory := range createdCategories {
				if readCategory.Name == createdCategory.Name && readCategory.ID != createdCategory.ID {
					test.Fatal("Created category ID must be same as a read category ID")
				}
			}

		}
	}
}
