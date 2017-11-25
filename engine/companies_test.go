package engine

import "testing"

func TestIntegrationCompanyCanBeCreated(test *testing.T) {
	// test.Skip()
	// var err error
	// puffer := New()

	// err = puffer.DatabaseSetUp("192.168.99.100", 9080)
	// if err != nil {
	// 	test.Error(err)
	// }

	// err = puffer.SetUpIndexes()
	// if err != nil {
	// 	test.Error(err)
	// }

	// testCategories := []Category{{Name: "First test category"}, {Name: "Second test category"}}

	// testCompany := Company{
	// 	Name:       "Test company",
	// 	IRI:        "http://www.test-company.ru/",
	// 	Categories: testCategories,
	// }

	// company, err := puffer.CreateCompany(&testCompany)
	// if err != nil {
	// 	if err != ErrCategoriesAlreadyExists {
	// 		test.Error(err)
	// 	}
	// }

	// if int(company.ID) == 0 {
	// 	test.Fail()
	// }

	// // puffer.DeleteCategories(company.Categories)
	// puffer.DeleteCompany(&company)
}
