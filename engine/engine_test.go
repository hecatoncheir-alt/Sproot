package engine

import (
	"testing"
	"time"
)

func TestSprootCanSaveAndGetCompany(test *testing.T) {
	var err error
	puffer := New()

	_, _, err = puffer.DatabaseSetUp("root", "192.168.99.100", 26257, "disable", "Items")
	if err != nil {
		test.Error(err)
	}

	_, err = puffer.GetCompany("Test company")
	if err != ErrCompanyNotExists {
		test.Fail()
	}

	company := Company{
		Name:       "Test company",
		IRI:        "test/",
		Categories: []string{"test company category"},
	}

	companyInStore, err := puffer.SaveCompany(&company)
	if err != nil {
		test.Error(err)
	}

	companyInStore, err = puffer.GetCompany("Test company")
	if err != ErrCompanyNotExists {
		test.Fail()
	}

	if companyInStore.Name != "Test company" {
		test.Fail()
	}
}

func TestSprootCanSaveGetAndDeleteData(test *testing.T) {
	var err error

	puffer := New()
	_, _, err = puffer.DatabaseSetUp("root", "192.168.99.100", 26257, "disable", "Items")
	if err != nil {
		test.Error(err)
	}

	_, err = puffer.GetPricesOfProductsByName("test item name")
	if err != nil {
		test.Fail()
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

	items, err := puffer.GetPricesOfProductsByName("test item name")
	if err != nil {
		test.Fail()
	}

	if len(items) == 0 {
		test.Fail()
	}
}
