package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	nsq "github.com/bitly/go-nsq"
)

// New constructor for Broker
func New() *Broker {
	broker := Broker{}
	broker.configuration = nsq.NewConfig()
	// broker.сonfiguration.MaxInFlight = 6
	// broker.сonfiguration.MsgTimeout = time.Duration(time.Second * 6)
	return &broker
}

// Broker is a object of message stream
type Broker struct {
	IP            string
	Port          int
	configuration *nsq.Config
	Producer      *nsq.Producer
}

// connectToMessageBroker method for connect to message broker
func (broker *Broker) connectToMessageBroker(host string, port int) (*nsq.Producer, error) {
	if host != "" && string(port) != "" {
		broker.IP = host
		broker.Port = port
	}

	hostAddr := fmt.Sprintf("%v:%v", broker.IP, strconv.Itoa(broker.Port))
	producer, err := nsq.NewProducer(hostAddr, broker.configuration)

	if err != nil {
		log.Print("Could not connect to message broker")
	}

	return producer, err

}

// Connect to message broker for publish events
func (broker *Broker) Connect(host string, port int) error {
	producer, err := broker.connectToMessageBroker(host, port)
	broker.Producer = producer
	return err
}

// WriteToTopic method for publish message to topic
func (broker *Broker) WriteToTopic(topic string, message interface{}) error {
	event, err := json.Marshal(message)
	if err != nil {
		return err
	}

	go broker.Producer.Publish(topic, event)
	return nil
}

// ListenTopic get events in channel of topic
func (broker *Broker) ListenTopic(topic string, channel string) (<-chan []byte, error) {
	consumer, err := nsq.NewConsumer(topic, channel, broker.configuration)
	if err != nil {
		return nil, err
	}

	events := make(chan []byte, 6)

	handler := nsq.HandlerFunc(func(message *nsq.Message) error {
		events <- message.Body
		return nil
	})

	consumer.AddConcurrentHandlers(handler, 6)

	hostAddr := fmt.Sprintf("%v:%v", broker.IP, strconv.Itoa(broker.Port))
	go consumer.ConnectToNSQD(hostAddr)

	return events, nil
}
