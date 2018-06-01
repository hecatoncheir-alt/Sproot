package storage

import (
	"log"
	"sync"
	"testing"

	"github.com/hecatoncheir/Configuration"
)

var once sync.Once
var storage *Storage

func prepareStorage() {
	var err error

	config := configuration.New()
	storage = New(config.Development.Database.Host,
		config.Development.Database.Port)

	err = storage.SetUp()
	if err != nil {
		log.Fatal(err)
	}
}

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	config := configuration.New()
	storage = New(config.Development.Database.Host, config.Development.Database.Port)

	err := storage.SetUp()
	if err != nil {
		test.Fail()
	}
}
