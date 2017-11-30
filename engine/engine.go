package engine

import (
	"context"
	"fmt"
	"log"

	dataBaseClient "github.com/dgraph-io/dgraph/client"
	"google.golang.org/grpc"
)

// Engine is a main object of engine pkg
type Engine struct {
	GraphAddress  string
	GraphGRPCHost string
	GraphGRPCPort int
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// DatabaseSetUp is a method for setup SQL database for graph engine
func (engine *Engine) DatabaseSetUp(host string, port int) error {
	engine.GraphGRPCHost = host
	engine.GraphGRPCPort = port
	engine.GraphAddress = fmt.Sprintf("%v:%v", host, port)
	return nil
}

// SetUpIndexes needed for setup schema for dgraph database
func (engine *Engine) SetUpIndexes() error {
	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return err
	}

	operation := &protos.Operation{Schema: `
		name: string @index(exact, term) .
	`}

	err = client.Alter(context.Background(), operation)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// PrepareDataBaseClient make all for needed checks for work with Dgraph database
func (engine *Engine) PrepareDataBaseClient() (*dataBaseClient.Dgraph, error) {
	conn, err := grpc.Dial(engine.GraphAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client.NewDgraphClient(
		protos.NewDgraphClient(conn),
	), nil

}
