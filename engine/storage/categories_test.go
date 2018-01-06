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
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	if createdCategory.ID == "" {
		test.Fail()
	}

	existCategory, err := storage.Categories.CreateCategory(categoryForCreate)
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

	createdCategory, err := storage.Categories.CreateCategory(categoryForSearch)
	if err != nil || createdCategory.ID == "" {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

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
	if err != ErrCategoryDoesNotExist {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := storage.Categories.CreateCategory(categoryForSearch)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

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
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	updatedCategory, err := storage.Categories.UpdateCategory(Category{Name: "Updated test category"})
	if err != nil {
		if err != ErrCategoryCanNotBeWithoutID {
			test.Error(err)
		}
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate)
	if err != nil {
		test.Error(err)
	}

	defer storage.Categories.DeleteCategory(createdCategory)

	categoryForUpdate := Category{ID: createdCategory.ID, Name: "Updated test category", IsActive: createdCategory.IsActive}
	updatedCategory, err = storage.Categories.UpdateCategory(categoryForUpdate)
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryForUpdate.Name {
		test.Fail()
	}

	categoryInStore, err := storage.Categories.ReadCategoryByID(createdCategory.ID)
	if err != nil {
		test.Error(err)
	}

	if updatedCategory.Name != categoryInStore.Name {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeDeactivate(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate)
	if err != nil {
		test.Error(err)
	}

	deactivatedCategoryID, err := storage.Categories.DeactivateCategory(createdCategory)
	if err != nil {
		test.Error(err)
	}

	if deactivatedCategoryID != createdCategory.ID {
		test.Fail()
	}

	_, err = storage.Categories.ReadCategoryByID(deactivatedCategoryID)
	if err != ErrCategoryDoesNotExist {
		test.Error(err)
	}
}

func TestIntegrationCategoryCanBeDeleted(test *testing.T) {
	var err error
	storage := New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		test.Error(err)
	}

	categoryForCreate := Category{Name: "Test category"}
	createdCategory, err := storage.Categories.CreateCategory(categoryForCreate)
	if err != nil {
		test.Error(err)
	}

	deletedCategoryID, err := storage.Categories.DeleteCategory(createdCategory)
	if err != nil {
		test.Error(err)
	}

	if deletedCategoryID != createdCategory.ID {
		test.Fail()
	}

	_, err = storage.Categories.ReadCategoryByID(deletedCategoryID)
	if err != ErrCategoryDoesNotExist {
		test.Error(err)
	}
}
