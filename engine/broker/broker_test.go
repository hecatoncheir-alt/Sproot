package broker

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/google/uuid"

	"github.com/hecatoncheir/Sproot/configuration"
)

func TestBrokerCanConnectToNSQ(test *testing.T) {
	bro := New()
	uuidOfTopic := uuid.New().String()

	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println(err)
	}

	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println("Need started NSQ")
		log.Println(err)
	}

	message, err := json.Marshal(map[string]string{"test key": "test value"})

	bro.Producer.Publish(uuidOfTopic, message)
	defer bro.Producer.Stop()

	items, err := bro.ListenTopic(uuidOfTopic, config.Development.SprootTopic)
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
	bro := New()
	uuidOfTopic := uuid.New().String()

	config, err := configuration.GetConfiguration()
	if err != nil {
		log.Println(err)
	}

	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		log.Println("Need started NSQ")
		log.Println(err)
	}

	item := EventData{Message: "test item"}

	items, err := bro.ListenTopic(uuidOfTopic, "Sproot")
	if err != nil {
		test.Error(err)
	}

	err = bro.WriteToTopic(uuidOfTopic, item)
	if err != nil {
		test.Error(err)
	}

	defer bro.Producer.Stop()

	for item := range items {
		data := EventData{}
		json.Unmarshal(item, &data)
		if data.Message == "test item" {
			break
		}
	}
}
