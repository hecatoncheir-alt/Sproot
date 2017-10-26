package engine

import (
	"fmt"
	"log"
	"io/ioutil"
	dataBaseClient "github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/x"
	"google.golang.org/grpc"
	"os"
)

// Engine is a main object of engine pkg
type Engine struct {
	GraphAddress string
	GraphIRI     string
	GraphHost    string
	GraphPort    int
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// DatabaseSetUp is a method for setup SQL database for graph engine
func (engine *Engine) DatabaseSetUp(protocol string, host string, port int) error {

	engine.GraphAddress = fmt.Sprintf("%v:%v", host, port)
	engine.GraphIRI = fmt.Sprintf("%v://%v:%v", protocol, host, port)
	engine.GraphHost = host
	engine.GraphPort = port

	return nil
}

func (engine *Engine) PrepareDataBaseClient() (*dataBaseClient.Dgraph, *dataBaseClient.Req, error) {
	conn, err := grpc.Dial(engine.GraphAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	x.Checkf(err, "While trying to dial gRPC")
	defer conn.Close()

	clientDir, err := ioutil.TempDir("", "client_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(clientDir)

	client := dataBaseClient.NewDgraphClient([]*grpc.ClientConn{conn}, dataBaseClient.DefaultOptions, clientDir)
	defer client.Close()

	return client, &dataBaseClient.Req{}, nil
}
