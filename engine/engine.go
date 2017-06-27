package engine

import (
	"fmt"
)

// Engine is a main object of engine pkg
type Engine struct {
	GraphAddress string
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// DatabaseSetUp is a method for setup SQL database for graph engine
func (engine *Engine) DatabaseSetUp(protocol string, host string, port int) error {

	dbAddr := fmt.Sprintf("%v://%v:%v", protocol, host, port)
	engine.GraphAddress = dbAddr

	return nil
}
