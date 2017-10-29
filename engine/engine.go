package engine

import (
	"fmt"
	"log"
	"io/ioutil"
	dataBaseClient "github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/x"
	"google.golang.org/grpc"
	"os"
	"context"
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

func (engine *Engine) SetUpIndexes() error {
	client, err := engine.PrepareDataBaseClient()
	if err != nil {
		log.Println(err)
		return err
	}

	defer client.Close()

	request := &dataBaseClient.Req{}

	request.SetSchema(`
				name: string @index(exact, term) .
	`)

	_, err = client.Run(context.Background(), request)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (engine *Engine) PrepareDataBaseClient() (*dataBaseClient.Dgraph, error) {
	conn, err := grpc.Dial(engine.GraphAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	x.Checkf(err, "While trying to dial gRPC")

	clientDir, err := ioutil.TempDir("", "sproot_dgraph_client_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(clientDir)

	client := dataBaseClient.NewDgraphClient([]*grpc.ClientConn{conn}, dataBaseClient.DefaultOptions, clientDir)

	return client, nil
}
