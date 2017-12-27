package storage

import (
	"testing"
)

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	storage := New(databaseHost, databasePort)
	err := storage.SetUp()
	if err != nil {
		test.Fail()
	}
}

// TODO
func TestIntegrationStorageCanGetCompanyWithCategories(test *testing.T) {
	test.Skip()
}
