package engine

import (
	"github.com/hecatoncheir/Sproot/engine/broker"
	"github.com/hecatoncheir/Sproot/engine/storage"
)

// Engine is a main object of engine pkg
type Engine struct {
	Storage *storage.Storage
	Broker  *broker.Broker
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// SetUpStorage for make connect to database and prepare client for requests
func (engine *Engine) SetUpStorage(host string, port int) error {
	engine.Storage = storage.New(host, port)
	err := engine.Storage.SetUp()
	if err != nil {
		return err
	}

	return nil
}

// SetUpBroker for make connect to broker and prepare client for requests
func (engine *Engine) SetUpBroker(host string, port int) error {
	bro := broker.New()
	engine.Broker = bro

	err := bro.Connect(host, port)
	if err != nil {
		return err
	}

	return nil
}
