package broker

import (
	"encoding/json"
	"github.com/hecatoncheir/Sproot/configuration"
	"log"
	"sync"
	"testing"
)

var broker *Broker
var once sync.Once

func SetUp() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println(err)
	}

	broker = New()
	err = broker.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println(err)
	}
}

func TestBrokerCanConnectToNSQ(test *testing.T) {
	once.Do(SetUp)

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	broker.Producer.Publish("test", message)

	items, err := broker.ListenTopic("test", "testing")
	if err != nil {
		test.Error(err)
	}

	for item := range items {
		data := map[string]string{}
		json.Unmarshal(item, &data)
		if data["test key"] == "test value" {
			break
		}
	}
}
