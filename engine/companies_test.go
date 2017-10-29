package engine

import "testing"

func TestIntegrationCompanyCanBeCreated(test *testing.T) {
	test.Skip()
	var err error
	puffer := New()

	err = puffer.DatabaseSetUp("192.168.99.100", 9080)
	if err != nil {
		test.Error(err)
	}

	testCategories := []string{"First test category", "Second test category"}

	testCompany := Company{
		Name:       "Test company",
		IRI:        "http://www.test-company.ru/",
		Categories: testCategories,
	}

	id, err := puffer.CreateCompany(&testCompany)
	if err != nil {
		test.Error(err)
	}

	if id == "" {
		test.Fail()
	}
}
