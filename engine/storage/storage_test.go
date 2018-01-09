package storage

import (
	"testing"
	"log"
	"sync"
)

var once sync.Once
var storage *Storage

func prepareStorage() {
	var err error
	storage = New(databaseHost, databasePort)

	err = storage.SetUp()
	if err != nil {
		log.Fatal(err)
	}
}

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	storage := New(databaseHost, databasePort)
	err := storage.SetUp()
	if err != nil {
		test.Fail()
	}
}
