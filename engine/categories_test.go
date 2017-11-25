package engine

import (
	"strconv"
	"testing"
)

func TestIntegrationCategoryCanBeCreated(test *testing.T) {
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
	createdCategory, err := puffer.CreateCategory(categoryForCreate.Name)
	if err != nil {
		test.Error(err)
	}

	defer puffer.DeleteCategory(createdCategory)

	if strconv.Itoa(int(createdCategory.ID)) == "" {
		test.Fail()
	}

	if createdCategory.ID == 0 {
		test.Fail()
	}

	existCategory, err := puffer.CreateCategory(categoryForCreate.Name)
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	if existCategory.ID != createdCategory.ID {
		test.Fail()
	}
}

func TestIntegrationCategoryCanBeReadedByName(test *testing.T) {
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

	categoriesFromStore, err := puffer.ReadCategoriesByName(categoryForCreate.Name)
	if err != ErrCategoryDoesNotExist {
		test.Fail()
	}

	if categoriesFromStore != nil {
		test.Fail()
	}

	createdCategory, err := puffer.CreateCategory(categoryForCreate.Name)
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	defer puffer.DeleteCategory(createdCategory)

	categoriesFromStore, err = puffer.ReadCategoriesByName(createdCategory.Name)
	if err != nil {
		test.Fail()
	}

	if categoriesFromStore[0].Name != createdCategory.Name {
		test.Fail()
	}

	if categoriesFromStore[0].ID == 0 {
		test.Fail()
	}

}

func TestIntegrationCategoryCanBeReadedById(test *testing.T) {
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

	createdCategory, err := puffer.CreateCategory(categoryForCreate.Name)
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	defer puffer.DeleteCategory(createdCategory)

	categoryFromStore, err := puffer.ReadCategoryByID(createdCategory.ID)
	if err != nil {
		test.Fail()
	}

	if categoryFromStore.Name != createdCategory.Name {
		test.Fail()
	}

	if categoryFromStore.ID == 0 {
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

func TestIntegrationCategoryCanBeDeleted(test *testing.T) {
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

	createdCategory, err := puffer.CreateCategory("First test category")
	if err != nil {
		if err != ErrCategoryAlreadyExist {
			test.Error(err)
		}
	}

	deletedCategoryID, err := puffer.DeleteCategory(createdCategory)

	if err != nil {
		test.Error(err)
	}

	if deletedCategoryID != createdCategory.ID {
		test.Fail()
	}

	_, err = puffer.ReadCategoryByID(deletedCategoryID)
	if err != ErrCategoryDoesNotExist {
		test.Fail()
	}
}
