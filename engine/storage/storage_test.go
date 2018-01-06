package storage

import (
	"testing"
	"fmt"
)

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	storage := New(databaseHost, databasePort)
	err := storage.SetUp()
	if err != nil {
		test.Fail()
	}
}
