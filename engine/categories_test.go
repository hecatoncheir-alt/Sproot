package engine

import (
	"testing"
)

func TestIntegrationCategoriesCanBeCreated(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("http", "192.168.99.100", 8080)
	if err != nil {
		test.Error(err)
	}

	testCategories := []string{"First test category", "Second test category"}
	ids, err := puffer.CreateCategories(testCategories)
	if err != nil {
		test.Error(err)
	}

	if len(ids) <= 0 {
		test.Fail()
	}
}
