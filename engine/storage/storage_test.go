package storage

import (
	"github.com/hecatoncheir/Sproot/configuration"
	"log"
	"sync"
	"testing"
)

var once sync.Once
var storage *Storage

func prepareStorage() {
	var err error

	config, err := configuration.GetConfiguration()
	storage = New(config.Development.Database.Host, config.Development.Database.Port)

	err = storage.SetUp()
	if err != nil {
		log.Fatal(err)
	}
}

func TestIntegrationStorageCanConnectToDatabase(test *testing.T) {
	config, err := configuration.GetConfiguration()
	storage = New(config.Development.Database.Host, config.Development.Database.Port)

	err = storage.SetUp()
	if err != nil {
		test.Fail()
	}
}
