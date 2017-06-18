package engine

import (
	"testing"
	"time"
)

func TestSprootCanSaveGetAndDeleteData(test *testing.T) {
	puffer := New()
	_, _, err := puffer.DatabaseSetUp("root", "192.168.99.100", 26257, "disable", "Items")
	if err != nil {
		test.Error(err)
	}

	testCompany := Company{
		Name:       "Test company",
		IRI:        "http://www.test-company.ru/",
		Categories: []string{"Cмартфоны"},
	}

	testPriceData := Price{
		Value:    "100",
		City:     "Moscow",
		DateTime: time.Now().UTC(),
	}

	incomingTestItem := Item{
		Name:             "test item name",
		Price:            testPriceData,
		Link:             "/",
		Company:          testCompany,
		PreviewImageLink: "img/",
	}

	item, err := puffer.SavePriceForProductOfCompany(&incomingTestItem)
	if err != nil {
		test.Error(err)
	}

	if item.Name != "test item name" {
		test.Fail()
	}
}
