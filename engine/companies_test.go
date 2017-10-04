package engine

import "testing"

func TestCompanyCanBeCreated(test *testing.T) {
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("http", "192.168.99.100", 8080)
	if err != nil {
		test.Error(err)
	}

	testCompany := Company{
		Name:       "Test company",
		IRI:        "http://www.test-company.ru/",
		Categories: []string{"Cмартфоны", "Test category"},
	}

	id, err := puffer.CreateCompany(&testCompany)
	if err != nil {
		test.Error(err)
	}

	if id == "" {
		test.Fail()
	}
}