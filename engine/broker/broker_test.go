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

	config, _ := configuration.GetConfiguration()

	items, err := broker.ListenTopic("test", config.Development.Channel)
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

func TestBrokerCanSendMessageToNSQ(test *testing.T) {
	var err error
	once.Do(SetUp)

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	items, err := broker.ListenTopic("test", "test")
	if err != nil {
		test.Error(err)
	}

	err = broker.WriteToTopic("test", message)
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
