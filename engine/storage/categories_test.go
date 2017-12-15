package storage

import (
	"testing"
)

func TestIntegrationCategoryCanBeCreated(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(&categoryForCreate)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(*createdCategory)

	if createdCategory.ID == "" {
		test.Fail()
	}

	existCategory, err := storage.Categories.CreateCategory(&categoryForCreate)
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	if existCategory.ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeReadByName(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForSearch := Category{Name: "Test category"}

	categoriesFromStore, err := storage.Categories.ReadCategoriesByName(categoryForSearch.Name)
	if err != ErrCategoriesByNameNotFound {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := storage.Categories.CreateCategory(&categoryForSearch)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(*createdCategory)

	categoriesFromStore, err = storage.Categories.ReadCategoriesByName(createdCategory.Name)
	if err != nil {
		test.Fail()
	}

	if categoriesFromStore[0].Name != createdCategory.Name {
		test.Fail()
	}

	if categoriesFromStore[0].ID == "" {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeReadById(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForSearch := Category{Name: "Test category"}

	categoriesFromStore, err := storage.Categories.ReadCategoriesByName(categoryForSearch.Name)
	if err != ErrCategoriesByNameNotFound {
		test.Fail()
	}

	categoryFromStore, err := storage.Categories.ReadCategoryByID("0")
	if err != ErrCategoryByIDNotFound {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := storage.Categories.CreateCategory(&categoryForSearch)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(*createdCategory)

	categoryFromStore, err = storage.Categories.ReadCategoryByID(createdCategory.ID)
	if err != nil {
		test.Fail()
	}

	if categoryFromStore.Name != createdCategory.Name {
		test.Fail()
	}

	if categoryFromStore.ID == "" {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeUpdated(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("192.168.99.100", 9080)
	if err != nil {
		test.Error(err)
	}

	err = puffer.SetUpIndexes()
	if err != nil {
		test.Error(err)
	}

	categoryForCreate := Category{Name: "First test category"}

	updatedCategory, err := puffer.UpdateCategory(&Category{Name: "Updated test category"})
	if err != nil {
		if err != ErrCategoryWithoutID {
			test.Error(err)
		}
	}

	createdCategory, err := puffer.CreateCategory(categoryForCreate.Name)
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	defer puffer.DeleteCategory(createdCategory)

	categoryForUpdate := Category{ID: createdCategory.ID, Name: "Updated test category"}
	updatedCategory, err = puffer.UpdateCategory(&categoryForUpdate)
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryForUpdate.Name {
		test.Fail()
	}

	categoryInStore, err := puffer.ReadCategoryByID(createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryInStore.Name {
		test.Fail()
	}
}

//func TestIntegrationCategoryCanBeDeleted(test *testing.T) {
//	var err error
//	puffer := New()
//
//	err = puffer.DatabaseSetUp("192.168.99.100", 9080)
//	if err != nil {
//		test.Error(err)
//	}
//
//	err = puffer.SetUpIndexes()
//	if err != nil {
//		test.Error(err)
//	}
//
//	createdCategory, err := puffer.CreateCategory("First test category")
//	if err != nil {
//		if err != ErrCategoryAlreadyExist {
//			test.Error(err)
//		}
//	}
//
//	deletedCategoryID, err := puffer.DeleteCategory(createdCategory)
//
//	if err != nil {
//		test.Error(err)
//	}
//
//	if deletedCategoryID != createdCategory.ID {
//		test.Fail()
//	}
//
//	_, err = puffer.ReadCategoryByID(deletedCategoryID)
//	if err != ErrCategoryDoesNotExist {
//		test.Fail()
//	}
//}
